package suites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/multsig"
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

/*

****************************************************************************************************************************************
Taken from deprecated spec definition: https://filecoin-project.github.io/specs/docs/systems/filecoin_token/multisig/multisig_actor_old/
****************************************************************************************************************************************

A basic multisig account actor. Allows sending of messages like a normal account actor, but with the requirement of
M of N parties agreeing to the operation. Completed and/or cancelled operations stick around in the actors state until
explicitly cleared out. Proposers may cancel transactions they propose, or transactions by proposers who are no longer
approved signers.

Self modification methods (add/remove signer, change requirement) are called by doing a multisig transaction invoking
the desired method on the contract itself. This means the ‘signature threshold’ logic only needs to be implemented
once, in one place.

The initialize actor is used to create new instances of the multisig.
*/

func MultiSigActorConstructor(t testing.TB, factory Factories) {
	const aliceBal = 200000000000
	const valueSend = 10
	const requiredSigners = 3
	const unlockDuration = 10

	w := multiSigTestSetup(t, factory)

	alice := w.Driver.NewAccountActor(aliceBal)

	multisigAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	mustCreateMultisigActor(w, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice)
	w.Driver.AssertBalance(multisigAddr, valueSend)
	w.Driver.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice},
		Required:       requiredSigners,
		NextTxID:       0,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions:   nil,
	})
}

func MultiSigActorProposeApprove(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const valueSend = 10
	const requiredSigners = 2
	const unlockDuration = 10

	w := multiSigTestSetup(t, factory)

	alice := w.Driver.NewAccountActor(initialBal)
	bob := w.Driver.NewAccountActor(initialBal)

	multisigAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// create a multisig actor with a balance of 'valueSend' FIL.
	mustCreateMultisigActor(w, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice, bob)
	w.Driver.AssertBalance(multisigAddr, valueSend)
	w.Driver.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice, bob},
		Required:       requiredSigners,
		NextTxID:       0,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions:   nil,
	})
	w.Driver.AssertBalance(multisigAddr, valueSend)

	// alice proposes that outsider should receive 'valueSend' FIL.
	txID0 := multsig.MultiSigTxID{TxID: 0}
	outsider := w.Driver.NewAccountActor(initialBal)

	mustProposeMultisigTransfer(w, 1, 0, txID0, multisigAddr, alice, outsider, valueSend)
	w.Driver.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice, bob},
		Required:       requiredSigners,
		NextTxID:       1,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions: []multsig.MTransaction{{
			Created:  0,
			TxID:     txID0.TxID,
			To:       outsider,
			Value:    types.NewInt(valueSend),
			Method:   0,
			Params:   []byte{},
			Approved: []address.Address{alice},
			Complete: false,
			Canceled: false,
			RetCode:  0,
		}},
	})
	w.Driver.AssertBalance(multisigAddr, valueSend)

	// outsider proposes themselves to receive 'valueSend' FIL. This fails as they are not a signer.
	msg, err := w.Producer.MultiSigPropose(multisigAddr, outsider, 0, outsider, types.NewInt(valueSend), 0, []byte{}, chain.Value(0))
	require.NoError(t, err)
	mr, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.EqualError(w.T, err, "not authorized (RetCode=1)")
	w.Driver.AssertReceipt(mr, chain.MessageReceipt{
		ExitCode:    1,
		ReturnValue: nil,
		GasUsed:     0,
	})
	w.Driver.AssertBalance(multisigAddr, valueSend)

	// outsider approves the value transfer alice sent. This fails as they are not a signer.
	msg, err = w.Producer.MultiSigApprove(multisigAddr, outsider, 1, txID0.TxID, chain.Value(0))
	require.NoError(t, err)
	mr, err = w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.EqualError(w.T, err, "not authorized (RetCode=1)")
	w.Driver.AssertReceipt(mr, chain.MessageReceipt{
		ExitCode:    1,
		ReturnValue: nil,
		GasUsed:     0,
	})

	// bob approves transfer of 'valueSend' FIL to outsider.
	mustApproveMultisigActor(w, 0, 0, multisigAddr, bob, txID0)
	w.Driver.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice, bob},
		Required:       requiredSigners,
		NextTxID:       1,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions: []multsig.MTransaction{{
			Created:  0,
			TxID:     txID0.TxID,
			To:       outsider,
			Value:    types.NewInt(valueSend),
			Method:   0,
			Params:   []byte{},
			Approved: []address.Address{alice, bob},
			Complete: true,
			Canceled: false,
			RetCode:  0,
		}},
	})
	w.Driver.AssertBalance(multisigAddr, 0)

}

func MultiSigActorProposeCancel(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const valueSend = 10
	const requiredSigners = 2
	const unlockDuration = 10

	w := multiSigTestSetup(t, factory)
	alice := w.Driver.NewAccountActor(initialBal)
	bob := w.Driver.NewAccountActor(initialBal)

	multisigAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// create a multisig actor with a balance of 'valueSend' FIL.
	mustCreateMultisigActor(w, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice, bob)
	w.Driver.AssertBalance(multisigAddr, valueSend)
	w.Driver.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice, bob},
		Required:       requiredSigners,
		NextTxID:       0,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions:   nil,
	})

	// alice proposes that outsider should receive 'valueSend' FIL.
	txID0 := multsig.MultiSigTxID{TxID: 0}
	outsider := w.Driver.NewAccountActor(initialBal)
	mustProposeMultisigTransfer(w, 1, 0, txID0, multisigAddr, alice, outsider, valueSend)
	w.Driver.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice, bob},
		Required:       requiredSigners,
		NextTxID:       1,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions: []multsig.MTransaction{{
			Created:  0,
			TxID:     txID0.TxID,
			To:       outsider,
			Value:    types.NewInt(valueSend),
			Method:   0,
			Params:   []byte{},
			Approved: []address.Address{alice},
			Complete: false,
			Canceled: false,
			RetCode:  0,
		}},
	})

	// bob cancels alice's transaction. This fails as bob did not create alice's transaction.
	msg, err := w.Producer.MultiSigCancel(multisigAddr, bob, 0, txID0.TxID, chain.Value(0))
	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.EqualError(t, err, "cannot cancel another signers transaction (RetCode=4)")
	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    4,
		ReturnValue: nil,
		GasUsed:     0,
	})
	w.Driver.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice, bob},
		Required:       requiredSigners,
		NextTxID:       1,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions: []multsig.MTransaction{{
			Created:  0,
			TxID:     txID0.TxID,
			To:       outsider,
			Value:    types.NewInt(valueSend),
			Method:   0,
			Params:   []byte{},
			Approved: []address.Address{alice},
			Complete: false,
			Canceled: false,
			RetCode:  0,
		}},
	})

	// alice cancels their transaction. The outsider doesn't receive any FIL, the multisig actor's balance is empty, and the
	// transaction is canceled.
	mustCancelMultisigActor(w, 2, 0, multisigAddr, alice, txID0)
	w.Driver.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice, bob},
		Required:       requiredSigners,
		NextTxID:       1,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions: []multsig.MTransaction{{
			Created:  0,
			TxID:     txID0.TxID,
			To:       outsider,
			Value:    types.NewInt(valueSend),
			Method:   0,
			Params:   []byte{},
			Approved: []address.Address{alice},
			Complete: false,
			Canceled: true,
			RetCode:  0,
		}},
	})
	w.Driver.AssertBalance(multisigAddr, valueSend)
	w.Driver.AssertBalance(outsider, initialBal)
}

type multiSigTestingWrapper struct {
	T         testing.TB
	Driver    *StateDriver
	Producer  *chain.MessageProducer
	Validator *chain.Validator
	ExeCtx    *chain.ExecutionContext
}

func multiSigTestSetup(t testing.TB, factory Factories) *multiSigTestingWrapper {
	drv := NewStateDriver(t, factory.NewState())
	gasPrice := types.NewInt(1)
	gasLimit := types.GasUnit(1000000)

	_, _, err := drv.State().SetSingletonActor(actors.InitAddress, types.NewInt(0))
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.BurntFundsAddress, types.NewInt(0))
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.NetworkAddress, TotalNetworkBalance)
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.StoragePowerAddress, types.NewInt(0))
	require.NoError(t, err)

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	testMiner := drv.NewAccountActor(0)
	exeCtx := chain.NewExecutionContext(1, testMiner)

	return &multiSigTestingWrapper{
		T:         t,
		Driver:    drv,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}

}

func mustProposeMultisigTransfer(gdg *multiSigTestingWrapper, nonce, value uint64, txID multsig.MultiSigTxID, to, from, proposeTo address.Address, proposeValue uint64) {
	msg, err := gdg.Producer.MultiSigPropose(to, from, nonce, proposeTo, types.NewInt(proposeValue), 0, []byte{}, chain.Value(value))
	msgReceipt, err := gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Driver.State(), msg)
	require.NoError(gdg.T, err)

	btxid, err := types.Serialize(&txID)
	require.NoError(gdg.T, err)

	gdg.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode: 0,
		// since the first byte is the cbor type indicator.
		ReturnValue: btxid[1:],
		GasUsed:     0,
	})
}

func mustCreateMultisigActor(gdg *multiSigTestingWrapper, nonce, value uint64, required, unlockDuration uint64, ms, creator address.Address, signers ...address.Address) {
	multiSigConstuctParams, err := types.Serialize(&multsig.MultiSigConstructorParams{
		Signers:        append(signers, creator),
		Required:       required,
		UnlockDuration: unlockDuration,
	})
	require.NoError(gdg.T, err)

	msg, err := gdg.Producer.InitExec(creator, nonce, actors.MultisigActorCodeCid, multiSigConstuctParams, chain.Value(value))
	require.NoError(gdg.T, err)

	msgReceipt, err := gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Driver.State(), msg)
	require.NoError(gdg.T, err)

	gdg.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: ms.Bytes(),
		GasUsed:     0,
	})
}

func mustApproveMultisigActor(gdg *multiSigTestingWrapper, nonce, value uint64, ms, from address.Address, txID multsig.MultiSigTxID) {
	msg, err := gdg.Producer.MultiSigApprove(ms, from, nonce, txID.TxID, chain.Value(0))
	require.NoError(gdg.T, err)

	msgReceipt, err := gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Driver.State(), msg)
	require.NoError(gdg.T, err)

	gdg.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})
}

func mustCancelMultisigActor(gdg *multiSigTestingWrapper, nonce, value uint64, ms, from address.Address, txID multsig.MultiSigTxID) {
	msg, err := gdg.Producer.MultiSigCancel(ms, from, nonce, txID.TxID, chain.Value(value))
	require.NoError(gdg.T, err)

	msgReceipt, err := gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Driver.State(), msg)
	require.NoError(gdg.T, err)

	gdg.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})
}

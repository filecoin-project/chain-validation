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

func multiSigTestSetup(t testing.TB, factory Factories) (*StateDriver, types.BigInt, types.GasUnit) {
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

	return drv, gasPrice, gasLimit

}

type Gadget struct {
	T         testing.TB
	Drv       *StateDriver
	Producer  *chain.MessageProducer
	Validator *chain.Validator
	ExeCtx    *chain.ExecutionContext
}

func MultiSigActorConstructor(t testing.TB, factory Factories) {
	const aliceBal = 200000000000
	const valueSend = 10
	const requiredSigners = 3
	const unlockDuration = 10

	drv, gasPrice, gasLimit := multiSigTestSetup(t, factory)

	// miner that mines in this test

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	testMiner := drv.NewAccountActor(0)
	exeCtx := chain.NewExecutionContext(1, testMiner)

	gdg := &Gadget{
		T:         t,
		Drv:       drv,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}

	alice := drv.NewAccountActor(aliceBal)

	multisigAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	mustCreateMultisigActor(gdg, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice)

	drv.AssertBalance(multisigAddr, valueSend)
	drv.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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

	drv, gasPrice, gasLimit := multiSigTestSetup(t, factory)

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	testMiner := drv.NewAccountActor(0)
	exeCtx := chain.NewExecutionContext(1, testMiner)

	gdg := &Gadget{
		T:         t,
		Drv:       drv,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}

	alice := drv.NewAccountActor(initialBal)
	bob := drv.NewAccountActor(initialBal)

	multisigAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// create a multisig actor with a balance of 'valueSend' FIL.
	mustCreateMultisigActor(gdg, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice, bob)
	drv.AssertBalance(multisigAddr, valueSend)
	drv.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice, bob},
		Required:       requiredSigners,
		NextTxID:       0,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions:   nil,
	})
	drv.AssertBalance(multisigAddr, valueSend)

	// alice proposes that outsider should receive 'valueSend' FIL.
	txID0 := multsig.MultiSigTxID{TxID: 0}
	outsider := drv.NewAccountActor(initialBal)
	mustProposeMultisigTransfer(gdg, 1, 0, txID0, multisigAddr, alice, outsider, valueSend)
	drv.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	drv.AssertBalance(multisigAddr, valueSend)

	// outsider proposes themselves to receive 'valueSend' FIL. This fails as they are not a signer.
	msg, err := gdg.Producer.MultiSigPropose(multisigAddr, outsider, 0, outsider, types.NewInt(valueSend), 0, []byte{}, chain.Value(0))
	mr, err := gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Drv.State(), msg)
	require.EqualError(gdg.T, err, "not authorized (RetCode=1)")
	drv.AssertReceipt(mr, chain.MessageReceipt{
		ExitCode:    1,
		ReturnValue: nil,
		GasUsed:     0,
	})
	drv.AssertBalance(multisigAddr, valueSend)

	// outsider approves the value transfer alice sent. This fails as they are not a signer.
	msg, err = gdg.Producer.MultiSigApprove(multisigAddr, outsider, 1, txID0.TxID, chain.Value(0))
	mr, err = gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Drv.State(), msg)
	require.EqualError(gdg.T, err, "not authorized (RetCode=1)")
	drv.AssertReceipt(mr, chain.MessageReceipt{
		ExitCode:    1,
		ReturnValue: nil,
		GasUsed:     0,
	})

	// bob approves transfer of 'valueSend' FIL to outsider.
	msg, err = gdg.Producer.MultiSigApprove(multisigAddr, bob, 0, txID0.TxID, chain.Value(0))
	msgReceipt, err := gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})
	drv.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	drv.AssertBalance(multisigAddr, 0)

}

func MultiSigActorProposeCancel(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const valueSend = 10
	const requiredSigners = 2
	const unlockDuration = 10

	drv, gasPrice, gasLimit := multiSigTestSetup(t, factory)

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	testMiner := drv.NewAccountActor(0)
	exeCtx := chain.NewExecutionContext(1, testMiner)

	gdg := &Gadget{
		T:         t,
		Drv:       drv,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}

	alice := drv.NewAccountActor(initialBal)
	bob := drv.NewAccountActor(initialBal)

	multisigAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// create a multisig actor with a balance of 'valueSend' FIL.
	mustCreateMultisigActor(gdg, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice, bob)
	drv.AssertBalance(multisigAddr, valueSend)
	drv.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	outsider := drv.NewAccountActor(initialBal)
	mustProposeMultisigTransfer(gdg, 1, 0, txID0, multisigAddr, alice, outsider, valueSend)
	drv.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	msg, err := gdg.Producer.MultiSigCancel(multisigAddr, bob, 0, txID0.TxID, chain.Value(0))
	msgReceipt, err := gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Drv.State(), msg)
	require.EqualError(t, err, "cannot cancel another signers transaction (RetCode=4)")
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    4,
		ReturnValue: nil,
		GasUsed:     0,
	})
	drv.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	msg, err = gdg.Producer.MultiSigCancel(multisigAddr, alice, 2, txID0.TxID, chain.Value(0))
	msgReceipt, err = gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})
	drv.AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	drv.AssertBalance(multisigAddr, valueSend)
	drv.AssertBalance(outsider, initialBal)
}

func mustProposeMultisigTransfer(gdg *Gadget, nonce, value uint64, txID multsig.MultiSigTxID, to, from, proposeTo address.Address, proposeValue uint64) {
	msg, err := gdg.Producer.MultiSigPropose(to, from, nonce, proposeTo, types.NewInt(proposeValue), 0, []byte{}, chain.Value(value))
	msgReceipt, err := gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Drv.State(), msg)
	require.NoError(gdg.T, err)

	btxid, err := types.Serialize(&txID)
	require.NoError(gdg.T, err)

	gdg.Drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode: 0,
		// since the first byte is the cbor type indicator.
		ReturnValue: btxid[1:],
		GasUsed:     0,
	})
}

func mustCreateMultisigActor(gdg *Gadget, nonce, value uint64, required, unlockDuration uint64, ms, creator address.Address, signers ...address.Address) {
	multiSigConstuctParams, err := types.Serialize(&multsig.MultiSigConstructorParams{
		Signers:        append(signers, creator),
		Required:       required,
		UnlockDuration: unlockDuration,
	})
	require.NoError(gdg.T, err)

	msg, err := gdg.Producer.InitExec(creator, nonce, actors.MultisigActorCodeCid, multiSigConstuctParams, chain.Value(value))
	require.NoError(gdg.T, err)

	msgReceipt, err := gdg.Validator.ApplyMessage(gdg.ExeCtx, gdg.Drv.State(), msg)
	require.NoError(gdg.T, err)

	require.NoError(gdg.T, err)
	gdg.Drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: ms.Bytes(),
		GasUsed:     0,
	})
}

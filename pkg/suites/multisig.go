package suites

import (
	"github.com/filecoin-project/chain-validation/pkg/state"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/multsig"
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

	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})

	alice := td.Driver().NewAccountActor(aliceBal)

	multisigAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	mustCreateMultisigActor(td, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice)
	td.Driver().AssertBalance(multisigAddr, valueSend)
	td.Driver().AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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

	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})

	alice := td.Driver().NewAccountActor(initialBal)
	bob := td.Driver().NewAccountActor(initialBal)

	multisigAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// create a multisig actor with a balance of 'valueSend' FIL.
	mustCreateMultisigActor(td, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice, bob)
	td.Driver().AssertBalance(multisigAddr, valueSend)
	td.Driver().AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
		Signers:        []address.Address{alice, bob},
		Required:       requiredSigners,
		NextTxID:       0,
		InitialBalance: types.NewInt(valueSend),
		StartingBlock:  1,
		UnlockDuration: unlockDuration,
		Transactions:   nil,
	})
	td.Driver().AssertBalance(multisigAddr, valueSend)

	// alice proposes that outsider should receive 'valueSend' FIL.
	txID0 := multsig.MultiSigTxID{TxID: 0}
	outsider := td.Driver().NewAccountActor(initialBal)

	mustProposeMultisigTransfer(td, 1, 0, txID0, multisigAddr, alice, outsider, valueSend)
	td.Driver().AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	td.Driver().AssertBalance(multisigAddr, valueSend)

	// outsider proposes themselves to receive 'valueSend' FIL. This fails as they are not a signer.
	msg, err := td.Producer().MultiSigPropose(multisigAddr, outsider, 0, outsider, types.NewInt(valueSend), 0, []byte{}, chain.Value(0))
	require.NoError(t, err)
	mr, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.EqualError(td.TB(), err, "not authorized (RetCode=1)")
	td.Driver().AssertReceipt(mr, chain.MessageReceipt{
		ExitCode:    1,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
	td.Driver().AssertBalance(multisigAddr, valueSend)

	// outsider approves the value transfer alice sent. This fails as they are not a signer.
	msg, err = td.Producer().MultiSigApprove(multisigAddr, outsider, 1, txID0.TxID, chain.Value(0))
	require.NoError(t, err)
	mr, err = td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.EqualError(td.TB(), err, "not authorized (RetCode=1)")
	td.Driver().AssertReceipt(mr, chain.MessageReceipt{
		ExitCode:    1,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})

	// bob approves transfer of 'valueSend' FIL to outsider.
	mustApproveMultisigActor(td, 0, 0, multisigAddr, bob, txID0)
	td.Driver().AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	td.Driver().AssertBalance(multisigAddr, 0)

}

func MultiSigActorProposeCancel(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const valueSend = 10
	const requiredSigners = 2
	const unlockDuration = 10

	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})

	alice := td.Driver().NewAccountActor(initialBal)
	bob := td.Driver().NewAccountActor(initialBal)

	multisigAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// create a multisig actor with a balance of 'valueSend' FIL.
	mustCreateMultisigActor(td, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice, bob)
	td.Driver().AssertBalance(multisigAddr, valueSend)
	td.Driver().AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	outsider := td.Driver().NewAccountActor(initialBal)
	mustProposeMultisigTransfer(td, 1, 0, txID0, multisigAddr, alice, outsider, valueSend)
	td.Driver().AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	msg, err := td.Producer().MultiSigCancel(multisigAddr, bob, 0, txID0.TxID, chain.Value(0))
	require.NoError(t, err)
	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.EqualError(t, err, "cannot cancel another signers transaction (RetCode=4)")
	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    4,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
	td.Driver().AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	mustCancelMultisigActor(td, 2, 0, multisigAddr, alice, txID0)
	td.Driver().AssertMultisigState(multisigAddr, multsig.MultiSigActorState{
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
	td.Driver().AssertBalance(multisigAddr, valueSend)
	td.Driver().AssertBalance(outsider, initialBal)
}

func mustProposeMultisigTransfer(td TestDriver, nonce, value uint64, txID multsig.MultiSigTxID, to, from, proposeTo address.Address, proposeValue uint64) {
	msg, err := td.Producer().MultiSigPropose(to, from, nonce, proposeTo, types.NewInt(proposeValue), 0, []byte{}, chain.Value(value))
	require.NoError(td.TB(), err)
	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	btxid, err := state.Serialize(&txID)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode: 0,
		// since the first byte is the cbor type indicator.
		ReturnValue: btxid[1:],
		GasUsed:     types.NewInt(0),
	})
}

func mustCreateMultisigActor(td TestDriver, nonce, value uint64, required, unlockDuration uint64, ms, creator address.Address, signers ...address.Address) {
	multiSigConstuctParams, err := state.Serialize(&multsig.MultiSigConstructorParams{
		Signers:        append(signers, creator),
		Required:       required,
		UnlockDuration: unlockDuration,
	})
	require.NoError(td.TB(), err)

	msg, err := td.Producer().InitExec(creator, nonce, actors.MultisigActorCodeCid, multiSigConstuctParams, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: ms.Bytes(),
		GasUsed:     types.NewInt(0),
	})
}

func mustApproveMultisigActor(td TestDriver, nonce, value uint64, ms, from address.Address, txID multsig.MultiSigTxID) {
	msg, err := td.Producer().MultiSigApprove(ms, from, nonce, txID.TxID, chain.Value(0))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

func mustCancelMultisigActor(td TestDriver, nonce, value uint64, ms, from address.Address, txID multsig.MultiSigTxID) {
	msg, err := td.Producer().MultiSigCancel(ms, from, nonce, txID.TxID, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

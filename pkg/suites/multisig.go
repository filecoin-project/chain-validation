package suites

import (
	"context"
	"fmt"
	"testing"

	address "github.com/filecoin-project/go-address"
	block "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	require "github.com/stretchr/testify/require"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"

	chain "github.com/filecoin-project/chain-validation/pkg/chain"
	state "github.com/filecoin-project/chain-validation/pkg/state"
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

type mockStore struct {
	data map[cid.Cid]block.Block
}

func newMockBlocks() *mockStore {
	return &mockStore{make(map[cid.Cid]block.Block)}
}

func (mb *mockStore) Get(c cid.Cid) (block.Block, error) {
	d, ok := mb.data[c]
	if ok {
		return d, nil
	}
	return nil, fmt.Errorf("Not Found")
}

func (mb *mockStore) Put(b block.Block) error {
	mb.data[b.Cid()] = b
	return nil
}

type contextStore struct {
	cbor.IpldStore
	ctx context.Context
}

func newContextStore(ctx context.Context) *contextStore {
	return &contextStore{
		IpldStore: cbor.NewCborStore(newMockBlocks()),
		ctx:       ctx,
	}
}

func (cs *contextStore) Context() context.Context {
	return cs.ctx
}

func MultiSigActorConstructor(t testing.TB, factory Factories) {
	const aliceBal = 200000000000
	const valueSend = 10
	const requiredSigners = 3
	const unlockDuration = 10

	pendingTxMap, err := adt_spec.MakeEmptyMap(newContextStore(context.Background()))
	require.NoError(t, err)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr:         big_spec.NewInt(0),
		builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
		builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
		builtin_spec.RewardActorAddr:       TotalNetworkBalance,
	})

	alice := td.Driver().NewAccountActor(aliceBal)

	multisigAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	mustCreateMultisigActor(td, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice)
	td.Driver().AssertBalance(multisigAddr, valueSend)
	td.Driver().AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice},
		NumApprovalsThreshold: requiredSigners,
		NextTxnID:             0,
		InitialBalance:        abi_spec.NewTokenAmount(valueSend),
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
		PendingTxns:           pendingTxMap.Root(),
	})
}

func MultiSigActorProposeApprove(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const valueSend = 10
	const requiredSigners = 2
	const unlockDuration = 10

	pendingTxMap, err := adt_spec.MakeEmptyMap(newContextStore(context.Background()))
	require.NoError(t, err)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr:         big_spec.NewInt(0),
		builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
		builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
		builtin_spec.RewardActorAddr:       TotalNetworkBalance,
	})

	alice := td.Driver().NewAccountActor(initialBal)
	bob := td.Driver().NewAccountActor(initialBal)

	multisigAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// create a multisig actor with a balance of 'valueSend' FIL.
	mustCreateMultisigActor(td, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice, bob)
	td.Driver().AssertBalance(multisigAddr, valueSend)
	td.Driver().AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: requiredSigners,
		NextTxnID:             0,
		InitialBalance:        abi_spec.NewTokenAmount(valueSend),
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
		PendingTxns:           pendingTxMap.Root(),
	})
	td.Driver().AssertBalance(multisigAddr, valueSend)

	// alice proposes that outsider should receive 'valueSend' FIL.
	outsider := td.Driver().NewAccountActor(initialBal)
	txID0 := multisig_spec.TxnID(0)
	err = pendingTxMap.Put(txID0, &multisig_spec.MultiSigTransaction{
		To:       outsider,
		Value:    abi_spec.NewTokenAmount(valueSend),
		Method:   builtin_spec.MethodSend,
		Params:   []byte{},
		Approved: []address.Address{alice},
	})
	require.NoError(t, err)

	mustProposeMultisigTransfer(td, 1, 0, txID0, multisigAddr, alice, outsider, valueSend)
	td.Driver().AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: requiredSigners,
		NextTxnID:             1,
		InitialBalance:        abi_spec.NewTokenAmount(valueSend),
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
		PendingTxns:           pendingTxMap.Root(),
	})
	td.Driver().AssertBalance(multisigAddr, valueSend)

	// outsider proposes themselves to receive 'valueSend' FIL. This fails as they are not a signer.
	msg, err := td.Producer().MultiSigPropose(multisigAddr, outsider, 0, outsider, abi_spec.NewTokenAmount(valueSend), 0, []byte{}, chain.Value(0))
	require.NoError(t, err)
	mr, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.EqualError(td.TB(), err, "not authorized (RetCode=1)")
	td.Driver().AssertReceipt(mr, chain.MessageReceipt{
		ExitCode:    1,
		ReturnValue: nil,
		GasUsed:     big_spec.NewInt(0),
	})
	td.Driver().AssertBalance(multisigAddr, valueSend)

	// outsider approves the value transfer alice sent. This fails as they are not a signer.
	msg, err = td.Producer().MultiSigApprove(multisigAddr, outsider, 1, txID0, chain.Value(0))
	require.NoError(t, err)
	mr, err = td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.EqualError(td.TB(), err, "not authorized (RetCode=1)")
	td.Driver().AssertReceipt(mr, chain.MessageReceipt{
		ExitCode:    1,
		ReturnValue: nil,
		GasUsed:     big_spec.NewInt(0),
	})

	// bob approves transfer of 'valueSend' FIL to outsider.
	mustApproveMultisigActor(td, 0, 0, multisigAddr, bob, txID0)
	require.NoError(t, pendingTxMap.Delete(txID0))
	td.Driver().AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: requiredSigners,
		NextTxnID:             txID0,
		InitialBalance:        abi_spec.NewTokenAmount(valueSend),
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
		PendingTxns:           pendingTxMap.Root(),
	})
	td.Driver().AssertBalance(multisigAddr, 0)

}

func MultiSigActorProposeCancel(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const valueSend = 10
	const requiredSigners = 2
	const unlockDuration = 10

	pendingTxMap, err := adt_spec.MakeEmptyMap(newContextStore(context.Background()))
	require.NoError(t, err)

	require.NoError(t, err)
	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr:         big_spec.NewInt(0),
		builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
		builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
		builtin_spec.RewardActorAddr:       TotalNetworkBalance,
	})

	alice := td.Driver().NewAccountActor(initialBal)
	bob := td.Driver().NewAccountActor(initialBal)

	multisigAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// create a multisig actor with a balance of 'valueSend' FIL.
	mustCreateMultisigActor(td, 0, valueSend, requiredSigners, unlockDuration, multisigAddr, alice, bob)
	td.Driver().AssertBalance(multisigAddr, valueSend)
	td.Driver().AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: requiredSigners,
		NextTxnID:             0,
		InitialBalance:        abi_spec.NewTokenAmount(valueSend),
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
		PendingTxns:           pendingTxMap.Root(),
	})

	// alice proposes that outsider should receive 'valueSend' FIL.
	txID0 := multisig_spec.TxnID(0)
	outsider := td.Driver().NewAccountActor(initialBal)
	mustProposeMultisigTransfer(td, 1, 0, txID0, multisigAddr, alice, outsider, valueSend)

	require.NoError(t, pendingTxMap.Put(txID0, &multisig_spec.MultiSigTransaction{
		To:       outsider,
		Value:    abi_spec.NewTokenAmount(valueSend),
		Method:   builtin_spec.MethodSend,
		Params:   []byte{},
		Approved: []address.Address{alice},
	}))

	td.Driver().AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: requiredSigners,
		NextTxnID:             1,
		InitialBalance:        abi_spec.NewTokenAmount(valueSend),
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
		PendingTxns:           pendingTxMap.Root(),
	})

	// bob cancels alice's transaction. This fails as bob did not create alice's transaction.
	msg, err := td.Producer().MultiSigCancel(multisigAddr, bob, 0, txID0, chain.Value(0))
	require.NoError(t, err)
	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.EqualError(t, err, "cannot cancel another signers transaction (RetCode=4)")
	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    4,
		ReturnValue: nil,
		GasUsed:     big_spec.NewInt(0),
	})
	td.Driver().AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: requiredSigners,
		NextTxnID:             1,
		InitialBalance:        abi_spec.NewTokenAmount(valueSend),
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
		PendingTxns:           pendingTxMap.Root(),
	})

	// alice cancels their transaction. The outsider doesn't receive any FIL, the multisig actor's balance is empty, and the
	// transaction is canceled.
	mustCancelMultisigActor(td, 2, 0, multisigAddr, alice, txID0)
	require.NoError(t, pendingTxMap.Delete(txID0))
	td.Driver().AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: requiredSigners,
		NextTxnID:             1,
		InitialBalance:        abi_spec.NewTokenAmount(valueSend),
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
		PendingTxns:           pendingTxMap.Root(),
	})

	td.Driver().AssertBalance(multisigAddr, valueSend)
	td.Driver().AssertBalance(outsider, initialBal)
}

func mustProposeMultisigTransfer(td TestDriver, nonce, value int64, txID multisig_spec.TxnID, to, from, proposeTo address.Address, proposeValue int64) {
	msg, err := td.Producer().MultiSigPropose(to, from, nonce, proposeTo, abi_spec.NewTokenAmount(proposeValue), 0, []byte{}, chain.Value(value))
	require.NoError(td.TB(), err)
	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	btxid, err := state.Serialize(&multisig_spec.TxnIDParams{ID: txID})
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode: 0,
		// since the first byte is the cbor type indicator.
		ReturnValue: btxid[1:],
		GasUsed:     big_spec.NewInt(0),
	})
}

func mustCreateMultisigActor(td TestDriver, nonce, value int64, required, unlockDuration int64, ms, creator address.Address, signers ...address.Address) {
	multiSigConstuctParams, err := state.Serialize(&multisig_spec.ConstructorParams{
		Signers:               append(signers, creator),
		NumApprovalsThreshold: required,
		UnlockDuration:        abi_spec.ChainEpoch(unlockDuration),
	})
	require.NoError(td.TB(), err)

	msg, err := td.Producer().InitExec(creator, nonce, builtin_spec.MultisigActorCodeID, multiSigConstuctParams, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: ms.Bytes(),
		GasUsed:     big_spec.NewInt(0),
	})
}

func mustApproveMultisigActor(td TestDriver, nonce, value int64, ms, from address.Address, txID multisig_spec.TxnID) {
	msg, err := td.Producer().MultiSigApprove(ms, from, nonce, txID, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     big_spec.NewInt(0),
	})
}

func mustCancelMultisigActor(td TestDriver, nonce, value int64, ms, from address.Address, txID multisig_spec.TxnID) {
	msg, err := td.Producer().MultiSigCancel(ms, from, nonce, txID, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     big_spec.NewInt(0),
	})
}

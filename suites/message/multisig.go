package message

import (
	"context"
	"testing"

	address "github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	exitcode_spec "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/minio/blake2b-simd"
	"github.com/stretchr/testify/require"
	typegen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestMultiSigActor(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("constructor test", func(t *testing.T) {
		const numApprovals = 3
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(10)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)
		defer td.Complete()

		// creator of the multisig actor
		alice, aliceId := td.NewAccountActor(drivers.SECP, initialBal)

		// expected address of the actor
		multisigAddr := utils.NewIDAddr(t, 1+utils.IdFromAddress(aliceId))

		createRet := td.ComputeInitActorExecReturn(alice, 0, 0, multisigAddr)
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               []address.Address{aliceId},
				NumApprovalsThreshold: numApprovals,
				UnlockDuration:        unlockDuration,
			},
			exitcode_spec.Ok, chain.MustSerialize(&createRet))
	})

	t.Run("propose and cancel", func(t *testing.T) {
		const numApprovals = 2
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(10)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)
		defer td.Complete()

		alice, aliceId := td.NewAccountActor(drivers.SECP, initialBal)

		bob, bobId := td.NewAccountActor(drivers.SECP, initialBal)
		outsider, outsiderId := td.NewAccountActor(drivers.SECP, initialBal)

		multisigAddr := utils.NewIDAddr(t, 1+utils.IdFromAddress(outsiderId))

		createRet := td.ComputeInitActorExecReturn(alice, 0, 0, multisigAddr)
		// create the multisig actor
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               []address.Address{aliceId, bobId},
				NumApprovalsThreshold: numApprovals,
				UnlockDuration:        unlockDuration,
			},
			exitcode_spec.Ok, chain.MustSerialize(&createRet))
		td.AssertBalance(multisigAddr, valueSend)

		// alice proposes that outsider should receive 'valueSend' FIL.
		pparams := multisig_spec.ProposeParams{
			To:     outsider,
			Value:  valueSend,
			Method: builtin_spec.MethodSend,
			Params: []byte{},
		}

		// propose the transaction and assert it exists in the actor state
		txID0 := multisig_spec.TxnID(0)
		expected := typegen.CborInt(txID0)
		td.ApplyExpect(
			td.MessageProducer.MultisigPropose(multisigAddr, alice, &pparams, chain.Nonce(1)),
			chain.MustSerialize(&expected))

		txn0 := multisig_spec.Transaction{
			To:       pparams.To,
			Value:    pparams.Value,
			Method:   pparams.Method,
			Params:   pparams.Params,
			Approved: []address.Address{aliceId},
		}
		ph := makeProposalHash(t, &txn0)
		td.AssertMultisigTransaction(multisigAddr, txID0, txn0)

		// bob cancels alice's transaction. This fails as bob did not create alice's transaction.
		td.ApplyFailure(
			td.MessageProducer.MultisigCancel(multisigAddr, bob, &multisig_spec.TxnIDParams{ID: txID0, ProposalHash: ph}, chain.Nonce(0)),
			exitcode_spec.ErrForbidden)

		// alice cancels their transaction. The outsider doesn't receive any FIL, the multisig actor's balance is empty, and the
		// transaction is canceled.
		td.ApplyOk(
			td.MessageProducer.MultisigCancel(multisigAddr, alice, &multisig_spec.TxnIDParams{ID: txID0, ProposalHash: ph}, chain.Nonce(2)),
		)
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{aliceId, bobId},
			NumApprovalsThreshold: numApprovals,
			NextTxnID:             1,
			InitialBalance:        valueSend,
			StartEpoch:            1,
			UnlockDuration:        unlockDuration,
		})
		td.AssertBalance(multisigAddr, valueSend)
	})

	t.Run("propose and approve", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		var initialBal = abi_spec.NewTokenAmount(200000000000)
		const numApprovals = 2
		const unlockDuration = 1
		var valueSend = abi_spec.NewTokenAmount(10)

		// Signers
		alice, aliceId := td.NewAccountActor(drivers.SECP, initialBal)
		bob, bobId := td.NewAccountActor(drivers.SECP, initialBal)

		// Not Signer
		outsider, outsiderId := td.NewAccountActor(drivers.SECP, initialBal)

		// Multisig actor address
		multisigAddr := utils.NewIDAddr(t, 1+utils.IdFromAddress(outsiderId))
		createRet := td.ComputeInitActorExecReturn(alice, 0, 0, multisigAddr)

		// create the multisig actor
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               []address.Address{aliceId, bobId},
				NumApprovalsThreshold: numApprovals,
				UnlockDuration:        unlockDuration,
			},
			exitcode_spec.Ok, chain.MustSerialize(&createRet))

		// setup propose expected values and params
		pparams := multisig_spec.ProposeParams{
			To:     outsider,
			Value:  valueSend,
			Method: builtin_spec.MethodSend,
			Params: []byte{},
		}

		// propose the transaction and assert it exists in the actor state
		txID0 := multisig_spec.TxnID(0)
		expected := typegen.CborInt(txID0)
		td.ApplyExpect(
			td.MessageProducer.MultisigPropose(multisigAddr, alice, &pparams, chain.Nonce(1)),
			chain.MustSerialize(&expected))

		txn0 := multisig_spec.Transaction{
			To:       pparams.To,
			Value:    pparams.Value,
			Method:   pparams.Method,
			Params:   pparams.Params,
			Approved: []address.Address{aliceId},
		}
		ph := makeProposalHash(t, &txn0)
		td.AssertMultisigTransaction(multisigAddr, txID0, txn0)

		// outsider proposes themselves to receive 'valueSend' FIL. This fails as they are not a signer.
		td.ApplyFailure(
			td.MessageProducer.MultisigPropose(multisigAddr, outsider, &pparams, chain.Nonce(0)),
			exitcode_spec.ErrForbidden)

		// outsider approves the value transfer alice sent. This fails as they are not a signer.
		td.ApplyFailure(
			td.MessageProducer.MultisigApprove(multisigAddr, outsider, &multisig_spec.TxnIDParams{ID: txID0, ProposalHash: ph}, chain.Nonce(1)),
			exitcode_spec.ErrForbidden)

		// increment the epoch to unlock the funds
		td.ExeCtx.Epoch += unlockDuration
		balanceBefore := td.GetBalance(outsider)

		// bob approves transfer of 'valueSend' FIL to outsider.
		td.ApplyOk(
			td.MessageProducer.MultisigApprove(multisigAddr, bob, &multisig_spec.TxnIDParams{ID: txID0, ProposalHash: ph}, chain.Nonce(0)),
		)
		txID1 := multisig_spec.TxnID(1)
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{aliceId, bobId},
			NumApprovalsThreshold: numApprovals,
			NextTxnID:             txID1,
			InitialBalance:        valueSend,
			StartEpoch:            1,
			UnlockDuration:        unlockDuration,
		})
		td.AssertMultisigContainsTransaction(multisigAddr, txID0, false)
		// Multisig balance has been transferred to outsider.
		td.AssertBalance(multisigAddr, big_spec.Zero())
		td.AssertBalance(outsider, big_spec.Add(balanceBefore, valueSend))
	})

	t.Run("add signer", func(t *testing.T) {
		const initialNumApprovals = 1
		var msValue = abi_spec.NewTokenAmount(100000000000)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)
		defer td.Complete()

		alice, aliceId := td.NewAccountActor(drivers.SECP, initialBal) // 101
		_, bobId := td.NewAccountActor(drivers.SECP, initialBal)       // 102
		var initialSigners = []address.Address{aliceId}

		multisigAddr := utils.NewIDAddr(t, 1+utils.IdFromAddress(bobId))
		createRet := td.ComputeInitActorExecReturn(alice, 0, 0, multisigAddr)

		td.MustCreateAndVerifyMultisigActor(0, msValue, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               initialSigners,
				NumApprovalsThreshold: initialNumApprovals,
				UnlockDuration:        0,
			},
			exitcode_spec.Ok,
			chain.MustSerialize(&createRet),
		)

		addSignerParams := multisig_spec.AddSignerParams{
			Signer:   bobId,
			Increase: false,
		}

		// alice fails to call directly since AddSigner
		td.ApplyFailure(
			td.MessageProducer.MultisigAddSigner(multisigAddr, alice, &addSignerParams, chain.Nonce(1)),
			exitcode_spec.SysErrForbidden,
		)

		// AddSigner must be staged through the multisig itself
		// Alice proposes the AddSigner.
		// Since approvals = 1 this auto-approves the transaction.
		txID0 := multisig_spec.TxnID(0)
		expected := typegen.CborInt(txID0)
		td.ApplyExpect(
			td.MessageProducer.MultisigPropose(multisigAddr, alice, &multisig_spec.ProposeParams{
				To:     multisigAddr,
				Value:  big_spec.Zero(),
				Method: builtin_spec.MethodsMultisig.AddSigner,
				Params: chain.MustSerialize(&addSignerParams),
			}, chain.Nonce(2)),
			chain.MustSerialize(&expected),
		)

		// TODO also exercise the approvals = 2 case with explicit approval.

		// Check that bob is now a signer
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               append(initialSigners, bobId),
			NumApprovalsThreshold: initialNumApprovals,
			NextTxnID:             multisig_spec.TxnID(1),
			InitialBalance:        big_spec.Zero(),
			StartEpoch:            0,
			UnlockDuration:        0,
		})
	})
}

func makeProposalHash(t *testing.T, txn *multisig_spec.Transaction) []byte {
	txnHash, err := multisig_spec.ComputeProposalHash(txn, blake2b.Sum256)
	require.NoError(t, err)
	return txnHash
}

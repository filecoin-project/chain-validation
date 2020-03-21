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

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
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
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: chain.MustSerialize(&createRet),
				GasUsed:     big_spec.Zero(), // Ignored
			})
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
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: chain.MustSerialize(&createRet),
				GasUsed:     big_spec.Zero(), // Ignored
			})
		td.AssertBalance(multisigAddr, valueSend)

		// alice proposes that outsider should receive 'valueSend' FIL.
		txID0 := multisig_spec.TxnID(0)
		pparams := multisig_spec.ProposeParams{
			To:     outsider,
			Value:  valueSend,
			Method: builtin_spec.MethodSend,
			Params: []byte{},
		}

		// propose the transaction and assert it exists in the actor state
		// TODO: stop using TxnIDParams since it's not being passed as params to any method.
		// Just get the right serialized bytes for the return value.
		btxid := chain.MustSerialize(&multisig_spec.TxnIDParams{ID: txID0})
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigPropose(multisigAddr, alice, pparams, chain.Value(big_spec.Zero()), chain.Nonce(1)),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: btxid[1:], GasUsed: big_spec.Zero()}, // Gas ignored
		)
		td.AssertMultisigTransaction(multisigAddr, txID0, multisig_spec.Transaction{
			To:       pparams.To,
			Value:    pparams.Value,
			Method:   pparams.Method,
			Params:   pparams.Params,
			Approved: []address.Address{aliceId},
		})

		// bob cancels alice's transaction. This fails as bob did not create alice's transaction.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigCancel(multisigAddr, bob, multisig_spec.TxnIDParams{ID: txID0}, chain.Value(big_spec.Zero()), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode_spec.ErrForbidden, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.Zero()}, // Gas ignored
		)

		// alice cancels their transaction. The outsider doesn't receive any FIL, the multisig actor's balance is empty, and the
		// transaction is canceled.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigCancel(multisigAddr, alice, multisig_spec.TxnIDParams{ID: txID0}, chain.Nonce(2), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.Zero()}, // Gas ignored
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
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: chain.MustSerialize(&createRet),
				GasUsed:     big_spec.NewInt(1542),
			})

		// setup propose expected values and params
		txID0 := multisig_spec.TxnID(0)
		pparams := multisig_spec.ProposeParams{
			To:     outsider,
			Value:  valueSend,
			Method: builtin_spec.MethodSend,
			Params: []byte{},
		}

		// propose the transaction and assert it exists in the actor state
		btxid := chain.MustSerialize(&multisig_spec.TxnIDParams{ID: txID0})
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigPropose(multisigAddr, alice, pparams, chain.Value(big_spec.Zero()), chain.Nonce(1)),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: btxid[1:], GasUsed: big_spec.Zero()}, // Gas ignored
		)
		td.AssertMultisigTransaction(multisigAddr, txID0, multisig_spec.Transaction{
			To:       pparams.To,
			Value:    pparams.Value,
			Method:   pparams.Method,
			Params:   pparams.Params,
			Approved: []address.Address{aliceId},
		})

		// outsider proposes themselves to receive 'valueSend' FIL. This fails as they are not a signer.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigPropose(multisigAddr, outsider, multisig_spec.ProposeParams{
				To:     outsider,
				Value:  valueSend,
				Method: builtin_spec.MethodSend,
				Params: []byte{},
			}, chain.Value(big_spec.Zero()), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode_spec.ErrForbidden, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.Zero()}, // Gas ignored
		)

		// outsider approves the value transfer alice sent. This fails as they are not a signer.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigApprove(multisigAddr, outsider, multisig_spec.TxnIDParams{ID: txID0}, chain.Value(big_spec.Zero()), chain.Nonce(1)),
			types.MessageReceipt{ExitCode: exitcode_spec.ErrForbidden, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.Zero()}, // Gas ignored
		)

		// increment the epoch to unlock the funds
		td.ExeCtx.Epoch += unlockDuration
		balanceBefore := td.GetBalance(outsider)

		// bob approves transfer of 'valueSend' FIL to outsider.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigApprove(multisigAddr, bob, multisig_spec.TxnIDParams{ID: txID0}, chain.Value(big_spec.Zero()), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.Zero()}, // Gas ignored
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
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: chain.MustSerialize(&createRet),
				GasUsed:     big_spec.Zero(), // Gas ignored
			})

		addSignerParams := multisig_spec.AddSignerParams{
			Signer:   bobId,
			Increase: false,
		}

		// alice fails to call directly since AddSigner
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigAddSigner(multisigAddr, alice, addSignerParams, chain.Nonce(1)),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.SysErrForbidden,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.Zero(), // Gas ignored
			})

		// AddSigner must be staged through the multisig itself
		txID0 := multisig_spec.TxnID(0)

		// Alice proposes the AddSigner.
		// Since approvals = 1 this auto-approves the transaction.
		btxid := chain.MustSerialize(&multisig_spec.TxnIDParams{ID: txID0})
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigPropose(multisigAddr, alice, multisig_spec.ProposeParams{
				To:     multisigAddr,
				Value:  big_spec.Zero(),
				Method: builtin_spec.MethodsMultisig.AddSigner,
				Params: chain.MustSerialize(&addSignerParams),
			}, chain.Nonce(2)),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: btxid[1:], GasUsed: big_spec.Zero()}, // Gas ignored
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

package suites

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
		WithActorState([]drivers.ActorState{
			drivers.DefaultInitActorState,
			drivers.DefaultRewardActorState,
			drivers.DefaultBurntFundsActorState,
		})

	t.Run("constructor test", func(t *testing.T) {
		const numApprovals = 3
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(10)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)

		// creator of the multisig actor
		alice, aliceId := td.NewAccountActor(drivers.SECP, initialBal)

		// expected address of the actor
		multisigAddr := utils.NewIDAddr(t, 101)

		createRet := td.ComputeInitActorExecReturn(aliceId, 0, multisigAddr)
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               []address.Address{alice},
				NumApprovalsThreshold: numApprovals,
				UnlockDuration:        unlockDuration,
			},
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: chain.MustSerialize(&createRet),
				GasUsed:     big_spec.NewInt(1125),
			})
	})

	t.Run("propose and cancel", func(t *testing.T) {
		const numApprovals = 2
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(10)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)

		alice, aliceId := td.NewAccountActor(drivers.SECP, initialBal)

		bob, _ := td.NewAccountActor(drivers.SECP, initialBal)
		outsider, _ := td.NewAccountActor(drivers.SECP, initialBal)

		multisigAddr := utils.NewIDAddr(t, 103)

		createRet := td.ComputeInitActorExecReturn(aliceId, 0, multisigAddr)
		// create the multisig actor
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               []address.Address{alice, bob},
				NumApprovalsThreshold: numApprovals,
				UnlockDuration:        unlockDuration,
			},
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: chain.MustSerialize(&createRet),
				GasUsed:     big_spec.NewInt(1387),
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
		btxid := chain.MustSerialize(&multisig_spec.TxnIDParams{ID: txID0})
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigPropose(multisigAddr, alice, pparams, chain.Value(big_spec.Zero()), chain.Nonce(1)),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: btxid[1:], GasUsed: big_spec.NewInt(1280)},
		)
		td.AssertMultisigTransaction(multisigAddr, txID0, multisig_spec.Transaction{
			To:       pparams.To,
			Value:    pparams.Value,
			Method:   pparams.Method,
			Params:   pparams.Params,
			Approved: []address.Address{alice},
		})

		// bob cancels alice's transaction. This fails as bob did not create alice's transaction.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigCancel(multisigAddr, bob, multisig_spec.TxnIDParams{ID: txID0}, chain.Value(big_spec.Zero()), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode_spec.ErrForbidden, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.NewInt(1000000)},
		)

		// alice cancels their transaction. The outsider doesn't receive any FIL, the multisig actor's balance is empty, and the
		// transaction is canceled.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigCancel(multisigAddr, alice, multisig_spec.TxnIDParams{ID: txID0}, chain.Nonce(1), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.NewInt(639)},
		)
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, bob},
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
		var initialBal = abi_spec.NewTokenAmount(200000000000)
		const numApprovals = 2
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(10)

		// Signers
		alice, _ := td.NewAccountActor(drivers.SECP, initialBal)
		bob, _ := td.NewAccountActor(drivers.SECP, initialBal)

		// Not Signer
		outsider, _ := td.NewAccountActor(drivers.SECP, initialBal)

		// Multisig actor address
		multisigAddr := utils.NewIDAddr(t, 104)

		// create the multisig actor
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               []address.Address{alice, bob},
				NumApprovalsThreshold: numApprovals,
				UnlockDuration:        unlockDuration,
			},
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: multisigAddr.Bytes(),
				GasUsed:     big_spec.NewInt(1387),
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
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: btxid[1:], GasUsed: big_spec.NewInt(1280)},
		)
		td.AssertMultisigTransaction(multisigAddr, txID0, multisig_spec.Transaction{
			To:       pparams.To,
			Value:    pparams.Value,
			Method:   pparams.Method,
			Params:   pparams.Params,
			Approved: []address.Address{alice},
		})

		// outsider proposes themselves to receive 'valueSend' FIL. This fails as they are not a signer.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigPropose(multisigAddr, outsider, multisig_spec.ProposeParams{
				To:     outsider,
				Value:  valueSend,
				Method: builtin_spec.MethodSend,
				Params: []byte{},
			}, chain.Value(big_spec.Zero()), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode_spec.ErrForbidden, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.NewInt(1000000)},
		)

		// outsider approves the value transfer alice sent. This fails as they are not a signer.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigApprove(multisigAddr, outsider, multisig_spec.TxnIDParams{ID: txID0}, chain.Value(big_spec.Zero()), chain.Nonce(1)),
			types.MessageReceipt{ExitCode: exitcode_spec.ErrForbidden, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.NewInt(1000000)},
		)

		// bob approves transfer of 'valueSend' FIL to outsider.
		txID1 := multisig_spec.TxnID(1)
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigApprove(multisigAddr, bob, multisig_spec.TxnIDParams{ID: txID1}, chain.Value(big_spec.Zero()), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.NewInt(1691)},
		)
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, bob},
			NumApprovalsThreshold: numApprovals,
			NextTxnID:             txID1,
			InitialBalance:        valueSend,
			StartEpoch:            1,
			UnlockDuration:        unlockDuration,
		})
		td.AssertMultisigContainsTransaction(multisigAddr, txID0, false)
	})

	t.Run("add signer and increase threshold", func(t *testing.T) {
		const initialNumApprovals = 2
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(100000000000)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)

		alice, _ := td.NewAccountActor(drivers.SECP, initialBal) // 101
		bob, _ := td.NewAccountActor(drivers.SECP, initialBal)   // 102
		chuck, _ := td.NewAccountActor(drivers.SECP, initialBal) // 103
		duck, _ := td.NewAccountActor(drivers.SECP, initialBal)  // 104
		var initialSigners = []address.Address{alice, bob}

		multisigAddr := utils.NewIDAddr(t, 105)

		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               initialSigners,
				NumApprovalsThreshold: initialNumApprovals,
				UnlockDuration:        unlockDuration,
			},
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: multisigAddr.Bytes(),
				GasUsed:     big_spec.NewInt(1508),
			})

		// alice fails to add a singer since this method can only be called by the multisig actors wallet address
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigAddSigner(multisigAddr, alice, multisig_spec.AddSignerParams{
				Signer:   chuck,
				Increase: false,
			}, chain.Value(big_spec.Zero()), chain.Nonce(1)),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.SysErrActorNotFound, // TODO set the correct error code here, lotus returns 'SysErrActorNotFound`, which is probably wrong.
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(1_000_000),
			})

		// success when multisig actor calls the add signer method
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigAddSigner(multisigAddr, multisigAddr, multisig_spec.AddSignerParams{
				Signer:   chuck,
				Increase: false,
			}, chain.Value(big_spec.Zero()), chain.Nonce(0)),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(517),
			})
		// assert that chuck is now a signer
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               append(initialSigners, chuck),
			NumApprovalsThreshold: initialNumApprovals,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})

		// add another signer and increase the number of signers required
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigAddSigner(multisigAddr, multisigAddr, multisig_spec.AddSignerParams{
				Signer:   duck,
				Increase: true,
			}, chain.Value(big_spec.Zero()), chain.Nonce(1)),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(583),
			})
		// assert that duck is noe a signer and the number of approvals required increased
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               append(initialSigners, chuck, duck),
			NumApprovalsThreshold: initialNumApprovals + 1,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})

	})

	t.Run("remove signer and decreases threshold", func(t *testing.T) {
		const initialNumApprovals = 2
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(100000000000)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)

		alice, _ := td.NewAccountActor(drivers.SECP, initialBal) // 101
		bob, _ := td.NewAccountActor(drivers.SECP, initialBal)   // 102
		chuck, _ := td.NewAccountActor(drivers.SECP, initialBal) // 103
		duck, _ := td.NewAccountActor(drivers.SECP, initialBal)  // 104
		var initialSigners = []address.Address{alice, bob, chuck, duck}

		multisigAddr := utils.NewIDAddr(t, 105)

		// create a ms actor with 4 signers and 3 approvals required
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               initialSigners,
				NumApprovalsThreshold: initialNumApprovals,
				UnlockDuration:        unlockDuration,
			},
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: multisigAddr.Bytes(),
				GasUsed:     big_spec.NewInt(1684),
			})

		// alice fails to remove a singer since this method can only be called by the multisig actors wallet address
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigRemoveSigner(multisigAddr, alice, multisig_spec.RemoveSignerParams{
				Signer:   chuck,
				Decrease: false,
			}, chain.Value(big_spec.Zero()), chain.Nonce(1)),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.SysErrActorNotFound, // TODO set the correct error code here, lotus returns 'SysErrActorNotFound`, which is probably wrong.
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(1_000_000),
			})

		// success when multisig actor calls the remove signer method and removes duck
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigRemoveSigner(multisigAddr, multisigAddr, multisig_spec.RemoveSignerParams{
				Signer:   duck,
				Decrease: false,
			}, chain.Value(big_spec.Zero()), chain.Nonce(0)),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(561),
			})
		// assert that duck is no longer a signer and that the number of required approvals has remained unchanged.
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, bob, chuck},
			NumApprovalsThreshold: initialNumApprovals,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})

		// remove chuck and decrease the number of signers required
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigRemoveSigner(multisigAddr, multisigAddr, multisig_spec.RemoveSignerParams{
				Signer:   chuck,
				Decrease: true,
			}, chain.Value(big_spec.Zero()), chain.Nonce(1)),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(495),
			})
		// assert that duck is no a signer and the number of approvals required decreased.
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, bob},
			NumApprovalsThreshold: initialNumApprovals - 1,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})

	})

	t.Run("swap signers and change number of approvals", func(t *testing.T) {
		const initialNumApprovals = 2
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(100000000000)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)

		alice, _ := td.NewAccountActor(drivers.SECP, initialBal) // 101
		bob, _ := td.NewAccountActor(drivers.SECP, initialBal)   // 102
		// chuck will be swapped in below
		chuck, _ := td.NewAccountActor(drivers.SECP, initialBal) // 103

		var initialSigners = []address.Address{alice, bob}

		multisigAddr := utils.NewIDAddr(t, 104)

		// create a ms actor with 4 signers and 3 approvals required
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               initialSigners,
				NumApprovalsThreshold: initialNumApprovals,
				UnlockDuration:        unlockDuration,
			},
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: multisigAddr.Bytes(),
				GasUsed:     big_spec.NewInt(1403),
			})

		// create parameters to swap bob for chuck
		swapParams := multisig_spec.SwapSignerParams{
			From: bob,
			To:   chuck,
		}
		// alice fails to since they are not the multisig address.
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigSwapSigner(multisigAddr, alice, swapParams, chain.Nonce(1), chain.Value(big_spec.Zero())),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.SysErrActorNotFound, // TODO set the correct error code here, lotus returns 'SysErrActorNotFound`, which is probably wrong.
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(1_000_000),
			})

		// swap operation success
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigSwapSigner(multisigAddr, multisigAddr, swapParams, chain.Nonce(0), chain.Value(big_spec.Zero())),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(515),
			})
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, chuck},
			NumApprovalsThreshold: initialNumApprovals,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})

		// decrease the threshold and assert state change
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.MultisigChangeNumApprovalsThreshold(multisigAddr, multisigAddr, multisig_spec.ChangeNumApprovalsThresholdParams{NewThreshold: initialNumApprovals - 1}, chain.Nonce(1), chain.Value(big_spec.Zero())),
			types.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(427),
			})
		td.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, chuck},
			NumApprovalsThreshold: initialNumApprovals - 1,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})
	})
}

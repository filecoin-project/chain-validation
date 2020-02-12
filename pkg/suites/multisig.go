package suites

import (
	"context"
	"testing"

	address "github.com/filecoin-project/go-address"
	require "github.com/stretchr/testify/require"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	exitcode_spec "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	chain "github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
)

func TestMultiSigActor(t *testing.T, factory Factories) {
	defaultMiner, err := address.NewSecp256k1Address([]byte{'m', 'i', 'n', 'e', 'r'})
	require.NoError(t, err)

	builder := NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(big_spec.NewInt(1000000)).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithSingletonActors(map[address.Address]big_spec.Int{
			builtin_spec.InitActorAddr:         big_spec.NewInt(0),
			builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
			builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
			builtin_spec.RewardActorAddr:       TotalNetworkBalance,
		}).
		WithDefaultMiner(defaultMiner)

	t.Run("constructor test", func(t *testing.T) {
		const numApprovals = 3
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(10)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)

		// creator of the multisig actor
		alice := td.Driver.NewAccountActor(SECP, initialBal)

		// expected address of the actor
		multisigAddr, err := address.NewIDAddress(102)
		require.NoError(t, err)

		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               []address.Address{alice},
				NumApprovalsThreshold: numApprovals,
				UnlockDuration:        unlockDuration,
			},
			chain.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: multisigAddr.Bytes(),
				GasUsed:     big_spec.NewInt(1125),
			})
	})

	t.Run("propose and cancel", func(t *testing.T) {
		const numApprovals = 2
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(10)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)

		alice := td.Driver.NewAccountActor(SECP, initialBal)
		bob := td.Driver.NewAccountActor(SECP, initialBal)
		outsider := td.Driver.NewAccountActor(SECP, initialBal)

		multisigAddr, err := address.NewIDAddress(104)
		require.NoError(t, err)

		// create the multisig actor
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               []address.Address{alice, bob},
				NumApprovalsThreshold: numApprovals,
				UnlockDuration:        unlockDuration,
			},
			chain.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: multisigAddr.Bytes(),
				GasUsed:     big_spec.NewInt(1387),
			})
		td.Driver.AssertBalance(multisigAddr, valueSend)

		// alice proposes that outsider should receive 'valueSend' FIL.
		txID0 := multisig_spec.TxnID(0)
		pparams := multisig_spec.ProposeParams{
			To:     outsider,
			Value:  valueSend,
			Method: builtin_spec.MethodSend,
			Params: []byte{},
		}

		// propose the transaction and assert it exists in the actor state
		btxid, err := state.Serialize(&multisig_spec.TxnIDParams{ID: txID0})
		require.NoError(t, err)
		td.MustProposeMultisigTransfer(1, big_spec.Zero(), txID0, multisigAddr, alice, pparams, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: btxid[1:],
			GasUsed:     big_spec.NewInt(1280),
		})
		td.Driver.AssertMultisigTransaction(multisigAddr, txID0, multisig_spec.Transaction{
			To:       pparams.To,
			Value:    pparams.Value,
			Method:   pparams.Method,
			Params:   pparams.Params,
			Approved: []address.Address{alice},
		})

		// bob cancels alice's transaction. This fails as bob did not create alice's transaction.
		td.ApplyMessageExpectReceipt(
			func() (*chain.Message, error) {
				return td.Producer.MultiSigCancel(multisigAddr, bob, txID0, chain.Value(big_spec.Zero()), chain.Nonce(0))
			},
			chain.MessageReceipt{
				ExitCode:    exitcode_spec.ErrForbidden,
				ReturnValue: nil,
				GasUsed:     big_spec.NewInt(1000000),
			},
		)

		// alice cancels their transaction. The outsider doesn't receive any FIL, the multisig actor's balance is empty, and the
		// transaction is canceled.
		td.MustCancelMultisigActor(2, big_spec.Zero(), multisigAddr, alice, txID0, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: EmptyRetrunValueBytes,
			GasUsed:     big_spec.NewInt(639),
		})
		td.Driver.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, bob},
			NumApprovalsThreshold: numApprovals,
			NextTxnID:             1,
			InitialBalance:        valueSend,
			StartEpoch:            1,
			UnlockDuration:        unlockDuration,
		})
		td.Driver.AssertBalance(multisigAddr, valueSend)
	})

	t.Run("propose and approve", func(t *testing.T) {
		td := builder.Build(t)
		var initialBal = abi_spec.NewTokenAmount(200000000000)
		const numApprovals = 2
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(10)

		// Signers
		alice := td.Driver.NewAccountActor(SECP, initialBal)
		bob := td.Driver.NewAccountActor(SECP, initialBal)

		// Not Signer
		outsider := td.Driver.NewAccountActor(SECP, initialBal)

		// Multisig actor address
		multisigAddr, err := address.NewIDAddress(104)
		require.NoError(t, err)

		// create the multisig actor
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               []address.Address{alice, bob},
				NumApprovalsThreshold: numApprovals,
				UnlockDuration:        unlockDuration,
			},
			chain.MessageReceipt{
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
		btxid, err := state.Serialize(&multisig_spec.TxnIDParams{ID: txID0})
		require.NoError(t, err)
		td.MustProposeMultisigTransfer(1, big_spec.Zero(), txID0, multisigAddr, alice, pparams, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: btxid[1:],
			GasUsed:     big_spec.NewInt(1280),
		})
		td.Driver.AssertMultisigTransaction(multisigAddr, txID0, multisig_spec.Transaction{
			To:       pparams.To,
			Value:    pparams.Value,
			Method:   pparams.Method,
			Params:   pparams.Params,
			Approved: []address.Address{alice},
		})

		// outsider proposes themselves to receive 'valueSend' FIL. This fails as they are not a signer.
		td.ApplyMessageExpectReceipt(
			func() (*chain.Message, error) {
				return td.Producer.MultiSigPropose(multisigAddr, outsider, multisig_spec.ProposeParams{
					To:     outsider,
					Value:  valueSend,
					Method: builtin_spec.MethodSend,
					Params: []byte{},
				}, chain.Value(big_spec.Zero()), chain.Nonce(0))
			},
			chain.MessageReceipt{
				ExitCode:    exitcode_spec.ErrForbidden,
				ReturnValue: nil,
				GasUsed:     big_spec.NewInt(1000000),
			},
		)

		// outsider approves the value transfer alice sent. This fails as they are not a signer.
		td.ApplyMessageExpectReceipt(
			func() (*chain.Message, error) {
				return td.Producer.MultiSigApprove(multisigAddr, outsider, txID0, chain.Value(big_spec.Zero()), chain.Nonce(1))
			},
			chain.MessageReceipt{
				ExitCode:    exitcode_spec.ErrForbidden,
				ReturnValue: nil,
				GasUsed:     big_spec.NewInt(1000000),
			},
		)

		// bob approves transfer of 'valueSend' FIL to outsider.
		txID1 := multisig_spec.TxnID(1)
		td.MustApproveMultisigActor(0, big_spec.Zero(), multisigAddr, bob, txID0, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: EmptyRetrunValueBytes,
			GasUsed:     big_spec.NewInt(1691),
		})
		td.Driver.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, bob},
			NumApprovalsThreshold: numApprovals,
			NextTxnID:             txID1,
			InitialBalance:        valueSend,
			StartEpoch:            1,
			UnlockDuration:        unlockDuration,
		})
		td.Driver.AssertMultisigContainsTransaction(multisigAddr, txID0, false)
	})

	t.Run("add signer and increase threshold", func(t *testing.T) {
		const initialNumApprovals = 2
		const unlockDuration = 10
		var valueSend = abi_spec.NewTokenAmount(100000000000)
		var initialBal = abi_spec.NewTokenAmount(200000000000)

		td := builder.Build(t)

		alice := td.Driver.NewAccountActor(SECP, initialBal) // 101
		bob := td.Driver.NewAccountActor(SECP, initialBal)   // 102
		chuck := td.Driver.NewAccountActor(SECP, initialBal) // 103
		duck := td.Driver.NewAccountActor(SECP, initialBal)  // 104
		var initialSigners = []address.Address{alice, bob}

		multisigAddr, err := address.NewIDAddress(105) // 105
		require.NoError(t, err)

		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               initialSigners,
				NumApprovalsThreshold: initialNumApprovals,
				UnlockDuration:        unlockDuration,
			},
			chain.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: multisigAddr.Bytes(),
				GasUsed:     big_spec.NewInt(1508),
			})

		// alice fails to add a singer since this method can only be called by the multisig actors wallet address
		msg, err := td.Producer.MultiSigAddSigner(multisigAddr, alice, multisig_spec.AddSignerParams{
			Signer:   chuck,
			Increase: false,
		}, chain.Value(big_spec.Zero()), chain.Nonce(1))
		require.NoError(t, err)

		msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    exitcode_spec.SysErrActorNotFound, // TODO set the correct error code here, lotus returns 'SysErrActorNotFound`, which is probably wrong.
			ReturnValue: nil,
			GasUsed:     big_spec.NewInt(1_000_000),
		})

		// success when multisig actor calls the add signer method
		msg, err = td.Producer.MultiSigAddSigner(multisigAddr, multisigAddr, multisig_spec.AddSignerParams{
			Signer:   chuck,
			Increase: false,
		}, chain.Value(big_spec.Zero()), chain.Nonce(0))
		require.NoError(t, err)

		msgReceipt, err = td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: EmptyRetrunValueBytes,
			GasUsed:     big_spec.NewInt(517),
		})
		// assert that chuck is now a signer
		td.Driver.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               append(initialSigners, chuck),
			NumApprovalsThreshold: initialNumApprovals,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})

		// add another signer and increase the number of signers required
		msg, err = td.Producer.MultiSigAddSigner(multisigAddr, multisigAddr, multisig_spec.AddSignerParams{
			Signer:   duck,
			Increase: true,
		}, chain.Value(big_spec.Zero()), chain.Nonce(1))
		require.NoError(t, err)

		msgReceipt, err = td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: EmptyRetrunValueBytes,
			GasUsed:     big_spec.NewInt(583),
		})
		// assert that duck is noe a signer and the number of approvals required increased
		td.Driver.AssertMultisigState(multisigAddr, multisig_spec.State{
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

		alice := td.Driver.NewAccountActor(SECP, initialBal) // 101
		bob := td.Driver.NewAccountActor(SECP, initialBal)   // 102
		chuck := td.Driver.NewAccountActor(SECP, initialBal) // 103
		duck := td.Driver.NewAccountActor(SECP, initialBal)  // 104
		var initialSigners = []address.Address{alice, bob, chuck, duck}

		multisigAddr, err := address.NewIDAddress(105) // 105
		require.NoError(t, err)

		// create a ms actor with 4 signers and 3 approvals required
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               initialSigners,
				NumApprovalsThreshold: initialNumApprovals,
				UnlockDuration:        unlockDuration,
			},
			chain.MessageReceipt{
				ExitCode:    exitcode_spec.Ok,
				ReturnValue: multisigAddr.Bytes(),
				GasUsed:     big_spec.NewInt(1684),
			})

		// alice fails to remove a singer since this method can only be called by the multisig actors wallet address
		msg, err := td.Producer.MultiSigRemoveSigner(multisigAddr, alice, multisig_spec.RemoveSignerParams{
			Signer:   chuck,
			Decrease: false,
		}, chain.Value(big_spec.Zero()), chain.Nonce(1))
		require.NoError(t, err)

		msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    exitcode_spec.SysErrActorNotFound, // TODO set the correct error code here, lotus returns 'SysErrActorNotFound`, which is probably wrong.
			ReturnValue: nil,
			GasUsed:     big_spec.NewInt(1_000_000),
		})

		// success when multisig actor calls the remove signer method and removes duck
		msg, err = td.Producer.MultiSigRemoveSigner(multisigAddr, multisigAddr, multisig_spec.RemoveSignerParams{
			Signer:   duck,
			Decrease: false,
		}, chain.Value(big_spec.Zero()), chain.Nonce(0))
		require.NoError(t, err)

		msgReceipt, err = td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: EmptyRetrunValueBytes,
			GasUsed:     big_spec.NewInt(561),
		})
		// assert that duck is no longer a signer and that the number of required approvals has remained unchanged.
		td.Driver.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, bob, chuck},
			NumApprovalsThreshold: initialNumApprovals,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})

		// remove chuck and decrease the number of signers required
		msg, err = td.Producer.MultiSigRemoveSigner(multisigAddr, multisigAddr, multisig_spec.RemoveSignerParams{
			Signer:   chuck,
			Decrease: true,
		}, chain.Value(big_spec.Zero()), chain.Nonce(1))
		require.NoError(t, err)

		msgReceipt, err = td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: EmptyRetrunValueBytes,
			GasUsed:     big_spec.NewInt(495),
		})
		// assert that duck is no a signer and the number of approvals required decreased.
		td.Driver.AssertMultisigState(multisigAddr, multisig_spec.State{
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

		alice := td.Driver.NewAccountActor(SECP, initialBal) // 101
		bob := td.Driver.NewAccountActor(SECP, initialBal)   // 102
		// chuck will be swapped in below
		chuck := td.Driver.NewAccountActor(SECP, initialBal) // 103

		var initialSigners = []address.Address{alice, bob}

		multisigAddr, err := address.NewIDAddress(104) // 104
		require.NoError(t, err)

		// create a ms actor with 4 signers and 3 approvals required
		td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice,
			&multisig_spec.ConstructorParams{
				Signers:               initialSigners,
				NumApprovalsThreshold: initialNumApprovals,
				UnlockDuration:        unlockDuration,
			},
			chain.MessageReceipt{
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
		msg, err := td.Producer.MultiSigSwapSigner(multisigAddr, alice, swapParams, chain.Nonce(1), chain.Value(big_spec.Zero()))
		require.NoError(t, err)

		msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    exitcode_spec.SysErrActorNotFound, // TODO set the correct error code here, lotus returns 'SysErrActorNotFound`, which is probably wrong.
			ReturnValue: nil,
			GasUsed:     big_spec.NewInt(1_000_000),
		})

		// swap operation success
		msg, err = td.Producer.MultiSigSwapSigner(multisigAddr, multisigAddr, swapParams, chain.Nonce(0), chain.Value(big_spec.Zero()))
		require.NoError(t, err)

		msgReceipt, err = td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: EmptyRetrunValueBytes,
			GasUsed:     big_spec.NewInt(515),
		})
		td.Driver.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, chuck},
			NumApprovalsThreshold: initialNumApprovals,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})

		// decrease the threshold and assert state change
		msg, err = td.Producer.MultiSigChangeApprovalsThreshold(multisigAddr, multisigAddr, initialNumApprovals-1, chain.Nonce(1), chain.Value(big_spec.Zero()))
		require.NoError(t, err)

		msgReceipt, err = td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    exitcode_spec.Ok,
			ReturnValue: EmptyRetrunValueBytes,
			GasUsed:     big_spec.NewInt(427),
		})
		td.Driver.AssertMultisigState(multisigAddr, multisig_spec.State{
			Signers:               []address.Address{alice, chuck},
			NumApprovalsThreshold: initialNumApprovals - 1,
			NextTxnID:             0,
			InitialBalance:        valueSend,
			StartEpoch:            td.ExeCtx.Epoch,
			UnlockDuration:        unlockDuration,
		})
	})

}

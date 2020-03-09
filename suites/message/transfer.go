package message

import (
	"context"
	"testing"

	address "github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	require "github.com/stretchr/testify/require"

	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	account_spec "github.com/filecoin-project/specs-actors/actors/builtin/account"

	chain "github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

type valueTransferTestCases struct {
	desc string

	sender    address.Address
	senderBal big_spec.Int

	transferAmnt big_spec.Int

	receiver    address.Address
	receiverBal big_spec.Int

	receipt types.MessageReceipt
}

func TestValueTransferSimple(t *testing.T, factories state.Factories) {
	alice := utils.NewSECP256K1Addr(t, "1")
	bob := utils.NewSECP256K1Addr(t, "2")

	builder := drivers.NewBuilder(context.Background(), factories).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	testCases := []valueTransferTestCases{
		{
			desc: "successfully transfer funds from sender to receiver",

			sender:    alice,
			senderBal: big_spec.NewInt(10_000_000),

			transferAmnt: big_spec.NewInt(50),

			receiver:    bob,
			receiverBal: big_spec.Zero(),

			receipt: types.MessageReceipt{
				ExitCode:    exitcode.Ok,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(128),
			},
		},
		{
			desc: "successfully transfer zero funds from sender to receiver",

			sender:    alice,
			senderBal: big_spec.NewInt(10_000_000),

			transferAmnt: big_spec.NewInt(0),

			receiver:    bob,
			receiverBal: big_spec.Zero(),

			receipt: types.MessageReceipt{
				ExitCode:    exitcode.Ok,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(114),
			},
		},
		{
			// Note: this test current fails for lotus as it returns an error instead of a message receipt
			desc: "fail to transfer more funds than sender balance > 0",

			sender:    alice,
			senderBal: big_spec.NewInt(10_000_000),

			transferAmnt: big_spec.NewInt(10_000_001),

			receiver:    bob,
			receiverBal: big_spec.Zero(),

			receipt: types.MessageReceipt{
				ExitCode:    exitcode.SysErrInsufficientFunds,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(1_000_000),
			},
		},
		{
			// Note: this test current fails for lotus as it returns an error instead of a message receipt
			desc: "fail to transfer more funds than sender has when sender balance == zero",

			sender:    alice,
			senderBal: big_spec.Zero(),

			transferAmnt: big_spec.NewInt(1),

			receiver:    bob,
			receiverBal: big_spec.Zero(),

			receipt: types.MessageReceipt{
				ExitCode:    exitcode.SysErrInsufficientFunds,
				ReturnValue: drivers.EmptyReturnValue,
				GasUsed:     big_spec.NewInt(1_000_000),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			td := builder.Build(t)

			// Create the to and from actors with balance in the state tree
			_, _, err := td.State().CreateActor(builtin_spec.AccountActorCodeID, tc.sender, tc.senderBal, &account_spec.State{Address: tc.sender})
			require.NoError(t, err)
			_, _, err = td.State().CreateActor(builtin_spec.AccountActorCodeID, tc.receiver, tc.receiverBal, &account_spec.State{Address: tc.receiver})
			require.NoError(t, err)

			gasUsed := td.ApplyMessageExpectResult(
				td.MessageProducer.Transfer(tc.receiver, tc.sender, chain.Value(tc.transferAmnt), chain.Nonce(0)),
				tc.receipt.ExitCode,
				tc.receipt.ReturnValue,
			)
			// create a message to transfer funds from `to` to `from` for amount `transferAmnt` and apply it to the state tree
			// assert the actor balances changed as expected, the receiver balance should not change if transfer fails
			if tc.receipt.ExitCode.IsSuccess() {
				td.AssertBalance(tc.sender, big_spec.Sub(big_spec.Sub(tc.senderBal, tc.transferAmnt), abi_spec.NewTokenAmount(gasUsed)))
				td.AssertBalance(tc.receiver, tc.transferAmnt)
			} else {
				td.AssertBalance(tc.sender, tc.senderBal)
			}

		})
	}
}

func TestValueTransferAdvance(t *testing.T, factory state.Factories) {
	var aliceBal = abi_spec.NewTokenAmount(1_000_000_000)

	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("self transfer", func(t *testing.T) {
		td := builder.Build(t)
		alice, _ := td.NewAccountActor(drivers.SECP, aliceBal)
		gasCost := abi_spec.NewTokenAmount(128)
		transferAmnt := abi_spec.NewTokenAmount(10)

		td.ApplyMessageExpectReceipt(
			td.MessageProducer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: gasCost},
		)
		td.AssertBalance(alice, big_spec.Sub(aliceBal, gasCost))
	})

	t.Run("transfer from known address to unknown account", func(t *testing.T) {
		td := builder.Build(t)

		alice, _ := td.NewAccountActor(drivers.SECP, aliceBal)
		unknown := td.Wallet().NewSECP256k1AccountAddress()
		gasCost := abi_spec.NewTokenAmount(128)
		transferAmnt := abi_spec.NewTokenAmount(10)

		td.ApplyMessageExpectReceipt(
			td.MessageProducer.Transfer(unknown, alice, chain.Value(transferAmnt), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: gasCost},
		)
		td.AssertBalance(alice, big_spec.Sub(big_spec.Sub(aliceBal, gasCost), transferAmnt))
		td.AssertBalance(unknown, transferAmnt)
	})

	t.Run("fail to transfer from unknown account to known address", func(t *testing.T) {
		td := builder.Build(t)
		alice, _ := td.NewAccountActor(drivers.SECP, aliceBal)
		unknown := td.Wallet().NewSECP256k1AccountAddress()
		gasCost := abi_spec.NewTokenAmount(1_000_000)
		transferAmnt := abi_spec.NewTokenAmount(10)

		td.ApplyMessageExpectReceipt(
			td.MessageProducer.Transfer(alice, unknown, chain.Value(transferAmnt), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.SysErrActorNotFound, ReturnValue: drivers.EmptyReturnValue, GasUsed: gasCost},
		)
	})

	t.Run("fail to transfer from unknown address to unknown address", func(t *testing.T) {
		td := builder.Build(t)
		unknown := td.Wallet().NewSECP256k1AccountAddress()
		nobody := td.Wallet().NewSECP256k1AccountAddress()
		gasCost := abi_spec.NewTokenAmount(1_000_000)
		transferAmnt := abi_spec.NewTokenAmount(10)

		td.ApplyMessageExpectReceipt(
			td.MessageProducer.Transfer(nobody, unknown, chain.Value(transferAmnt), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.SysErrActorNotFound, ReturnValue: drivers.EmptyReturnValue, GasUsed: gasCost},
		)
	})
}

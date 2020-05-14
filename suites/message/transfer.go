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

	code exitcode.ExitCode
}

func MessageTest_ValueTransferSimple(t *testing.T, factories state.Factories) {
	alice := utils.NewSECP256K1Addr(t, "1")
	bob := utils.NewSECP256K1Addr(t, "2")

	builder := drivers.NewBuilder(context.Background(), factories).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState...)

	testCases := []valueTransferTestCases{
		{
			desc: "successfully transfer funds from sender to receiver",

			sender:    alice,
			senderBal: big_spec.NewInt(10_000_000),

			transferAmnt: big_spec.NewInt(50),

			receiver:    bob,
			receiverBal: big_spec.Zero(),

			code: exitcode.Ok,
		},
		{
			desc: "successfully transfer zero funds from sender to receiver",

			sender:    alice,
			senderBal: big_spec.NewInt(10_000_000),

			transferAmnt: big_spec.NewInt(0),

			receiver:    bob,
			receiverBal: big_spec.Zero(),

			code: exitcode.Ok,
		},
		{
			desc: "fail to transfer more funds than sender balance > 0",

			sender:    alice,
			senderBal: big_spec.NewInt(10_000_000),

			transferAmnt: big_spec.NewInt(10_000_001),

			receiver:    bob,
			receiverBal: big_spec.Zero(),

			code: exitcode.SysErrSenderStateInvalid,
		},
		{
			desc: "fail to transfer more funds than sender has when sender balance == zero",

			sender:    alice,
			senderBal: big_spec.Zero(),

			transferAmnt: big_spec.NewInt(1),

			receiver:    bob,
			receiverBal: big_spec.Zero(),

			code: exitcode.SysErrSenderStateInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			td := builder.Build(t)
			defer td.Complete()

			// Create the to and from actors with balance in the state tree
			_, _, err := td.State().CreateActor(builtin_spec.AccountActorCodeID, tc.sender, tc.senderBal, &account_spec.State{Address: tc.sender})
			require.NoError(t, err)
			if tc.sender.String() != tc.receiver.String() {
				_, _, err := td.State().CreateActor(builtin_spec.AccountActorCodeID, tc.receiver, tc.receiverBal, &account_spec.State{Address: tc.receiver})
				require.NoError(t, err)
			}

			sendAct, err := td.State().Actor(tc.sender)
			require.NoError(t, err)
			require.Equal(t, tc.senderBal.String(), sendAct.Balance().String())

			result := td.ApplyFailure(
				td.MessageProducer.Transfer(tc.sender, tc.receiver, chain.Value(tc.transferAmnt), chain.Nonce(0)),
				tc.code,
			)
			// create a message to transfer funds from `to` to `from` for amount `transferAmnt` and apply it to the state tree
			// assert the actor balances changed as expected, the receiver balance should not change if transfer fails
			if tc.code.IsSuccess() {
				td.AssertBalance(tc.sender, big_spec.Sub(big_spec.Sub(tc.senderBal, tc.transferAmnt), result.Receipt.GasUsed.Big()))
				td.AssertBalance(tc.receiver, tc.transferAmnt)
			} else {
				td.AssertBalance(tc.sender, tc.senderBal)
			}
		})
	}
}

func MessageTest_ValueTransferAdvance(t *testing.T, factory state.Factories) {
	var aliceInitialBalance = abi_spec.NewTokenAmount(1_000_000_000)

	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState...)

	t.Run("self transfer", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		alice, _ := td.NewAccountActor(drivers.SECP, aliceInitialBalance)
		transferAmnt := abi_spec.NewTokenAmount(10)

		result := td.ApplyOk(
			td.MessageProducer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(0)))
		// since this is a self transfer expect alice's balance to only decrease by the gasUsed
		td.AssertBalance(alice, big_spec.Sub(aliceInitialBalance, result.Receipt.GasUsed.Big()))
	})

	t.Run("ok transfer from known address to new account", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		alice, _ := td.NewAccountActor(drivers.SECP, aliceInitialBalance)
		receiver := td.Wallet().NewSECP256k1AccountAddress()
		transferAmnt := abi_spec.NewTokenAmount(10)

		result := td.ApplyOk(
			td.MessageProducer.Transfer(alice, receiver, chain.Value(transferAmnt), chain.Nonce(0)),
		)
		td.AssertBalance(alice, big_spec.Sub(big_spec.Sub(aliceInitialBalance, result.Receipt.GasUsed.Big()), transferAmnt))
		td.AssertBalance(receiver, transferAmnt)
	})

	t.Run("fail to transfer from unknown account to known address", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		alice, _ := td.NewAccountActor(drivers.SECP, aliceInitialBalance)
		unknown := td.Wallet().NewSECP256k1AccountAddress()
		transferAmnt := abi_spec.NewTokenAmount(10)

		td.ApplyFailure(
			td.MessageProducer.Transfer(unknown, alice, chain.Value(transferAmnt), chain.Nonce(0)),
			exitcode.SysErrSenderInvalid)
		td.AssertBalance(alice, aliceInitialBalance)
	})

	t.Run("fail to transfer from unknown address to unknown address", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		sender := td.Wallet().NewSECP256k1AccountAddress()
		receiver := td.Wallet().NewSECP256k1AccountAddress()
		transferAmnt := abi_spec.NewTokenAmount(10)

		td.ApplyFailure(
			td.MessageProducer.Transfer(sender, receiver, chain.Value(transferAmnt), chain.Nonce(0)),
			exitcode.SysErrSenderInvalid)
		td.AssertNoActor(sender)
		td.AssertNoActor(receiver)
	})
}

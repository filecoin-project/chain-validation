package suites

import (
	"context"
	"testing"

	address "github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	init_spec "github.com/filecoin-project/specs-actors/actors/builtin/init"
	reward_spec "github.com/filecoin-project/specs-actors/actors/builtin/reward"
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
	defaultMiner := utils.NewBLSAddr(t, 123)
	alice := utils.NewSECP256K1Addr(t, "1")
	bob := utils.NewSECP256K1Addr(t, "2")

	builder := drivers.NewBuilder(context.Background(), factories).
		WithDefaultGasLimit(big_spec.NewInt(1000000)).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithDefaultMiner(defaultMiner).
		WithActorState([]drivers.ActorState{
			{
				Addr:    builtin_spec.InitActorAddr,
				Balance: big_spec.Zero(),
				Code:    builtin_spec.InitActorCodeID,
				State:   init_spec.ConstructState(drivers.EmptyMapCid, "chain-validation"),
			},
			{
				Addr:    builtin_spec.RewardActorAddr,
				Balance: TotalNetworkBalance,
				Code:    builtin_spec.RewardActorCodeID,
				State:   reward_spec.ConstructState(drivers.EmptyMultiMapCid),
			},
			{
				Addr:    builtin_spec.BurntFundsActorAddr,
				Balance: big_spec.Zero(),
				Code:    builtin_spec.AccountActorCodeID,
				State:   &account_spec.State{Address: builtin_spec.BurntFundsActorAddr},
			},
		})

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
				ReturnValue: EmptyReturnValue,
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
				ReturnValue: EmptyReturnValue,
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
				ReturnValue: EmptyReturnValue,
				GasUsed:     big_spec.NewInt(0),
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
				ReturnValue: EmptyReturnValue,
				GasUsed:     big_spec.NewInt(0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			td := builder.Build(t)

			// Create the to and from actors with balance in the state tree
			_, err := td.State().CreateActor(builtin_spec.AccountActorCodeID, tc.sender, tc.senderBal, &account_spec.State{Address: tc.sender})
			require.NoError(t, err)
			if tc.sender.String() != tc.receiver.String() {
				_, err := td.State().CreateActor(builtin_spec.AccountActorCodeID, tc.receiver, tc.receiverBal, &account_spec.State{Address: tc.receiver})
				require.NoError(t, err)
			}

			sendAct, err := td.State().Actor(tc.sender)
			require.NoError(t, err)
			require.Equal(t, tc.senderBal.String(), sendAct.Balance().String())

			td.ApplyMessageExpectReceipt(
				td.Producer.Transfer(tc.receiver, tc.sender, chain.Value(tc.transferAmnt), chain.Nonce(0)),
				tc.receipt,
			)
			// create a message to transfer funds from `to` to `from` for amount `transferAmnt` and apply it to the state tree
			// assert the actor balances changed as expected, the receiver balance should not change if transfer fails
			if tc.receipt.ExitCode.IsSuccess() {
				td.AssertBalanceWithGas(tc.sender, big_spec.Sub(tc.senderBal, tc.transferAmnt), tc.receipt.GasUsed)
				td.AssertBalance(tc.receiver, tc.transferAmnt)
			} else {
				td.AssertBalanceWithGas(tc.sender, big_spec.Sub(tc.senderBal, tc.transferAmnt), tc.receipt.GasUsed)
			}

		})
	}

}

func TestValueTransferAdvance(t *testing.T, factory state.Factories) {
	var gasCost = big_spec.Zero()
	var aliceBal = abi_spec.NewTokenAmount(1_000_000_000)
	var transferAmnt = abi_spec.NewTokenAmount(10)

	defaultMiner := utils.NewBLSAddr(t, 123)

	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(big_spec.NewInt(1000000)).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithDefaultMiner(defaultMiner).
		WithActorState([]drivers.ActorState{
			{
				Addr:    builtin_spec.InitActorAddr,
				Balance: big_spec.Zero(),
				Code:    builtin_spec.InitActorCodeID,
				State:   init_spec.ConstructState(drivers.EmptyMapCid, "chain-validation"),
			},
			{
				Addr:    builtin_spec.RewardActorAddr,
				Balance: TotalNetworkBalance,
				Code:    builtin_spec.RewardActorCodeID,
				State:   reward_spec.ConstructState(drivers.EmptyMultiMapCid),
			},
			{
				Addr:    builtin_spec.BurntFundsActorAddr,
				Balance: big_spec.Zero(),
				Code:    builtin_spec.AccountActorCodeID,
				State:   &account_spec.State{Address: builtin_spec.BurntFundsActorAddr},
			},
		})

	t.Run("self transfer", func(t *testing.T) {
		td := builder.Build(t)
		alice := td.NewAccountActor(drivers.SECP, aliceBal)

		td.ApplyMessageExpectReceipt(
			td.Producer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: EmptyReturnValue, GasUsed: gasCost},
		)
		td.AssertBalanceWithGas(alice, aliceBal, gasCost)
	})
	t.Run("transfer from known address to unknown account", func(t *testing.T) {
		td := builder.Build(t)

		alice := td.NewAccountActor(drivers.SECP, aliceBal)
		unknown := td.Wallet().NewSECP256k1AccountAddress()

		td.ApplyMessageExpectReceipt(
			td.Producer.Transfer(unknown, alice, chain.Value(transferAmnt), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: EmptyReturnValue, GasUsed: gasCost},
		)
		td.AssertBalanceWithGas(alice, aliceBal, gasCost)
	})

	t.Run("fail to transfer from unknown account to known address", func(t *testing.T) {
		td := builder.Build(t)
		alice := td.NewAccountActor(drivers.SECP, aliceBal)
		unknown := td.Wallet().NewSECP256k1AccountAddress()

		td.ApplyMessageExpectReceipt(
			td.Producer.Transfer(alice, unknown, chain.Value(transferAmnt), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.SysErrActorNotFound, ReturnValue: EmptyReturnValue, GasUsed: gasCost},
		)
		td.AssertBalanceWithGas(alice, aliceBal, gasCost)
	})

	t.Run("fail to transfer from unknown address to unknown address", func(t *testing.T) {
		td := builder.Build(t)
		unknown := td.Wallet().NewSECP256k1AccountAddress()
		nobody := td.Wallet().NewSECP256k1AccountAddress()

		td.ApplyMessageExpectReceipt(
			td.Producer.Transfer(nobody, unknown, chain.Value(transferAmnt), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.SysErrActorNotFound, ReturnValue: EmptyReturnValue, GasUsed: gasCost},
		)
	})
}

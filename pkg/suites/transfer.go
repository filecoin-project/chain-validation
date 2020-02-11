package suites

import (
	"context"
	"testing"

	address "github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	require "github.com/stretchr/testify/require"

	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"

	chain "github.com/filecoin-project/chain-validation/pkg/chain"
)

type valueTransferTestCases struct {
	desc string

	sender    address.Address
	senderBal big_spec.Int

	transferAmnt big_spec.Int

	receiver    address.Address
	receiverBal big_spec.Int

	receipt chain.MessageReceipt
}

func TestValueTransferSimple(t *testing.T, factories Factories) {
	defaultMiner, err := address.NewSecp256k1Address([]byte{'m', 'i', 'n', 'e', 'r'})
	require.NoError(t, err)

	alice, err := address.NewSecp256k1Address([]byte{'1'})
	require.NoError(t, err)

	bob, err := address.NewSecp256k1Address([]byte{'2'})
	require.NoError(t, err)

	builder := NewBuilder(context.Background(), factories).
		WithDefaultGasLimit(big_spec.NewInt(1000000)).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithSingletonActors(map[address.Address]big_spec.Int{
			builtin_spec.InitActorAddr: big_spec.Zero(),
		}).
		WithDefaultMiner(defaultMiner)

	testCases := []valueTransferTestCases{
		{
			desc: "successfully transfer funds from sender to receiver",

			sender:    alice,
			senderBal: big_spec.NewInt(10_000_000),

			transferAmnt: big_spec.NewInt(50),

			receiver:    bob,
			receiverBal: big_spec.Zero(),

			receipt: chain.MessageReceipt{
				ExitCode:    exitcode.Ok,
				ReturnValue: nil,
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

			receipt: chain.MessageReceipt{
				ExitCode:    exitcode.Ok,
				ReturnValue: nil,
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

			receipt: chain.MessageReceipt{
				ExitCode:    exitcode.ErrInsufficientFunds,
				ReturnValue: nil,
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

			receipt: chain.MessageReceipt{
				ExitCode:    exitcode.ErrInsufficientFunds,
				ReturnValue: nil,
				GasUsed:     big_spec.NewInt(0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			td := builder.Build(t)

			// Create the to and from actors with balance in the state tree
			_, _, err = td.Driver.st.SetActor(tc.sender, builtin_spec.AccountActorCodeID, tc.senderBal)
			require.NoError(t, err)
			if tc.sender.String() != tc.receiver.String() {
				_, _, err := td.Driver.st.SetActor(tc.receiver, builtin_spec.AccountActorCodeID, tc.receiverBal)
				require.NoError(t, err)
			}

			// create a message to transfer funds from `to` to `from` for amount `transferAmnt` and apply it to the state tree
			transferMsg := td.Producer.Transfer(tc.receiver, tc.sender, chain.Value(tc.transferAmnt), chain.Nonce(0))
			transferReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), transferMsg)
			require.NoError(t, err)

			// assert the transfer message application returned the expected exitcode and gas cast
			td.Driver.AssertReceipt(transferReceipt, tc.receipt)

			// assert the actor balances changed as expected, the receiver balance should not change if transfer fails
			if tc.receipt.ExitCode.IsSuccess() {
				td.Driver.AssertBalance(tc.sender, big_spec.Sub(big_spec.Sub(tc.senderBal, tc.transferAmnt), tc.receipt.GasUsed))
				td.Driver.AssertBalance(tc.receiver, tc.transferAmnt)
				td.Driver.AssertBalance(td.ExeCtx.MinerOwner, tc.receipt.GasUsed)
			} else {
				td.Driver.AssertBalance(tc.sender, big_spec.Sub(big_spec.Sub(tc.senderBal, tc.transferAmnt), tc.receipt.GasUsed))
				td.Driver.AssertBalance(td.ExeCtx.MinerOwner, tc.receipt.GasUsed)
			}

		})
	}

}

func TestValueTransferAdvance(t *testing.T, factory Factories) {
	var gasCost = big_spec.Zero()
	var aliceBal = abi_spec.NewTokenAmount(100)
	var transferAmnt = abi_spec.NewTokenAmount(10)

	defaultMiner, err := address.NewSecp256k1Address([]byte{'m', 'i', 'n', 'e', 'r'})
	require.NoError(t, err)

	builder := NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(big_spec.NewInt(1000000)).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithSingletonActors(map[address.Address]big_spec.Int{
			builtin_spec.InitActorAddr: big_spec.Zero(),
		}).
		WithDefaultMiner(defaultMiner)

	t.Run("fail to self transfer", func(t *testing.T) {
		td := builder.Build(t)
		alice := td.Driver.NewAccountActor(SECP, aliceBal)

		msg := td.Producer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(0))

		msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    0,
			ReturnValue: nil,
			GasUsed:     gasCost,
		})

		td.Driver.AssertBalance(alice, big_spec.Sub(aliceBal, gasCost))
		td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
	})
	t.Run("fail to transfer from known address to unknown account", func(t *testing.T) {
		td := builder.Build(t)
		alice := td.Driver.NewAccountActor(SECP, aliceBal)

		unknown := td.Driver.State().NewSecp256k1AccountAddress()

		msg := td.Producer.Transfer(unknown, alice, chain.Value(transferAmnt), chain.Nonce(0))

		msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    0,
			ReturnValue: nil,
			GasUsed:     gasCost,
		})

		td.Driver.AssertBalance(alice, big_spec.Sub(aliceBal, gasCost))
		td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
	})

	t.Run("fail to transfer from unknown account to known address", func(t *testing.T) {
		td := builder.Build(t)
		alice := td.Driver.NewAccountActor(SECP, aliceBal)
		unknown := td.Driver.State().NewSecp256k1AccountAddress()

		msg := td.Producer.Transfer(alice, unknown, chain.Value(transferAmnt), chain.Nonce(0))

		msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    0,
			ReturnValue: nil,
			GasUsed:     gasCost,
		})

		td.Driver.AssertBalance(alice, big_spec.Sub(aliceBal, gasCost))
		td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
	})

	t.Run("fail to transfer from unknown address to unknown address", func(t *testing.T) {
		td := builder.Build(t)
		unknown := td.Driver.State().NewSecp256k1AccountAddress()
		nobody := td.Driver.State().NewSecp256k1AccountAddress()

		msg := td.Producer.Transfer(nobody, unknown, chain.Value(transferAmnt), chain.Nonce(0))

		msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
		require.NoError(t, err)
		td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
			ExitCode:    0,
			ReturnValue: nil,
			GasUsed:     gasCost,
		})

		td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
	})
}

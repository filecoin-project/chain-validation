package suites

import (
	"testing"

	address "github.com/filecoin-project/go-address"
	require "github.com/stretchr/testify/require"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"

	chain "github.com/filecoin-project/chain-validation/pkg/chain"
)

func AccountValueTransferSuccess(t *testing.T, factory Factories) {
	var gasCost = big_spec.NewInt(128)
	var aliceBal = abi_spec.NewTokenAmount(20000000000)
	var bobBal = big_spec.Zero()
	var transferValue = abi_spec.NewTokenAmount(50)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr: big_spec.Zero(),
	})

	alice := td.Driver.NewAccountActor(SECP, aliceBal)
	bob := td.Driver.NewAccountActor(SECP, bobBal)

	msg := td.Producer.Transfer(bob, alice, chain.Value(transferValue), chain.Nonce(0))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(t, err)
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     gasCost,
	})

	td.Driver.AssertBalance(alice, big_spec.Sub(big_spec.Sub(aliceBal, transferValue), gasCost))
	td.Driver.AssertBalance(bob, transferValue)
	td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
}

func AccountValueTransferZeroFunds(t *testing.T, factory Factories) {
	var gasCost = big_spec.NewInt(114)
	var aliceBal = abi_spec.NewTokenAmount(20000000000)
	var bobBal = abi_spec.NewTokenAmount(0)
	var transferValue = big_spec.Zero()

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr: big_spec.Zero(),
	})

	alice := td.Driver.NewAccountActor(SECP, aliceBal)
	bob := td.Driver.NewAccountActor(SECP, bobBal)

	msg := td.Producer.Transfer(bob, alice, chain.Value(transferValue), chain.Nonce(0))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(t, err)
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     gasCost,
	})

	td.Driver.AssertBalance(alice, big_spec.Sub(big_spec.Sub(aliceBal, transferValue), gasCost))
	td.Driver.AssertBalance(bob, transferValue)
	td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
}

func AccountValueTransferOverBalanceNonZero(t *testing.T, factory Factories) {
	var gasCost = big_spec.Zero()
	var aliceBal = abi_spec.NewTokenAmount(2000)
	var bobBal = big_spec.Zero()
	var transferAmnt = big_spec.Add(aliceBal, abi_spec.NewTokenAmount(1))

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr: big_spec.Zero(),
	})

	alice := td.Driver.NewAccountActor(SECP, aliceBal)
	bob := td.Driver.NewAccountActor(SECP, bobBal)

	msg := td.Producer.Transfer(bob, alice, chain.Value(transferAmnt), chain.Nonce(0))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.Error(t, err)
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     gasCost,
	})

	td.Driver.AssertBalance(alice, big_spec.Sub(aliceBal, gasCost))
	td.Driver.AssertBalance(bob, bobBal)
	td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
}

func AccountValueTransferOverBalanceZero(t *testing.T, factory Factories) {
	var gasCost = big_spec.Zero()
	var aliceBal = big_spec.Zero()
	var bobBal = big_spec.Zero()
	var transferAmnt = big_spec.Add(aliceBal, abi_spec.NewTokenAmount(1))

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr: big_spec.Zero(),
	})

	alice := td.Driver.NewAccountActor(SECP, aliceBal)
	bob := td.Driver.NewAccountActor(SECP, bobBal)

	msg := td.Producer.Transfer(bob, alice, chain.Value(transferAmnt), chain.Nonce(0))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.Error(t, err)
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     gasCost,
	})

	td.Driver.AssertBalance(alice, aliceBal)
	td.Driver.AssertBalance(bob, bobBal)
	td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
}

func AccountValueTransferToSelf(t *testing.T, factory Factories) {
	var gasCost = big_spec.Zero()
	var aliceBal = abi_spec.NewTokenAmount(100)
	var transferAmnt = abi_spec.NewTokenAmount(10)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr: big_spec.Zero(),
	})

	alice := td.Driver.NewAccountActor(SECP, aliceBal)

	msg := td.Producer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(0))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.Error(t, err)
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     gasCost,
	})

	td.Driver.AssertBalance(alice, big_spec.Sub(aliceBal, gasCost))
	td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
}

func AccountValueTransferFromKnownToUnknownAccount(t *testing.T, factory Factories) {
	var gasCost = big_spec.Zero()
	var aliceBal = abi_spec.NewTokenAmount(100)
	var transferAmnt = abi_spec.NewTokenAmount(10)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr: big_spec.Zero(),
	})

	alice := td.Driver.NewAccountActor(SECP, aliceBal)

	unknown := td.Driver.State().NewSecp256k1AccountAddress()

	msg := td.Producer.Transfer(unknown, alice, chain.Value(transferAmnt), chain.Nonce(0))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.Error(t, err)
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     gasCost,
	})

	td.Driver.AssertBalance(alice, big_spec.Sub(aliceBal, gasCost))
	td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
}

func AccountValueTransferFromUnknownToKnownAccount(t *testing.T, factory Factories) {
	var gasCost = big_spec.Zero()
	var aliceBal = abi_spec.NewTokenAmount(100)
	var transferAmnt = abi_spec.NewTokenAmount(10)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr: big_spec.Zero(),
	})

	alice := td.Driver.NewAccountActor(SECP, aliceBal)
	unknown := td.Driver.State().NewSecp256k1AccountAddress()

	msg := td.Producer.Transfer(alice, unknown, chain.Value(transferAmnt), chain.Nonce(0))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.Error(t, err)
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     gasCost,
	})

	td.Driver.AssertBalance(alice, big_spec.Sub(aliceBal, gasCost))
	td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
}

func AccountValueTransferFromUnknownToUnknownAccount(t *testing.T, factory Factories) {
	var gasCost = big_spec.Zero()
	var transferAmnt = abi_spec.NewTokenAmount(10)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr: big_spec.Zero(),
	})

	unknown := td.Driver.State().NewSecp256k1AccountAddress()
	nobody := td.Driver.State().NewSecp256k1AccountAddress()

	msg := td.Producer.Transfer(nobody, unknown, chain.Value(transferAmnt), chain.Nonce(0))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.Error(t, err)
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     gasCost,
	})

	td.Driver.AssertBalance(td.ExeCtx.MinerOwner, gasCost)
}

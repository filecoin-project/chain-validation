package suites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

func AccountValueTransferSuccess(t *testing.T, factory Factories, expGasUsed uint64) {
	const initialBal = 20000000000
	const transferValue = 50
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress: types.NewInt(0),
	})

	alice := c.Driver().NewAccountActor(initialBal)
	bob := c.Driver().NewAccountActor(0)

	msg, err := c.Producer().Transfer(alice, bob, 0, transferValue)
	require.NoError(t, err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(t, err)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})

	c.Driver().AssertBalance(alice, initialBal-transferValue-expGasUsed)
	c.Driver().AssertBalance(bob, transferValue)
	// This should become non-zero after gas tracking and payments are integrated.
	c.Driver().AssertBalance(c.ExeCtx().MinerOwner, expGasUsed)
}

func AccountValueTransferZeroFunds(t *testing.T, factory Factories, expGasUsed uint64) {
	const initialBal = 20000000000
	const transferValue = 0
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress: types.NewInt(0),
	})

	alice := c.Driver().NewAccountActor(initialBal)
	bob := c.Driver().NewAccountActor(0)

	msg, err := c.Producer().Transfer(alice, bob, 0, transferValue)
	require.NoError(t, err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(t, err)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})

	c.Driver().AssertBalance(alice, initialBal-transferValue-expGasUsed)
	c.Driver().AssertBalance(bob, transferValue)
	// This should become non-zero after gas tracking and payments are integrated.
	c.Driver().AssertBalance(c.ExeCtx().MinerOwner, expGasUsed)
}

func AccountValueTransferOverBalanceNonZero(t *testing.T, factory Factories, expGasUsed uint64) {
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress: types.NewInt(0),
	})

	alice := c.Driver().NewAccountActor(2000)
	bob := c.Driver().NewAccountActor(0)

	msg, err := c.Producer().Transfer(alice, bob, 0, 2001)
	require.NoError(t, err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.Error(t, err)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(expGasUsed),
	})

	c.Driver().AssertBalance(alice, 2000-expGasUsed)
	c.Driver().AssertBalance(bob, 0)
	// This should become non-zero after gas tracking and payments are integrated.
	c.Driver().AssertBalance(c.ExeCtx().MinerOwner, expGasUsed)
}

func AccountValueTransferOverBalanceZero(t *testing.T, factory Factories, expGasUsed uint64) {
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress: types.NewInt(0),
	})

	alice := c.Driver().NewAccountActor(0)
	bob := c.Driver().NewAccountActor(0)

	msg, err := c.Producer().Transfer(alice, bob, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.Error(t, err)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(expGasUsed),
	})

	c.Driver().AssertBalance(alice, 0)
	c.Driver().AssertBalance(bob, 0)
	// This should become non-zero after gas tracking and payments are integrated.
	c.Driver().AssertBalance(c.ExeCtx().MinerOwner, expGasUsed)
}

func AccountValueTransferToSelf(t *testing.T, factory Factories, expGasUsed uint64) {
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress: types.NewInt(0),
	})

	alice := c.Driver().NewAccountActor(1)

	msg, err := c.Producer().Transfer(alice, alice, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.Error(t, err)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(expGasUsed),
	})

	c.Driver().AssertBalance(alice, 1)
	// This should become non-zero after gas tracking and payments are integrated.
	c.Driver().AssertBalance(c.ExeCtx().MinerOwner, expGasUsed)
}

func AccountValueTransferFromKnownToUnknownAccount(t *testing.T, factory Factories, expGasUsed uint64) {
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress: types.NewInt(0),
	})

	alice := c.Driver().NewAccountActor(1)
	unknown, err := c.Driver().State().NewAccountAddress()
	require.NoError(t, err)

	msg, err := c.Producer().Transfer(alice, unknown, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.Error(t, err)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(expGasUsed),
	})

	c.Driver().AssertBalance(alice, 1)
	// This should become non-zero after gas tracking and payments are integrated.
	c.Driver().AssertBalance(c.ExeCtx().MinerOwner, expGasUsed)
}

func AccountValueTransferFromUnknownToKnownAccount(t *testing.T, factory Factories, expGasUsed uint64) {
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress: types.NewInt(0),
	})

	alice := c.Driver().NewAccountActor(1)
	unknown, err := c.Driver().State().NewAccountAddress()
	require.NoError(t, err)

	msg, err := c.Producer().Transfer(unknown, alice, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.Error(t, err)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(expGasUsed),
	})

	c.Driver().AssertBalance(alice, 1)
	// This should become non-zero after gas tracking and payments are integrated.
	c.Driver().AssertBalance(c.ExeCtx().MinerOwner, expGasUsed)
}

func AccountValueTransferFromUnknownToUnknownAccount(t *testing.T, factory Factories, expGasUsed uint64) {
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress: types.NewInt(0),
	})

	alice := c.Driver().NewAccountActor(1)
	unknown, err := c.Driver().State().NewAccountAddress()
	require.NoError(t, err)

	nobody, err := c.Driver().State().NewAccountAddress()
	require.NoError(t, err)

	msg, err := c.Producer().Transfer(unknown, nobody, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.Error(t, err)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(expGasUsed),
	})

	c.Driver().AssertBalance(alice, 1)
	// This should become non-zero after gas tracking and payments are integrated.
	c.Driver().AssertBalance(c.ExeCtx().MinerOwner, expGasUsed)
}

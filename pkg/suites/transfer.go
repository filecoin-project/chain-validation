package suites

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

func transferTestSetup(t *testing.T, factory Factories) (*StateDriver, *chain.MessageProducer, *chain.Validator) {
	drv := NewStateDriver(t, factory.NewState())

	_, _, err := drv.State().SetSingletonActor(state.InitAddress, types.NewInt(0))
	require.NoError(t, err)

	gasPrice := types.NewInt(1)
	gasLimit := types.GasUnit(1000)

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	return drv, producer, validator
}

func AccountValueTransferSuccess(t *testing.T, factory Factories, expGasUsed uint64) {
	drv, producer, validator := transferTestSetup(t, factory)

	alice := drv.NewAccountActor(2000)
	bob := drv.NewAccountActor(0)
	miner := drv.NewAccountActor(0) // Miner owner

	exeCtx := chain.NewExecutionContext(1, miner)

	msg, err := producer.Transfer(alice, bob, 0, 50)
	require.NoError(t, err)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Empty(t, msgReceipt.ReturnValue)
	assert.Equal(t, types.GasUnit(expGasUsed), msgReceipt.GasUsed)

	drv.AssertBalance(alice, 1950-expGasUsed)
	drv.AssertBalance(bob, 50)
	// This should become non-zero after gas tracking and payments are integrated.
	drv.AssertBalance(miner, expGasUsed)
}

func AccountValueTransferZeroFunds(t *testing.T, factory Factories, expGasUsed uint64) {
	drv, producer, validator := transferTestSetup(t, factory)

	alice := drv.NewAccountActor(2000)
	bob := drv.NewAccountActor(0)
	miner := drv.NewAccountActor(0) // Miner owner

	exeCtx := chain.NewExecutionContext(1, miner)

	msg, err := producer.Transfer(alice, bob, 0, 0)
	require.NoError(t, err)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Empty(t, msgReceipt.ReturnValue)
	assert.Equal(t, types.GasUnit(expGasUsed), msgReceipt.GasUsed)

	drv.AssertBalance(alice, 2000-expGasUsed)
	drv.AssertBalance(bob, 0)
	// This should become non-zero after gas tracking and payments are integrated.
	drv.AssertBalance(miner, expGasUsed)
}

func AccountValueTransferOverBalanceNonZero(t *testing.T, factory Factories, expGasUsed uint64) {
	drv, producer, validator := transferTestSetup(t, factory)

	alice := drv.NewAccountActor(2000)
	bob := drv.NewAccountActor(0)
	miner := drv.NewAccountActor(0) // Miner owner

	exeCtx := chain.NewExecutionContext(1, miner)

	msg, err := producer.Transfer(alice, bob, 0, 2001)
	require.NoError(t, err)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.Error(t, err)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Empty(t, msgReceipt.ReturnValue)
	assert.Equal(t, types.GasUnit(expGasUsed), msgReceipt.GasUsed)

	drv.AssertBalance(alice, 2000-expGasUsed)
	drv.AssertBalance(bob, 0)
	// This should become non-zero after gas tracking and payments are integrated.
	drv.AssertBalance(miner, expGasUsed)
}

func AccountValueTransferOverBalanceZero(t *testing.T, factory Factories, expGasUsed uint64) {
	drv, producer, validator := transferTestSetup(t, factory)

	alice := drv.NewAccountActor(0)
	bob := drv.NewAccountActor(0)
	miner := drv.NewAccountActor(0) // Miner owner

	exeCtx := chain.NewExecutionContext(1, miner)

	msg, err := producer.Transfer(alice, bob, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.Error(t, err)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Empty(t, msgReceipt.ReturnValue)
	assert.Equal(t, types.GasUnit(expGasUsed), msgReceipt.GasUsed)

	drv.AssertBalance(alice, 0)
	drv.AssertBalance(bob, 0)
	// This should become non-zero after gas tracking and payments are integrated.
	drv.AssertBalance(miner, expGasUsed)
}

func AccountValueTransferToSelf(t *testing.T, factory Factories, expGasUsed uint64) {
	drv, producer, validator := transferTestSetup(t, factory)

	alice := drv.NewAccountActor(1)
	miner := drv.NewAccountActor(0) // Miner owner

	exeCtx := chain.NewExecutionContext(1, miner)

	msg, err := producer.Transfer(alice, alice, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.Error(t, err)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Empty(t, msgReceipt.ReturnValue)
	assert.Equal(t, types.GasUnit(expGasUsed), msgReceipt.GasUsed)

	drv.AssertBalance(alice, 1)
	// This should become non-zero after gas tracking and payments are integrated.
	drv.AssertBalance(miner, expGasUsed)
}

func AccountValueTransferFromKnownToUnknownAccount(t *testing.T, factory Factories, expGasUsed uint64) {
	drv, producer, validator := transferTestSetup(t, factory)

	alice := drv.NewAccountActor(1)
	miner := drv.NewAccountActor(0) // Miner owner
	unknown, err := drv.State().NewAccountAddress()
	require.NoError(t, err)

	exeCtx := chain.NewExecutionContext(1, miner)

	msg, err := producer.Transfer(alice, unknown, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.Error(t, err)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Empty(t, msgReceipt.ReturnValue)
	assert.Equal(t, types.GasUnit(expGasUsed), msgReceipt.GasUsed)

	drv.AssertBalance(alice, 1)
	// This should become non-zero after gas tracking and payments are integrated.
	drv.AssertBalance(miner, expGasUsed)

}

func AccountValueTransferFromUnknownToKnownAccount(t *testing.T, factory Factories, expGasUsed uint64) {
	drv, producer, validator := transferTestSetup(t, factory)

	alice := drv.NewAccountActor(1)
	miner := drv.NewAccountActor(0) // Miner owner
	unknown, err := drv.State().NewAccountAddress()
	require.NoError(t, err)

	exeCtx := chain.NewExecutionContext(1, miner)

	msg, err := producer.Transfer(unknown, alice, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.Error(t, err)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Empty(t, msgReceipt.ReturnValue)
	assert.Equal(t, types.GasUnit(expGasUsed), msgReceipt.GasUsed)

	drv.AssertBalance(alice, 1)
	// This should become non-zero after gas tracking and payments are integrated.
	drv.AssertBalance(miner, expGasUsed)

}

func AccountValueTransferFromUnknownToUnknownAccount(t *testing.T, factory Factories, expGasUsed uint64) {
	drv, producer, validator := transferTestSetup(t, factory)

	alice := drv.NewAccountActor(1)
	miner := drv.NewAccountActor(0) // Miner owner
	unknown, err := drv.State().NewAccountAddress()
	require.NoError(t, err)

	nobody, err := drv.State().NewAccountAddress()
	require.NoError(t, err)

	exeCtx := chain.NewExecutionContext(1, miner)

	msg, err := producer.Transfer(unknown, nobody, 0, 1)
	require.NoError(t, err)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.Error(t, err)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Empty(t, msgReceipt.ReturnValue)
	assert.Equal(t, types.GasUnit(expGasUsed), msgReceipt.GasUsed)

	drv.AssertBalance(alice, 1)
	// This should become non-zero after gas tracking and payments are integrated.
	drv.AssertBalance(miner, expGasUsed)

}

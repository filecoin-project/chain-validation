package suites

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
)

// A basic example validation test.
// At present this code is verbose and demonstrates the opportunity for helper methods.
func Example(t *testing.T, factories Factories) {
	drv := NewStateDriver(t, factories.NewState())

	alice := drv.NewAccountActor(2000)
	bob := drv.NewAccountActor(0)
	miner := drv.NewAccountActor(0) // Miner owner

	gasPrice := big.NewInt(1)
	gasLimit := state.GasUnit(1000)

	producer := chain.NewMessageProducer(factories.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	msg, err := producer.Transfer(alice, bob, 0, 50)
	require.NoError(t, err)

	validator := chain.NewValidator(factories)
	exeCtx := chain.NewExecutionContext(1, miner)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Equal(t, []byte{}, msgReceipt.ReturnValue)
	assert.Equal(t, state.GasUnit(0), msgReceipt.GasUsed)

	drv.AssertBalance(alice, 1950)
	drv.AssertBalance(bob, 50)
	// This should become non-zero after gas tracking and payments are integrated.
	drv.AssertBalance(miner, 0)
}

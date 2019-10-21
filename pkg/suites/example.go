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
func Example(t *testing.T, driver Driver) {
	actors := make(map[state.Address]state.Actor)

	alice, err := driver.NewAddress()
	require.NoError(t, err)
	actors[alice] = driver.NewActor(state.AccountActorCodeCid, big.NewInt(2000))

	bob, err := driver.NewAddress()
	require.NoError(t, err)
	actors[bob] = driver.NewActor(state.AccountActorCodeCid, big.NewInt(0))

	miner, err := driver.NewAddress()
	require.NoError(t, err)
	actors[miner] = driver.NewActor(state.AccountActorCodeCid, big.NewInt(0))

	tree, storage, err := driver.NewState(actors)
	require.NoError(t, err)
	require.NotNil(t, tree)
	require.NotNil(t, storage)

	producer := chain.NewMessageProducer(driver, state.GasUnit(1000), big.NewInt(1))
	msg, err := producer.Transfer(alice, bob, 0, 50)

	exeCtx := chain.NewExecutionContext(1, miner)
	validator := chain.NewValidator(driver)

	endState, msgReceipt, err := validator.ApplyMessage(exeCtx, tree, storage, msg)
	require.NoError(t, err)
	require.NotNil(t, endState)
	require.NotNil(t, msgReceipt)

	assert.Equal(t, uint8(0), msgReceipt.ExitCode)
	assert.Equal(t, []byte{}, msgReceipt.ReturnValue)
	assert.Equal(t, state.GasUnit(0), msgReceipt.GasUsed)

	actorAlice, err := endState.Actor(alice)
	require.NoError(t, err)
	assert.Equal(t, state.AttoFIL(big.NewInt(1950)), actorAlice.Balance())

	actorBob, err := endState.Actor(bob)
	require.NoError(t, err)
	assert.Equal(t, state.AttoFIL(big.NewInt(50)), actorBob.Balance())

	// This should become non-zero after gas tracking and payments are integrated.
	actorMiner, err := endState.Actor(miner)
	require.NoError(t, err)
	assert.Equal(t, state.AttoFIL(big.NewInt(0)), actorMiner.Balance())
}

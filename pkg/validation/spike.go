package validation

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
)

// Get this method to run
func TryItOut(t *testing.T, msgFactory chain.MessageFactory, stateFactory state.StateFactory) {
	actors := make(map[state.Address]state.Actor)

	alice, err := stateFactory.NewAddress()
	require.NoError(t, err)
	actors[alice] = stateFactory.NewActor(state.AccountActorCodeCid, big.NewInt(100))

	bob, err := stateFactory.NewAddress()
	require.NoError(t, err)
	actors[bob] = stateFactory.NewActor(state.AccountActorCodeCid, big.NewInt(0))

	tree, storage, err := stateFactory.NewState(actors)
	require.NoError(t, err)
	require.NotNil(t, tree)
	require.NotNil(t, storage)

	producer := chain.NewMessageProducer(msgFactory)
	require.NoError(t, producer.Transfer(alice, bob, big.NewInt(50)))

	exeCtx := state.NewExecutionContext(1, alice)
	validator := state.NewValidator(stateFactory)

	endState, err := validator.ApplyMessages(exeCtx, tree, storage, producer.Messages())
	require.NoError(t, err)
	require.NotNil(t, endState)

	actorAlice, err := endState.Actor(alice)
	require.NoError(t, err)
	assert.Equal(t, state.AttoFIL(big.NewInt(50)), actorAlice.Balance())

	actorBob, err := endState.Actor(bob)
	require.NoError(t, err)
	assert.Equal(t, state.AttoFIL(big.NewInt(50)), actorBob.Balance())
}

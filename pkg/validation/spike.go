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
	minerOwner, err := state.NewActorAddress([]byte("miner")) // This should really be a SECP address
	require.NoError(t, err)

	actors := make(map[state.Address]state.Actor)
	actors[state.NetworkAddress] = stateFactory.NewActor(state.AccountActorCodeCid, big.NewInt(1000000))
	actors[state.BurntFundsAddress] = stateFactory.NewActor(state.AccountActorCodeCid, big.NewInt(0))
	actors[minerOwner] = stateFactory.NewActor(state.AccountActorCodeCid, big.NewInt(0))
	tree, storage, err := stateFactory.NewState(actors)
	if err != nil {
		t.Fatal(err)
	}

	producer := chain.NewMessageProducer(msgFactory)
	require.NoError(t, producer.Transfer(state.NetworkAddress, state.BurntFundsAddress, big.NewInt(1)))

	context := state.NewExecutionContext(1, minerOwner)
	validator := state.NewValidator(stateFactory)

	endState, err := validator.ApplyMessages(context, tree, storage, producer.Messages())
	require.NoError(t, err)
	require.NotNil(t, endState)

	networkActor, err := endState.Actor(state.NetworkAddress)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(999999), networkActor.Balance())

	burntActor, err := endState.Actor(state.BurntFundsAddress)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), burntActor.Balance())
}

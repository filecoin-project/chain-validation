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
func TryItOut(t *testing.T, msgFactory chain.MessageFactory, stateFactory state.StateFactory, storageFactory state.StorageFactory) {
	actors := make(map[state.Address]state.Actor)
	actors[state.NetworkAddress] = stateFactory.NewActor(state.AccountActorCodeCid, big.NewInt(1000000))
	actors[state.BurntFundsAddress] = stateFactory.NewActor(state.AccountActorCodeCid, big.NewInt(0))
	initState, err := stateFactory.NewState(actors)
	if err != nil {
		t.Fatal(err)
	}

	producer := chain.NewMessageProducer(msgFactory)
	require.NoError(t, producer.Transfer(state.NetworkAddress, state.BurntFundsAddress, big.NewInt(1)))

	validator := state.NewValidator(stateFactory, storageFactory)
	endState, err := validator.ApplyMessages(initState, producer.Messages())
	require.NoError(t, err)

	networkActor, err := endState.Actor(state.NetworkAddress)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(999999), networkActor.Balance())

	burntActor, err := endState.Actor(state.BurntFundsAddress)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), burntActor.Balance())
}

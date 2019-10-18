package suites

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

// A basic example validation test.
// At present this code is verbose and demonstrates the opportunity for helper methods.
func Example(t *testing.T, driver Driver) {
	sugar := NewSugarDrive(t, driver)

	alice := sugar.MakeAccountActor(2000)
	bob := sugar.MakeAccountActor(0)
	miner := sugar.MakeStorageMinerActor(alice.address, []byte{1}, 1024)
	producer := sugar.ProducerWithActors(miner, alice, bob)

	producer.Transfer(alice, bob, 50)
	assert.Equal(t, state.AttoFIL(big.NewInt(1950)), alice.actor.Balance())
	assert.Equal(t, state.AttoFIL(big.NewInt(50)), bob.actor.Balance())
	// This should become non-zero after gas tracking and payments are integrated.
	assert.Equal(t, state.AttoFIL(big.NewInt(0)), miner.actor.Balance())

	producer.Transfer(alice, bob, 100)
	assert.Equal(t, state.AttoFIL(big.NewInt(1850)), alice.actor.Balance())
	assert.Equal(t, state.AttoFIL(big.NewInt(150)), bob.actor.Balance())
	// This should become non-zero after gas tracking and payments are integrated.
	assert.Equal(t, state.AttoFIL(big.NewInt(0)), miner.actor.Balance())

}

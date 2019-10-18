package suites

import (
	"testing"
)

// A basic example validation test.
// At present this code is verbose and demonstrates the opportunity for helper methods.
func Example(t *testing.T, driver Driver) {
	sugar := NewSugarDrive(t, driver)

	alice := sugar.MakeAccountActor(2000)
	bob := sugar.MakeAccountActor(0)
	miner := sugar.MakeStorageMinerActor(alice.address, []byte{1}, 1024)
	producer := sugar.ProducerWithActors(miner, alice, bob)

	producer.MustTransfer(alice, bob, 50)
	producer.MustTransfer(alice, bob, 100)
}

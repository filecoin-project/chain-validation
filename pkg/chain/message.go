package chain

import (
	"github.com/filecoin-project/chain-validation/pkg/state"
)

// MethodID identifies a VM actor method.
// The values here are not intended to match the spec's method IDs, though once implementations
// converge on those we could make it so.
// Integrations should map these method ids to the internal method handle representation.
type MethodID int

// An enumeration of all actor methods which a message could invoke.
// Note that some methods are not intended for direct invocation by account actors, but they are still
// listed here so that the behaviour of attempting to invoke them can be exercised.
const (
	NoMethod MethodID = iota
	InitConstructor
	InitExec
	InitGetActorIDForAddress
	StoragePowerConstructor
	StoragePowerCreateStorageMiner
	StoragePowerUpdatePower
	StoragePowerTotalStorage
	StoragePowerPowerLookup
	StoragePowerIncrementPower
	StoragePowerSuspendMiner
	StorageMarketConstructor
	StorageMarketWithdrawBalance
	StorageMarketAddBalance
	StorageMarketCheckLockedBalance
	StorageMarketPublishStorageDeal
	StorageMarketHandleCronAction
	CronConstructor
	CronTick
	// List not yet complete, pending specification.

	// Provides a value above which integrations can assign their own method identifiers without
	// collision with these "standard" ones.
	MethodCount
)

// MessageFactory creates a concrete, but opaque, message object.
// Integrations should implement this to provide a message value that will be accepted by the
// validation engine.
type MessageFactory interface {
	MakeMessage(from, to state.Address, method MethodID, nonce uint64, value, gasPrice state.AttoFIL, gasLimit state.GasUnit,
		params ...interface{}) (interface{}, error)
}

// MessageProducer presents a convenient API for scripting the creation of long and complex message sequences.
// The created messages are retained for subsequent export or evaluation in a VM.
// Actual message construction is delegated to a `MessageFactory`, and the message are opaque to the producer.
type MessageProducer struct {
	factory       MessageFactory
	accountNonces map[state.Address]uint64
}

// NewMessageProducer creates a new message producer, delegating message creation to `factory`.
func NewMessageProducer(factory MessageFactory) *MessageProducer {
	return &MessageProducer{
		factory:       factory,
		accountNonces: make(map[state.Address]uint64),
	}
}

// Build creates a single message and returns it.
func (mp *MessageProducer) Build(from, to state.Address, method MethodID, nonce uint64, value, gasPrice state.AttoFIL, gasLimit state.GasUnit,
	params ...interface{}) (interface{}, error) {
	fm, err := mp.factory.MakeMessage(from, to, method, nonce, value, gasPrice, gasLimit, params)
	if err != nil {
		return nil, err
	}

	return fm, nil
}

//
// Sugar methods for type-checked construction of specific messages.
//

// Transfer builds a simple value transfer message and returns it.
func (mp *MessageProducer) Transfer(from, to state.Address, value, gasPrice state.AttoFIL, gasLimit state.GasUnit) (interface{}, error) {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(from, to, NoMethod, nonce, value, gasPrice, gasLimit)
}

// InitExec builds a message invoking InitActor.Exec and returns it.
func (mp *MessageProducer) InitExec(from state.Address, value, gasPrice state.AttoFIL, gasLimit state.GasUnit, params ...interface{}) (interface{}, error) {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(state.InitAddress, from, InitExec, nonce, value, gasPrice, gasLimit, params)
}

// StoragePowerCreateStorageMiner builds a message invoking StoragePowerActor.CreateStorageMiner and returns it.
func (mp *MessageProducer) StoragePowerCreateStorageMiner(from state.Address, value, gasPrice state.AttoFIL, gasLimit state.GasUnit,
	owner state.Address, worker state.PubKey, sectorSize state.BytesAmount, peerID state.PeerID) (interface{}, error) {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(state.StorageMarketAddress, from, StoragePowerCreateStorageMiner, nonce, value, gasPrice, gasLimit, owner, worker, sectorSize, peerID)
}

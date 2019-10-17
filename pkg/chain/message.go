package chain

import (
	"github.com/filecoin-project/chain-validation/pkg/state"
)

// MethodID is the type of a VM actor method identifier.
// This is a string for generality at the moment, but should eventually become an integer.
type MethodID int

// An enumeration of all actor methods which a message could invoke.
// The enum values here are not intended to match the spec's method IDs, though once implementations
// converge on those we could make it so.
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
)

// MessageFactory creates a concrete, but opaque, message object.
type MessageFactory interface {
	MakeMessage(from, to state.Address, method MethodID, nonce uint64, value state.AttoFIL,
		params ...interface{}) (interface{}, error)
}

// MessageProducer presents a convenient API for scripting the creation of long and complex message sequences.
// The created messages are retained for subsequent export or evaluation in a VM.
// Actual message construction is delegated to a `MessageFactory`, and the message are opaque to the producer.
type MessageProducer struct {
	factory       MessageFactory
	messages      []interface{}
	accountNonces map[state.Address]uint64
}

// NewMessageProducer creates a new message producer, delegating message creation to `factory`.
func NewMessageProducer(factory MessageFactory) *MessageProducer {
	return &MessageProducer{
		factory:       factory,
		accountNonces: make(map[state.Address]uint64),
	}
}

// Messages returns a slice containing all messages created by the producer.
func (mp *MessageProducer) Messages() []interface{} {
	return mp.messages[:]
}

// Build creates and stores a single message.
func (mp *MessageProducer) Build(from, to state.Address, method MethodID, nonce uint64, value state.AttoFIL,
	params ...interface{}) error {
	fm, err := mp.factory.MakeMessage(from, to, method, nonce, value, params)
	if err != nil {
		return err
	}

	mp.messages = append(mp.messages, fm)
	return nil
}

//
// Sugar methods for type-checked construction of specific messages.
//

func (mp *MessageProducer) Transfer(from, to state.Address, value state.AttoFIL) error {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(from, to, NoMethod, nonce, value)
}

//
// Init actor
//

func (mp *MessageProducer) InitExec(from state.Address, value state.AttoFIL, params ...interface{}) error {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(state.InitAddress, from, InitExec, nonce, value, params)
}

//
// StoragePower actor
//

func (mp *MessageProducer) StoragePowerCreateStorageMiner(from state.Address, value state.AttoFIL, owner state.Address, worker state.PubKey, sectorSize state.BytesAmount, peerID state.PeerID) error {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(state.StorageMarketAddress, from, StoragePowerCreateStorageMiner, nonce, value, owner, worker, sectorSize, peerID)
}

package chain

import (
	"math/big"

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
	factory  MessageFactory
	defaults msgOpts // Note non-pointer reference.

	messages []interface{}
}

// NewMessageProducer creates a new message producer, delegating message creation to `factory`.
func NewMessageProducer(factory MessageFactory, defaultGasLimit state.GasUnit, defaultGasPrice state.AttoFIL) *MessageProducer {
	return &MessageProducer{
		factory: factory,
		defaults: msgOpts{
			gasLimit: defaultGasLimit,
			gasPrice: defaultGasPrice,
		},
	}
}

// Messages returns a slice containing all messages created by the producer.
func (mp *MessageProducer) Messages() []interface{} {
	return mp.messages[:]
}

// msgOpts specifies value and gas parameters for a message, supporting a functional options pattern
// for concise but customizable message construction.
type msgOpts struct {
	value    state.AttoFIL
	gasLimit state.GasUnit
	gasPrice state.AttoFIL
}

// MsgOpt is an option configuring message value or gas parameters.
type MsgOpt func(*msgOpts)

func Value(value uint64) MsgOpt {
	return func(opts *msgOpts) {
		opts.value = big.NewInt(0).SetUint64(value)
	}
}

func GasLimit(limit uint64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasLimit = state.GasUnit(limit)
	}
}

func GasPrice(price uint64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasPrice = big.NewInt(0).SetUint64(price)
	}
}

// Build creates and returns a single message, using default gas parameters unless modified by `opts`.
func (mp *MessageProducer) Build(from, to state.Address, nonce uint64, method MethodID, params []interface{},
	opts ...MsgOpt) (interface{}, error) {
	values := mp.defaults
	for _, opt := range opts {
		opt(&values)
	}

	return mp.BuildFull(from, to, method, nonce, values.value, values.gasLimit, values.gasPrice, params...)
}

// BuildFull creates and returns a single message.
func (mp *MessageProducer) BuildFull(from, to state.Address, method MethodID, nonce uint64, value state.AttoFIL,
	gasLimit state.GasUnit, gasPrice state.AttoFIL, params ...interface{}) (interface{}, error) {
	fm, err := mp.factory.MakeMessage(from, to, method, nonce, value, gasPrice, gasLimit, params...)
	if err != nil {
		return nil, err
	}

	mp.messages = append(mp.messages, fm)
	return fm, nil
}

//
// Sugar methods for type-checked construction of specific messages.
//

// Transfer builds a simple value transfer message and returns it.
func (mp *MessageProducer) Transfer(from, to state.Address, nonce uint64, value uint64, opts ...MsgOpt) (interface{}, error) {
	x := append([]MsgOpt{Value(value)}, opts...)
	return mp.Build(from, to, nonce, NoMethod, noParams, x...)
}

// InitExec builds a message invoking InitActor.Exec and returns it.
func (mp *MessageProducer) InitExec(from state.Address, nonce uint64, params []interface{}, opts ...MsgOpt) (interface{}, error) {
	return mp.Build(state.InitAddress, from, nonce, InitExec, params, opts...)
}

// StoragePowerCreateStorageMiner builds a message invoking StoragePowerActor.CreateStorageMiner and returns it.
func (mp *MessageProducer) StoragePowerCreateStorageMiner(from state.Address, nonce uint64,
	owner state.Address, worker state.PubKey, sectorSize state.BytesAmount, peerID state.PeerID,
	opts ...MsgOpt) (interface{}, error) {
	params := []interface{}{owner, worker, sectorSize, peerID}
	return mp.Build(from, state.StorageMarketAddress, nonce, StoragePowerCreateStorageMiner, params, opts...)
}

var noParams []interface{}


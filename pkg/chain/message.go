package chain

import (
	"math/big"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/chain-validation/pkg/state"
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
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

	StorageMinerUpdatePeerID
	StorageMinerGetOwner
	StorageMinerGetWorkerAddr
	StorageMinerGetPower
	StorageMinerGetPeerID
	StorageMinerGetSectorSize

	PaymentChannelCreate

	GetSectorSize
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
	MakeMessage(from, to address.Address, method MethodID, nonce uint64, value, gasPrice types.AttoFIL, gasLimit types.GasUnit,
		params ...interface{}) (interface{}, error)
	FromSingletonAddress(address state.SingletonActorID) address.Address
	FromActorCodeCid(cod state.ActorCodeID) cid.Cid
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
func NewMessageProducer(factory MessageFactory, defaultGasLimit types.GasUnit, defaultGasPrice types.AttoFIL) *MessageProducer {
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
	value    types.AttoFIL
	gasLimit types.GasUnit
	gasPrice types.AttoFIL
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
		opts.gasLimit = types.GasUnit(limit)
	}
}

func GasPrice(price uint64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasPrice = big.NewInt(0).SetUint64(price)
	}
}

// Build creates and returns a single message, using default gas parameters unless modified by `opts`.
func (mp *MessageProducer) Build(from, to address.Address, nonce uint64, method MethodID, params []interface{},
	opts ...MsgOpt) (interface{}, error) {
	values := mp.defaults
	for _, opt := range opts {
		opt(&values)
	}

	return mp.BuildFull(from, to, method, nonce, values.value, values.gasLimit, values.gasPrice, params...)
}

// BuildFull creates and returns a single message.
func (mp *MessageProducer) BuildFull(from, to address.Address, method MethodID, nonce uint64, value types.AttoFIL,
	gasLimit types.GasUnit, gasPrice types.AttoFIL, params ...interface{}) (interface{}, error) {
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
func (mp *MessageProducer) Transfer(from, to address.Address, nonce uint64, value uint64, opts ...MsgOpt) (interface{}, error) {
	x := append([]MsgOpt{Value(value)}, opts...)
	return mp.Build(from, to, nonce, NoMethod, noParams, x...)
}

// InitExec builds a message invoking InitActor.Exec and returns it.
func (mp *MessageProducer) InitExec(from address.Address, nonce uint64, params []interface{}, opts ...MsgOpt) (interface{}, error) {
	iaAddr := mp.factory.FromSingletonAddress(state.InitAddress)
	return mp.Build(iaAddr, from, nonce, InitExec, params, opts...)
}

//
// Storage Power Actor Methods
//

// StoragePowerCreateStorageMiner builds a message invoking StoragePowerActor.CreateStorageMiner and returns it.
func (mp *MessageProducer) StoragePowerCreateStorageMiner(from address.Address, nonce uint64,
	owner address.Address, worker address.Address, sectorSize types.BytesAmount, peerID types.PeerID,
	opts ...MsgOpt) (interface{}, error) {

	spaAddr := mp.factory.FromSingletonAddress(state.StoragePowerAddress)
	params := []interface{}{owner, worker, sectorSize, peerID}
	return mp.Build(from, spaAddr, nonce, StoragePowerCreateStorageMiner, params, opts...)
}

func (mp *MessageProducer) StoragePowerUpdateStorage(from address.Address, nonce uint64, delta types.BytesAmount, opts ...MsgOpt) (interface{}, error) {
	spaAddr := mp.factory.FromSingletonAddress(state.StoragePowerAddress)
	params := []interface{}{delta}
	return mp.Build(from, spaAddr, nonce, StoragePowerUpdatePower, params, opts...)
}

//
// Storage Miner Actor Methods
//

func (mp *MessageProducer) StorageMinerUpdatePeerID(to, from address.Address, nonce uint64, peerID types.PeerID, opts ...MsgOpt) (interface{}, error) {
	params := []interface{}{peerID}
	return mp.Build(from, to, nonce, StorageMinerUpdatePeerID, params, opts...)
}

func (mp *MessageProducer) StorageMinerGetOwner(to, from address.Address, nonce uint64, opts ...MsgOpt) (interface{}, error) {
	return mp.Build(from, to, nonce, StorageMinerGetOwner, noParams, opts...)
}

func (mp *MessageProducer) StorageMinerGetPower(to, from address.Address, nonce uint64, opts ...MsgOpt) (interface{}, error) {
	return mp.Build(from, to, nonce, StorageMinerGetPower, noParams, opts...)
}

func (mp *MessageProducer) StorageMinerGetWorkerAddr(to, from address.Address, nonce uint64, opts ...MsgOpt) (interface{}, error) {
	return mp.Build(from, to, nonce, StorageMinerGetWorkerAddr, noParams, opts...)
}

func (mp *MessageProducer) StorageMinerGetPeerID(to, from address.Address, nonce uint64, opts ...MsgOpt) (interface{}, error) {
	return mp.Build(from, to, nonce, StorageMinerGetPeerID, noParams, opts...)
}

func (mp *MessageProducer) StorageMinerGetSectorSize(to, from address.Address, nonce uint64, opts ...MsgOpt) (interface{}, error) {
	return mp.Build(from, to, nonce, StorageMinerGetSectorSize, noParams, opts...)
}

//
// Payment Channel Actor Methods
//

func (mp *MessageProducer) PaymentChannelCreate(to, from address.Address, nonce, value uint64, opts ...MsgOpt) (interface{}, error) {
	payChParams := []interface{}{to}
	msgOpt := append([]MsgOpt{Value(value)}, opts...)

	initParams := []interface{}{mp.factory.FromActorCodeCid(state.PaymentChannelActorCodeCid), payChParams}
	return mp.Build(from, mp.factory.FromSingletonAddress(state.InitAddress), nonce, InitExec, initParams, msgOpt...)
}

var noParams []interface{}

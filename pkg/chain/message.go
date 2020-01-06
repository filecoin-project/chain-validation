package chain

import (
	"github.com/filecoin-project/chain-validation/pkg/state"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/initialize"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/multsig"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/paych"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgminr"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgmrkt"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgpwr"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type Message struct {
	// Address of the receiving actor.
	To address.Address
	// Address of the sending actor.
	From address.Address
	// Expected CallSeqNum of the sending actor (only for top-level messages).
	CallSeqNum uint64

	// Amount of value to transfer from sender's to receiver's balance.
	Value types.BigInt

	// Optional method to invoke on receiver, zero for a plain value send.
	Method MethodID
	/// Serialized parameters to the method (if method is non-zero).
	Params []byte

	GasPrice types.BigInt
	GasLimit types.BigInt
}

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
	StorageMarketPublishStorageDeals
	StorageMarketActivateStorageDeals
	StorageMarketComputeDataCommitment

	StorageMinerUpdatePeerID
	StorageMinerGetOwner
	StorageMinerGetWorkerAddr
	StorageMinerGetPower
	StorageMinerGetPeerID
	StorageMinerGetSectorSize

	MultiSigConstructor
	MultiSigPropose
	MultiSigApprove
	MultiSigCancel
	MultiSigClearCompleted
	MultiSigAddSigner
	MultiSigRemoveSigner
	MultiSigSwapSigner
	MultiSigChangeRequirement

	PaymentChannelConstructor
	PaymentChannelUpdate
	PaymentChannelClose
	PaymentChannelCollect
	PaymentChannelGetOwner
	PaymentChannelGetToSend

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
	MakeMessage(from, to address.Address, method MethodID, nonce uint64, value, gasPrice, gasLimit types.BigInt, params []byte) (*Message, error)
}

type ActorInfoMapping interface {
	FromSingletonAddress(address actors.SingletonActorID) address.Address
	FromActorCodeCid(cod actors.ActorCodeID) cid.Cid
}

// MessageProducer presents a convenient API for scripting the creation of long and complex message sequences.
// The created messages are retained for subsequent export or evaluation in a VM.
// Actual message construction is delegated to a `MessageFactory`, and the message are opaque to the producer.
type MessageProducer struct {
	factory   MessageFactory
	actorInfo ActorInfoMapping
	defaults  msgOpts // Note non-pointer reference.

	messages []*Message
}

// NewMessageProducer creates a new message producer, delegating message creation to `factory`.
func NewMessageProducer(factory MessageFactory, ai ActorInfoMapping, defaultGasLimit, defaultGasPrice types.BigInt) *MessageProducer {
	return &MessageProducer{
		factory:   factory,
		actorInfo: ai,
		defaults: msgOpts{
			gasLimit: defaultGasLimit,
			gasPrice: defaultGasPrice,
		},
	}
}

// Messages returns a slice containing all messages created by the producer.
func (mp *MessageProducer) Messages() []*Message {
	return mp.messages
}

// msgOpts specifies value and gas parameters for a message, supporting a functional options pattern
// for concise but customizable message construction.
type msgOpts struct {
	value    types.BigInt
	gasLimit types.BigInt
	gasPrice types.BigInt
}

// MsgOpt is an option configuring message value or gas parameters.
type MsgOpt func(*msgOpts)

func Value(value uint64) MsgOpt {
	return func(opts *msgOpts) {
		opts.value = types.NewInt(value)
	}
}

func BigValue(value types.BigInt) MsgOpt {
	return func(opts *msgOpts) {
		opts.value = value
	}
}

func GasLimit(limit uint64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasLimit = types.NewInt(limit)
	}
}

func GasPrice(price uint64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasPrice = types.NewInt(price)
	}
}

// Build creates and returns a single message, using default gas parameters unless modified by `opts`.
func (mp *MessageProducer) Build(from, to address.Address, nonce uint64, method MethodID, params []byte, opts ...MsgOpt) (*Message, error) {
	values := mp.defaults
	for _, opt := range opts {
		opt(&values)
	}

	return mp.BuildFull(from, to, method, nonce, values.value, values.gasLimit, values.gasPrice, params)
}

// BuildFull creates and returns a single message.
func (mp *MessageProducer) BuildFull(from, to address.Address, method MethodID, nonce uint64, value types.BigInt,
	gasLimit, gasPrice types.BigInt, params []byte) (*Message, error) {
	fm, err := mp.factory.MakeMessage(from, to, method, nonce, value, gasPrice, gasLimit, params)
	if err != nil {
		return nil, err
	}

	mp.messages = append(mp.messages, fm)
	return fm, nil
}

//
// Helper methods until spec defines these
//

func (mp *MessageProducer) SingletonAddress(id actors.SingletonActorID) address.Address {
	return mp.actorInfo.FromSingletonAddress(id)
}

func (mp *MessageProducer) ActorCid(c actors.ActorCodeID) cid.Cid {
	return mp.actorInfo.FromActorCodeCid(c)
}

//
// Sugar methods for type-checked construction of specific messages.
//

// Transfer builds a simple value transfer message and returns it.
func (mp *MessageProducer) Transfer(from, to address.Address, nonce uint64, value uint64, opts ...MsgOpt) (*Message, error) {
	x := append([]MsgOpt{Value(value)}, opts...)
	return mp.Build(from, to, nonce, NoMethod, noParams, x...)
}

// InitExec builds a message invoking InitActor.Exec and returns it.
func (mp *MessageProducer) InitExec(from address.Address, nonce uint64, code actors.ActorCodeID, params []byte, opts ...MsgOpt) (*Message, error) {
	iaAddr := mp.actorInfo.FromSingletonAddress(actors.InitAddress)
	initParams, err := state.Serialize(&initialize.ExecParams{
		Code:   mp.actorInfo.FromActorCodeCid(code),
		Params: params,
	})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, iaAddr, nonce, InitExec, initParams, opts...)
}

//
// Storage Market Actor Methods
//

func (mp *MessageProducer) StorageMarketWithdrawBalance(from address.Address, nonce uint64, balance types.BigInt, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgmrkt.WithdrawBalanceParams{Balance: balance})
	if err != nil {
		return nil, err
	}
	smaddr := mp.actorInfo.FromSingletonAddress(actors.StorageMarketAddress)
	return mp.Build(from, smaddr, nonce, StorageMarketWithdrawBalance, params, opts...)
}

func (mp *MessageProducer) StorageMarketAddBalance(from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	smaddr := mp.actorInfo.FromSingletonAddress(actors.StorageMarketAddress)
	return mp.Build(from, smaddr, nonce, StorageMarketAddBalance, noParams, opts...)
}

func (mp *MessageProducer) StorageMarketPublishStorageDeals(from address.Address, nonce uint64, deals []strgmrkt.StorageDeal, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgmrkt.PublishStorageDealsParams{Deals: deals})
	if err != nil {
		return nil, err
	}
	smaddr := mp.actorInfo.FromSingletonAddress(actors.StorageMarketAddress)
	return mp.Build(from, smaddr, nonce, StorageMarketPublishStorageDeals, params, opts...)
}

func (mp *MessageProducer) StorageMarketActivateStorageDeals(from address.Address, nonce uint64, dealIDs []uint64, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgmrkt.ActivateStorageDealsParams{Deals: dealIDs})
	if err != nil {
		return nil, err
	}
	smaddr := mp.actorInfo.FromSingletonAddress(actors.StorageMarketAddress)
	return mp.Build(from, smaddr, nonce, StorageMarketActivateStorageDeals, params, opts...)
}

func (mp *MessageProducer) StorageMarketComputeDataCommitment(from address.Address, nonce uint64, sectorSize uint64, dealIDs []uint64, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgmrkt.ComputeDataCommitmentParams{
		DealIDs:    dealIDs,
		SectorSize: sectorSize,
	})
	if err != nil {
		return nil, err
	}
	smaddr := mp.actorInfo.FromSingletonAddress(actors.StorageMarketAddress)
	return mp.Build(from, smaddr, nonce, StorageMarketComputeDataCommitment, params, opts...)
}

//
// Storage Power Actor Methods
//

// StoragePowerCreateStorageMiner builds a message invoking StoragePowerActor.CreateStorageMiner and returns it.
func (mp *MessageProducer) StoragePowerCreateStorageMiner(from address.Address, nonce uint64,
	owner address.Address, worker address.Address, sectorSize uint64, peerID peer.ID,
	opts ...MsgOpt) (*Message, error) {

	spaAddr := mp.actorInfo.FromSingletonAddress(actors.StoragePowerAddress)
	params, err := state.Serialize(&strgpwr.CreateStorageMinerParams{
		Owner:      owner,
		Worker:     worker,
		SectorSize: sectorSize,
		PeerID:     peerID,
	})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, spaAddr, nonce, StoragePowerCreateStorageMiner, params, opts...)
}

func (mp *MessageProducer) StoragePowerUpdateStorage(from address.Address, nonce uint64, delta types.BigInt, nextppEnd, previousppEnd uint64, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgpwr.UpdateStorageParams{
		Delta:                    delta,
		NextProvingPeriodEnd:     nextppEnd,
		PreviousProvingPeriodEnd: previousppEnd,
	})
	if err != nil {
		return nil, err
	}
	spaAddr := mp.actorInfo.FromSingletonAddress(actors.StoragePowerAddress)
	return mp.Build(from, spaAddr, nonce, StoragePowerUpdatePower, params, opts...)
}

func (mp *MessageProducer) StoragePowerPledgeCollateralForSize(from address.Address, nonce uint64, size types.BigInt, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgpwr.PledgeCollateralParams{Size: size})
	if err != nil {
		return nil, err
	}
	spaAddr := mp.actorInfo.FromSingletonAddress(actors.StoragePowerAddress)
	return mp.Build(from, spaAddr, nonce, StoragePowerUpdatePower, params, opts...)
}

func (mp *MessageProducer) StoragePowerLookupPower(from address.Address, nonce uint64, miner address.Address, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgpwr.PowerLookupParams{Miner: miner})
	if err != nil {
		return nil, err
	}
	spaAddr := mp.actorInfo.FromSingletonAddress(actors.StoragePowerAddress)
	return mp.Build(from, spaAddr, nonce, StoragePowerUpdatePower, params, opts...)
}

//
// Storage Miner Actor Methods
//

func (mp *MessageProducer) StorageMinerUpdatePeerID(to, from address.Address, nonce uint64, peerID peer.ID, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgminr.UpdatePeerIDParams{PeerID: peerID})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, to, nonce, StorageMinerUpdatePeerID, params, opts...)
}

func (mp *MessageProducer) StorageMinerGetOwner(to, from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, StorageMinerGetOwner, noParams, opts...)
}

func (mp *MessageProducer) StorageMinerGetPower(to, from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, StorageMinerGetPower, noParams, opts...)
}

func (mp *MessageProducer) StorageMinerGetWorkerAddr(to, from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, StorageMinerGetWorkerAddr, noParams, opts...)
}

func (mp *MessageProducer) StorageMinerGetPeerID(to, from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, StorageMinerGetPeerID, noParams, opts...)
}

func (mp *MessageProducer) StorageMinerGetSectorSize(to, from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, StorageMinerGetSectorSize, noParams, opts...)
}

//
// Multi Signature Actor Methods
//

func (mp *MessageProducer) MultiSigPropose(to, from address.Address, nonce uint64, proposeTo address.Address, proposeValue types.BigInt, proposeMethod uint64, proposeParams []byte, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&multsig.MultiSigProposeParams{
		To:     proposeTo,
		Value:  proposeValue,
		Method: proposeMethod,
		Params: proposeParams,
	})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, to, nonce, MultiSigPropose, params, opts...)
}

func (mp *MessageProducer) MultiSigApprove(to, from address.Address, nonce uint64, txID uint64, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&multsig.MultiSigTxID{TxID: txID})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, to, nonce, MultiSigApprove, params, opts...)
}

func (mp *MessageProducer) MultiSigCancel(to, from address.Address, nonce uint64, txID uint64, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&multsig.MultiSigTxID{TxID: txID})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, to, nonce, MultiSigCancel, params, opts...)
}

func (mp *MessageProducer) MultiSigAddSigner(to, from address.Address, nonce uint64, signer address.Address, increase bool, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&multsig.MultiSigAddSignerParam{
		Signer:   signer,
		Increase: increase,
	})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, to, nonce, MultiSigAddSigner, params, opts...)
}

func (mp *MessageProducer) MultiSigRemoveSigner(to, from address.Address, nonce uint64, signer address.Address, decrease bool, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&multsig.MultiSigRemoveSignerParam{
		Signer:   signer,
		Decrease: decrease,
	})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, to, nonce, MultiSigRemoveSigner, params, opts...)
}

func (mp *MessageProducer) MultiSigSwapSigner(to, from address.Address, nonce uint64, swapFrom, swapTo address.Address, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&multsig.MultiSigSwapSignerParams{
		From: swapFrom,
		To:   swapTo,
	})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, to, nonce, MultiSigSwapSigner, params, opts...)
}

func (mp *MessageProducer) MultiSigChangeRequirement(to, from address.Address, nonce uint64, req uint64, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&multsig.MultiSigChangeReqParams{Req: req})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, to, nonce, MultiSigChangeRequirement, params, opts...)
}

//
// Payment Channel Actor Methods
//

func (mp *MessageProducer) PaychUpdateChannelState(to, from address.Address, nonce uint64, sv types.SignedVoucher, secret, proof []byte, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&paych.PaymentChannelUpdateParams{
		Sv:     sv,
		Secret: secret,
		Proof:  proof,
	})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, to, nonce, PaymentChannelUpdate, params, opts...)
}

func (mp *MessageProducer) PaychClose(to, from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, PaymentChannelClose, noParams, opts...)
}
func (mp *MessageProducer) PaychCollect(to, from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, PaymentChannelCollect, noParams, opts...)
}
func (mp *MessageProducer) PaychGetOwner(to, from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, PaymentChannelGetOwner, noParams, opts...)
}
func (mp *MessageProducer) PaychGetToSend(to, from address.Address, nonce uint64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, PaymentChannelGetToSend, noParams, opts...)
}

var noParams []byte

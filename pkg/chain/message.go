package chain

import (
	address "github.com/filecoin-project/go-address"
	cid "github.com/ipfs/go-cid"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	init_spec "github.com/filecoin-project/specs-actors/actors/builtin/init"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"

	state "github.com/filecoin-project/chain-validation/pkg/state"
)

type Message struct {
	// Address of the receiving actor.
	To address.Address
	// Address of the sending actor.
	From address.Address
	// Expected CallSeqNum of the sending actor (only for top-level messages).
	CallSeqNum int64

	// Amount of value to transfer from sender's to receiver's balance.
	Value big_spec.Int

	// Optional method to invoke on receiver, zero for a plain value send.
	Method abi_spec.MethodNum
	/// Serialized parameters to the method (if method is non-zero).
	Params []byte

	GasPrice big_spec.Int
	GasLimit big_spec.Int
}

type MessageFactory struct{}

func (b *MessageFactory) MakeMessage(to, from address.Address, method abi_spec.MethodNum, callSeq int64, value, gasPrice, gasLimit big_spec.Int, params []byte) *Message {
	return &Message{
		To:         to,
		From:       from,
		CallSeqNum: callSeq,
		Value:      value,
		Method:     method,
		Params:     params,
		GasPrice:   gasPrice,
		GasLimit:   gasLimit,
	}
}

// MessageProducer presents a convenient API for scripting the creation of long and complex message sequences.
// The created messages are retained for subsequent export or evaluation in a VM.
// Actual message construction is delegated to a `MessageFactory`.
type MessageProducer struct {
	factory  *MessageFactory
	defaults msgOpts // Note non-pointer reference.

	messages []*Message
}

// NewMessageProducer creates a new message producer, delegating message creation to `factory`.
func NewMessageProducer(defaultGasLimit, defaultGasPrice big_spec.Int) *MessageProducer {
	return &MessageProducer{
		factory: &MessageFactory{},
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
	nonce    int64
	value    big_spec.Int
	gasLimit big_spec.Int
	gasPrice big_spec.Int
}

// MsgOpt is an option configuring message value or gas parameters.
type MsgOpt func(*msgOpts)

func Value(value big_spec.Int) MsgOpt {
	return func(opts *msgOpts) {
		opts.value = value
	}
}

func Nonce(n int64) MsgOpt {
	return func(opts *msgOpts) {
		opts.nonce = n
	}
}

func GasLimit(limit int64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasLimit = big_spec.NewInt(limit)
	}
}

func GasPrice(price int64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasPrice = big_spec.NewInt(price)
	}
}

// BuildFull creates and returns a single message.
func (mp *MessageProducer) BuildFull(to, from address.Address, method abi_spec.MethodNum, nonce int64, value, gasLimit, gasPrice big_spec.Int, params []byte) *Message {
	fm := mp.factory.MakeMessage(to, from, method, nonce, value, gasPrice, gasLimit, params)
	mp.messages = append(mp.messages, fm)
	return fm
}

// Build creates and returns a single message, using default gas parameters unless modified by `opts`.
func (mp *MessageProducer) Build(to, from address.Address, method abi_spec.MethodNum, params []byte, opts ...MsgOpt) *Message {
	values := mp.defaults
	for _, opt := range opts {
		opt(&values)
	}

	return mp.BuildFull(to, from, method, values.nonce, values.value, values.gasLimit, values.gasPrice, params)
}

//
// Sugar methods for type-checked construction of specific messages.
//

// Transfer builds a simple value transfer message and returns it.
func (mp *MessageProducer) Transfer(to, from address.Address, opts ...MsgOpt) *Message {
	return mp.Build(to, from, builtin_spec.MethodSend, noParams, opts...)
}

//
// Init Actor Methods
//

// TODO add the rest of the actor methods

// InitExec builds a message invoking InitActor.Exec and returns it.
func (mp *MessageProducer) InitExec(from address.Address, code cid.Cid, params []byte, opts ...MsgOpt) (*Message, error) {
	initParams, err := state.Serialize(&init_spec.ExecParams{CodeCID: code, ConstructorParams: params})
	if err != nil {
		return nil, err
	}
	return mp.Build(builtin_spec.InitActorAddr, from, builtin_spec.MethodsInit.Exec, initParams, opts...), nil
}

//
// Storage Market Actor Methods
//
/*
func (mp *MessageProducer) StorageMarketWithdrawBalance(from address.Address, nonce int64, balance types.BigInt, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgmrkt.WithdrawBalanceParams{Balance: balance})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, builtin_spec.StorageMarketActorAddr, nonce, builtin_spec.Method_StorageMarketActor_WithdrawBalance, params, opts...)
}

func (mp *MessageProducer) StorageMarketAddBalance(from address.Address, nonce int64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, builtin_spec.StorageMarketActorAddr, nonce, builtin_spec.Method_StorageMarketActor_AddBalance, noParams, opts...)
}

func (mp *MessageProducer) StorageMarketPublishStorageDeals(from address.Address, nonce int64, deals []strgmrkt.StorageDeal, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgmrkt.PublishStorageDealsParams{Deals: deals})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, builtin_spec.StorageMarketActorAddr, nonce, builtin_spec.Method_StorageMarketActor_PublishStorageDeals, params, opts...)
}
*/

//
// Storage Power Actor Methods
//
// TODO add the rest of actor methods

/*
// StoragePowerCreateStorageMiner builds a message invoking StoragePowerActor.CreateStorageMiner and returns it.
func (mp *MessageProducer) StoragePowerCreateStorageMiner(from address.Address, nonce int64,
	owner address.Address, worker address.Address, sectorSize uint64, peerID peer.ID,
	opts ...MsgOpt) (*Message, error) {

	params, err := state.Serialize(&strgpwr.CreateStorageMinerParams{
		Owner:      owner,
		Worker:     worker,
		SectorSize: sectorSize,
		PeerID:     peerID,
	})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, builtin_spec.StoragePowerActorAddr, nonce, builtin_spec.Method_StoragePowerActor_CreateMiner, params, opts...)
}

func (mp *MessageProducer) StoragePowerPledgeCollateralForSize(from address.Address, nonce int64, size types.BigInt, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgpwr.PledgeCollateralParams{Size: size})
	if err != nil {
		return nil, err
	}
	// TODO verify this is the right method and param set
	return mp.Build(from, builtin_spec.StoragePowerActorAddr, nonce, builtin_spec.Method_StoragePowerActor_GetMinerUnmetPledgeCollateralRequirement, params, opts...)
}

func (mp *MessageProducer) StoragePowerLookupPower(from address.Address, nonce int64, miner address.Address, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&strgpwr.PowerLookupParams{Miner: miner})
	if err != nil {
		return nil, err
	}
	return mp.Build(from, builtin_spec.StoragePowerActorAddr, nonce, builtin_spec.Method_StoragePowerActor_GetMinerConsensusPower, params, opts...)
}
*/

/*
//
// Storage Miner Actor Methods
//
// TODO add the rest of actor methods

func (mp *MessageProducer) StorageMinerGetOwner(to, from address.Address, nonce int64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, builtin_spec.Method_StorageMinerActor_GetOwnerAddr, noParams, opts...)
}

func (mp *MessageProducer) StorageMinerGetWorkerAddr(to, from address.Address, nonce int64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, builtin_spec.Method_StorageMinerActor_GetWorkerAddr, noParams, opts...)
}
*/

//
// Multi Signature Actor Methods
//

func (mp *MessageProducer) MultiSigPropose(to, from address.Address, params multisig_spec.ProposeParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Propose, ser, opts...), nil
}

func (mp *MessageProducer) MultiSigApprove(to, from address.Address, txID multisig_spec.TxnID, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&multisig_spec.TxnIDParams{ID: txID})
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Approve, ser, opts...), nil
}

func (mp *MessageProducer) MultiSigCancel(to, from address.Address, txID multisig_spec.TxnID, opts ...MsgOpt) (*Message, error) {
	params, err := state.Serialize(&multisig_spec.TxnIDParams{ID: txID})
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Cancel, params, opts...), nil
}

func (mp *MessageProducer) MultiSigAddSigner(to, from address.Address, params multisig_spec.AddSignerParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.AddSigner, ser, opts...), nil
}

func (mp *MessageProducer) MultiSigRemoveSigner(to, from address.Address, params multisig_spec.RemoveSignerParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.RemoveSigner, ser, opts...), nil
}

func (mp *MessageProducer) MultiSigSwapSigner(to, from address.Address, params multisig_spec.SwapSignerParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.SwapSigner, ser, opts...), nil
}

func (mp *MessageProducer) MultiSigChangeApprovalsThreshold(to, from address.Address, threshold int64, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&multisig_spec.ChangeNumApprovalsThresholdParams{NewThreshold: threshold})
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.ChangeNumApprovalsThreshold, ser, opts...), nil
}

//
// Payment Channel Actor Methods
//
// TODO add these methods when the spec defines them...
/*
func (mp *MessageProducer) PaychUpdateChannelState(to, from address.Address, nonce int64, sv types.SignedVoucher, secret, proof []byte, opts ...MsgOpt) (*Message, error) {
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

func (mp *MessageProducer) PaychClose(to, from address.Address, nonce int64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, PaymentChannelClose, noParams, opts...)
}
func (mp *MessageProducer) PaychCollect(to, from address.Address, nonce int64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, PaymentChannelCollect, noParams, opts...)
}
func (mp *MessageProducer) PaychGetOwner(to, from address.Address, nonce int64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, PaymentChannelGetOwner, noParams, opts...)
}
func (mp *MessageProducer) PaychGetToSend(to, from address.Address, nonce int64, opts ...MsgOpt) (*Message, error) {
	return mp.Build(from, to, nonce, PaymentChannelGetToSend, noParams, opts...)
}
*/

var noParams []byte

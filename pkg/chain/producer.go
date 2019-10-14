package chain

import (
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/filecoin-project/chain-validation/pkg/state"
)


// MessageFactory creates a concrete, but opaque, message object.
type MessageFactory interface {
	MakeMessage(from, to state.Address, method state.MethodID, nonce uint64, value state.AttoFIL, params ...interface{}) (interface{}, error)
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

// Build creates and store a single message.
func (mp *MessageProducer) Build(from, to state.Address, method state.MethodID, nonce uint64, value state.AttoFIL, params ...interface{}) error {
	fm, err := mp.factory.MakeMessage(from, to, method, nonce, value, params)
	if err != nil {
		return err
	}

	mp.messages = append(mp.messages, fm)
	return nil
}

//
// Sugar
//

func (mp *MessageProducer) Transfer(from, to state.Address, value state.AttoFIL) error {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(from, to, state.MethodID(""), nonce, value)
}

//
// StorageMarket
//

func (mp *MessageProducer) CreateStorageMiner(from state.Address, collateral state.AttoFIL, sectorSize state.BytesAmount, peerID peer.ID) error {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(state.StorageMarketAddress, from, state.MethodID("createStorageMiner"), nonce, collateral, sectorSize, peerID)
}

func (mp *MessageProducer) UpdateStorage(from state.Address, value state.AttoFIL, delta state.BytesAmount) error {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(state.StorageMarketAddress, from, state.MethodID("updateStorage"), nonce, value, delta)
}

func (mp *MessageProducer) GetTotalStorage(from state.Address, value state.AttoFIL) error {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(state.StorageMarketAddress, from, state.MethodID("getTotalStorage"), nonce, value)
}

func (mp *MessageProducer) GetProofMode(from state.Address, value state.AttoFIL) error {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(state.StorageMarketAddress, from, state.MethodID("getProofsMode"), nonce, value)
}

func (mp *MessageProducer) GetLateMiners(from state.Address, value state.AttoFIL) error {
	nonce := mp.accountNonces[from]
	mp.accountNonces[from]++
	return mp.Build(state.StorageMarketAddress, from, state.MethodID("getLateMiners"), nonce, value)
}

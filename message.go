package chainvalidation

import (
	"math/big"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/filecoin-project/chain-validation/address"
)

// "as a principal this library should aim for abstractions"

//
// Types shared with Driver
//
type (
	BytesAmount *big.Int
	AttoFIL     *big.Int
	GasUnit     uint64
)

//
// Methods
//
type MethodID string

var (
	MethodSend               = MethodID("")
	MethodCreateStorageMiner = MethodID("createStorageMiner")
	MethodUpdateStorage      = MethodID("updateStorage")
	MethodGetTotalStorage    = MethodID("getTotalStorage")
	MethodGetProofsMode      = MethodID("getProofsMode")
	MethodGetLateMiners      = MethodID("getLateMiners")
)

type MessageFactory interface {
	MakeMessage(to, from address.Address, method MethodID, nonce uint64, value AttoFIL, params ...interface{}) (interface{}, error)
}

type MessageProducer struct {
	messages   []interface{}
	actorNonce map[address.Address]uint64
	factory    MessageFactory

	// TODO could track the nonce and other common params to the Build method
}

//
// Flour
//

func (mp *MessageProducer) Messages() []interface{} {
	return mp.messages
}

func (mp *MessageProducer) Build(to, from address.Address, method MethodID, nonce uint64, value AttoFIL, params ...interface{}) error {
	fm, err := mp.factory.MakeMessage(to, from, method, nonce, value, params)
	if err != nil {
		return err
	}

	mp.messages = append(mp.messages, fm)
	return nil
}

//
// Sugar
//

func (mp *MessageProducer) Transfer(to, from address.Address, value AttoFIL) error {
	nonce := mp.actorNonce[from]
	mp.actorNonce[from]++
	return mp.Build(to, from, MethodSend, nonce, value)
}

//
// StorageMarket
//

func (mp *MessageProducer) CreateStorageMiner(from address.Address, collateral AttoFIL, sectorSize BytesAmount, peerID peer.ID) error {
	nonce := mp.actorNonce[from]
	mp.actorNonce[from]++
	return mp.Build(address.StorageMarketAddress, from, MethodCreateStorageMiner, nonce, collateral, sectorSize, peerID)
}

func (mp *MessageProducer) UpdateStorage(from address.Address, value AttoFIL, delta BytesAmount) error {
	nonce := mp.actorNonce[from]
	mp.actorNonce[from]++
	return mp.Build(address.StorageMarketAddress, from, MethodUpdateStorage, nonce, value, delta)
}

func (mp *MessageProducer) GetTotalStorage(from address.Address, value AttoFIL) error {
	nonce := mp.actorNonce[from]
	mp.actorNonce[from]++
	return mp.Build(address.StorageMarketAddress, from, MethodGetTotalStorage, nonce, value)
}

func (mp *MessageProducer) GetProofMode(from address.Address, value AttoFIL) error {
	nonce := mp.actorNonce[from]
	mp.actorNonce[from]++
	return mp.Build(address.StorageMarketAddress, from, MethodGetProofsMode, nonce, value)
}

func (mp *MessageProducer) GetLateMiners(from address.Address, value AttoFIL) error {
	nonce := mp.actorNonce[from]
	mp.actorNonce[from]++
	return mp.Build(address.StorageMarketAddress, from, MethodGetLateMiners, nonce, value)
}

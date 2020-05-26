package types

import (
	"bytes"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	crypto "github.com/filecoin-project/specs-actors/actors/crypto"
	block "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

type Message struct {
	// Address of the receiving actor.
	To address.Address
	// Address of the sending actor.
	From address.Address
	// Expected CallSeqNum of the sending actor (only for top-level messages).
	CallSeqNum uint64

	// Amount of value to transfer from sender's to receiver's balance.
	Value big.Int

	GasPrice big.Int
	GasLimit int64

	// Optional method to invoke on receiver, zero for a plain value send.
	Method abi.MethodNum
	/// Serialized parameters to the method (if method is non-zero).
	Params []byte
}

func (m *Message) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := m.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Message) Cid() cid.Cid {
	data, err := m.Serialize()
	if err != nil {
		panic(err)
	}

	pref := cid.NewPrefixV1(cid.DagCBOR, multihash.BLAKE2B_MIN+31)
	c, err := pref.Sum(data)
	if err != nil {
		panic(err)
	}

	blk, err := block.NewBlockWithCid(data, c)
	if err != nil {
		panic(err)
	}
	return blk.Cid()

}

type SignedMessage struct {
	Message   Message
	Signature crypto.Signature
}

func (sm *SignedMessage) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := sm.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (sm *SignedMessage) Cid() cid.Cid {
	if sm.Signature.Type == crypto.SigTypeBLS {
		return sm.Message.Cid()
	}

	data, err := sm.Serialize()
	if err != nil {
		panic(err)
	}

	pref := cid.NewPrefixV1(cid.DagCBOR, multihash.BLAKE2B_MIN+31)
	c, err := pref.Sum(data)
	if err != nil {
		panic(err)
	}

	blk, err := block.NewBlockWithCid(data, c)
	if err != nil {
		panic(err)
	}
	return blk.Cid()
}

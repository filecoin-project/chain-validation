package types

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	crypto "github.com/filecoin-project/specs-actors/actors/crypto"
)

type Message struct {
	// Address of the receiving actor.
	To address.Address
	// Address of the sending actor.
	From address.Address
	// Expected CallSeqNum of the sending actor (only for top-level messages).
	CallSeqNum int64

	// Amount of value to transfer from sender's to receiver's balance.
	Value big.Int

	// Optional method to invoke on receiver, zero for a plain value send.
	Method abi.MethodNum
	/// Serialized parameters to the method (if method is non-zero).
	Params []byte

	GasPrice big.Int
	GasLimit int64
}

type SignedMessage struct {
	Message   Message
	Signature crypto.Signature
}

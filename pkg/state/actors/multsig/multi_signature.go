package multsig

import (
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type MultiSigActorState struct {
	Signers  []address.Address
	Required uint64
	NextTxID uint64

	InitialBalance types.BigInt
	StartingBlock  uint64
	UnlockDuration uint64

	//TODO: make this map/sharray/whatever
	Transactions []MTransaction
}

type MTransaction struct {
	Created uint64 // NOT USED ??
	TxID    uint64

	To     address.Address
	Value  types.BigInt
	Method uint64
	Params []byte

	Approved []address.Address
	Complete bool
	Canceled bool
	RetCode  uint64
}

type MultiSigConstructorParams struct {
	Signers        []address.Address
	Required       uint64
	UnlockDuration uint64
}

type MultiSigProposeParams struct {
	To     address.Address
	Value  types.BigInt
	Method uint64
	Params []byte
}

type MultiSigAddSignerParam struct {
	Signer   address.Address
	Increase bool
}

type MultiSigRemoveSignerParam struct {
	Signer   address.Address
	Decrease bool
}

type MultiSigSwapSignerParams struct {
	From address.Address
	To   address.Address
}

type MultiSigChangeReqParams struct {
	Req uint64
}

type MultiSigTxID struct {
	TxID uint64
}

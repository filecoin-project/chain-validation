package types

import "github.com/filecoin-project/chain-validation/pkg/state/address"

type SignedVoucher struct {
	TimeLock       uint64
	SecretPreimage []byte
	Extra          *ModVerifyParams
	Lane           uint64
	Nonce          uint64
	Amount         BigInt
	MinCloseHeight uint64

	Merges []Merge

	Signature *Signature
}

type Merge struct {
	Lane  uint64
	Nonce uint64
}

type ModVerifyParams struct {
	Actor  address.Address
	Method uint64
	Data   []byte
}

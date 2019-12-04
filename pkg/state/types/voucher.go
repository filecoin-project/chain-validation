package types

import (
	"bytes"

	"github.com/filecoin-project/chain-validation/pkg/state/address"
)

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

func (sv *SignedVoucher) SigningBytes() ([]byte, error) {
	osv := *sv
	osv.Signature = nil

	buf := new(bytes.Buffer)
	if err := osv.MarshalCBOR(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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

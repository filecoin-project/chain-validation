package types

import (
	"math/big"

	"github.com/filecoin-project/go-leb128"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/pkg/errors"

	"github.com/filecoin-project/chain-validation/pkg/state/address"
)

// Type aliases for state values and message method parameters.
type (
	BytesAmount *big.Int
	AttoFIL     *big.Int

	GasUnit uint64

	PubKey []byte
	PeerID string
)

// Given a slice of the above types encode them to CBOR byte array.
func EncodeValues(params ...interface{}) ([]byte, error) {
	if len(params) == 0 {
		return []byte{}, nil
	}
	var arr []interface{}
	for i, p := range params {
		bs, err := EncodeValue(p)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed at index: %d", i)
		}
		arr = append(arr, bs)
	}
	return cbor.DumpObject(arr)
}

// Given an above type encode it to CBOR.
func EncodeValue(p interface{}) (interface{}, error) {
	switch v := p.(type) {
	case address.Address:
		return v, nil
	case AttoFIL:
		return (*v).Bytes(), nil
	case BytesAmount:
		return (*v).Uint64(), nil
	case GasUnit:
		return leb128.FromUInt64(uint64(v)), nil
	case uint64:
		return v, nil
	case PubKey:
		return v, nil
	case PeerID:
		return v, nil
	case cid.Cid:
		return v, nil
	case []interface{}:
		return EncodeValues(v...)
	default:
		return []byte{}, errors.Errorf("invalid type: %T", p)
	}
}

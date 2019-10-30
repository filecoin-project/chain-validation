package state

import (
	"math/big"

	"github.com/filecoin-project/go-leb128"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/pkg/errors"
)

// Type aliases for state values and message method parameters.
type (
	BytesAmount *big.Int
	AttoFIL *big.Int

	GasUnit uint64

	PubKey []byte
	PeerID []byte
)

// Given a slice of the above types encode them to CBOR byte array.
func EncodeValues(params ...interface{}) ([]byte, error) {
	if len(params) == 0 {
		return []byte{}, nil
	}
	var arr [][]byte
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
func EncodeValue(p interface{}) ([]byte, error) {
	switch v := p.(type) {
	case Address:
		return []byte(v), nil
	case AttoFIL:
		return (*v).Bytes(), nil
	case BytesAmount:
		return (*v).Bytes(), nil
	case GasUnit:
		return leb128.FromUInt64(uint64(v)), nil
	case uint64:
		return leb128.FromUInt64(v), nil
	case PubKey:
		return v, nil
	case PeerID:
		return v, nil
	default:
		return []byte{}, errors.Errorf("invalid type: %T", p)
	}
}


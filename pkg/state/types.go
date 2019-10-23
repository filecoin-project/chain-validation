package state

import (
	"github.com/filecoin-project/go-leb128"
	"github.com/pkg/errors"
	"math/big"

	cbor "github.com/ipfs/go-ipld-cbor"
)

// Type aliases for state values and message method parameters.
type (
	BytesAmount *big.Int
	AttoFIL *big.Int

	GasUnit uint64

	PubKey []byte
	PeerID []byte
)

// Given a slice of the above types convert them to CBOR byte array.
func EncodeValues(params ...interface{}) ([]byte, error) {
	if len(params) == 0 {
		return []byte{}, nil
	}
	var arr [][]byte
	for i, p := range params {
		switch v := p.(type) {
		case Address:
			arr = append(arr, []byte(v))
		case AttoFIL:
			arr = append(arr, (*v).Bytes())
		case BytesAmount:
			arr = append(arr, (*v).Bytes())
		case GasUnit:
			arr = append(arr, leb128.FromUInt64(uint64(v)))
		case PubKey:
			arr = append(arr, v)
		case PeerID:
			arr = append(arr, v)
		default:
			return []byte{}, errors.Errorf("invalid type: %T at index: %d", p, i)
		}
	}
	return cbor.DumpObject(arr)
}

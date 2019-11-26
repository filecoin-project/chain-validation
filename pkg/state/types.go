package state

import (
	"bytes"
	"fmt"
	"github.com/filecoin-project/chain-validation/pkg/state/types"

	"github.com/filecoin-project/go-leb128"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/pkg/errors"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/pkg/state/address"
)

// Type aliases for state values and message method parameters.
type (
	GasUnit uint64

	PubKey []byte
	PeerID string
)

func Serialize(i cbg.CBORMarshaler) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := i.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Deserialize(b []byte, out interface{}) error {
	um, ok := out.(cbg.CBORUnmarshaler)
	if !ok {
		return fmt.Errorf("type %T does not implement UnmarshalCBOR", out)
	}
	return um.UnmarshalCBOR(bytes.NewReader(b))
}

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
	case types.BigInt:
		return v, nil
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

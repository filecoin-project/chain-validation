package types

import (
	"fmt"
	"io"
	"math/big"

	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/polydawn/refmt/obj/atlas"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func init() {
	cbor.RegisterCborType(atlas.BuildEntry(BigInt{}).Transform().
		TransformMarshal(atlas.MakeMarshalTransformFunc(
			func(i BigInt) ([]byte, error) {
				return i.cborBytes(), nil
			})).
		TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
			func(x []byte) (BigInt, error) {
				return fromCborBytes(x)
			})).
		Complete())
}

const BigIntMaxSerializedLen = 128

type BigInt struct {
	*big.Int
}

var EmptyInt = BigInt{}

func NewInt(i uint64) BigInt {
	return BigInt{big.NewInt(0).SetUint64(i)}
}

func (bi *BigInt) cborBytes() []byte {
	if bi.Int == nil {
		return []byte{}
	}

	switch {
	case bi.Sign() > 0:
		return append([]byte{0}, bi.Bytes()...)
	case bi.Sign() < 0:
		return append([]byte{1}, bi.Bytes()...)
	default: //  bi.Sign() == 0:
		return []byte{}
	}
}

func fromCborBytes(buf []byte) (BigInt, error) {
	if len(buf) == 0 {
		return NewInt(0), nil
	}

	var negative bool
	switch buf[0] {
	case 0:
		negative = false
	case 1:
		negative = true
	default:
		return EmptyInt, fmt.Errorf("big int prefix should be either 0 or 1, got %d", buf[0])
	}

	i := big.NewInt(0).SetBytes(buf[1:])
	if negative {
		i.Neg(i)
	}

	return BigInt{i}, nil
}

func (bi *BigInt) MarshalCBOR(w io.Writer) error {
	if bi.Int == nil {
		zero := NewInt(0)
		return zero.MarshalCBOR(w)
	}

	enc := bi.cborBytes()

	header := cbg.CborEncodeMajorType(cbg.MajByteString, uint64(len(enc)))
	if _, err := w.Write(header); err != nil {
		return err
	}

	if _, err := w.Write(enc); err != nil {
		return err
	}

	return nil
}

func (bi *BigInt) UnmarshalCBOR(br io.Reader) error {
	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if maj != cbg.MajByteString {
		return fmt.Errorf("cbor input for fil big int was not a byte string (%x)", maj)
	}

	if extra == 0 {
		bi.Int = big.NewInt(0)
		return nil
	}

	if extra > BigIntMaxSerializedLen {
		return fmt.Errorf("big integer byte array too long")
	}

	buf := make([]byte, extra)
	if _, err := io.ReadFull(br, buf); err != nil {
		return err
	}

	i, err := fromCborBytes(buf)
	if err != nil {
		return err
	}

	*bi = i

	return nil
}

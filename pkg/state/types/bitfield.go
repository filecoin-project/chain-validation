package types

import (
	"fmt"
	"io"
	"sort"

	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/chain-validation/extern/rleplus"
)

type BitField struct {
	bits map[uint64]struct{}
}

func NewBitField() BitField {
	return BitField{bits: make(map[uint64]struct{})}
}

func BitFieldFromSet(setBits []uint64) BitField {
	res := BitField{bits: make(map[uint64]struct{})}
	for _, b := range setBits {
		res.bits[b] = struct{}{}
	}
	return res
}

// Set ...s bit in the BitField
func (bf BitField) Set(bit uint64) {
	bf.bits[bit] = struct{}{}
}

// Clear ...s bit in the BitField
func (bf BitField) Clear(bit uint64) {
	delete(bf.bits, bit)
}

// Has checkes if bit is set in the BitField
func (bf BitField) Has(bit uint64) bool {
	_, ok := bf.bits[bit]
	return ok
}

// All returns all set bits, in random order
func (bf BitField) All() []uint64 {
	res := make([]uint64, 0, len(bf.bits))
	for i := range bf.bits {
		res = append(res, i)
	}

	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return res
}

func (bf BitField) MarshalCBOR(w io.Writer) error {
	ints := make([]uint64, 0, len(bf.bits))
	for i := range bf.bits {
		ints = append(ints, i)
	}

	rle, _, err := rleplus.Encode(ints) // Encode sorts internally
	if err != nil {
		return err
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajByteString, uint64(len(rle)))); err != nil {
		return err
	}
	if _, err = w.Write(rle); err != nil {
		return xerrors.Errorf("writing rle: %w", err)
	}
	return nil
}

func (bf *BitField) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if extra > 8192 {
		return fmt.Errorf("array too large")
	}

	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	rle := make([]byte, extra)
	if _, err := io.ReadFull(br, rle); err != nil {
		return err
	}

	ints, err := rleplus.Decode(rle)
	if err != nil {
		return xerrors.Errorf("could not decode rle+: %w", err)
	}
	bf.bits = make(map[uint64]struct{})
	for _, i := range ints {
		bf.bits[i] = struct{}{}
	}

	return nil
}

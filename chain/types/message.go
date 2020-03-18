package types

import (
	"bytes"
	"fmt"
	"io"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	crypto "github.com/filecoin-project/specs-actors/actors/crypto"
	block "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	cbg "github.com/whyrusleeping/cbor-gen"
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

func (t *Message) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{136}); err != nil {
		return err
	}

	// t.To (address.Address) (struct)
	if err := t.To.MarshalCBOR(w); err != nil {
		return err
	}

	// t.From (address.Address) (struct)
	if err := t.From.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Nonce (uint64) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajNegativeInt, uint64(t.CallSeqNum))); err != nil {
		return err
	}

	// t.Value (big.Int) (struct)
	if err := t.Value.MarshalCBOR(w); err != nil {
		return err
	}

	// t.GasPrice (big.Int) (struct)
	if err := t.GasPrice.MarshalCBOR(w); err != nil {
		return err
	}

	// t.GasLimit (big.Int) (struct)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajNegativeInt, uint64(t.GasLimit))); err != nil {
		return err
	}

	// t.Method (abi.MethodNum) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajNegativeInt, uint64(t.Method))); err != nil {
		return err
	}

	// t.Params ([]uint8) (slice)
	if len(t.Params) > cbg.ByteArrayMaxLen {
		return fmt.Errorf("Byte array in field t.Params was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajByteString, uint64(len(t.Params)))); err != nil {
		return err
	}
	if _, err := w.Write(t.Params); err != nil {
		return err
	}
	return nil
}

func (t *Message) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 8 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.To (address.Address) (struct)

	{

		if err := t.To.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.From (address.Address) (struct)

	{

		if err := t.From.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.Nonce (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajNegativeInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.CallSeqNum = int64(extra)

	}
	// t.Value (big.Int) (struct)

	{

		if err := t.Value.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.GasPrice (big.Int) (struct)

	{

		if err := t.GasPrice.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.GasLimit (big.Int) (struct)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajNegativeInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.GasLimit = int64(extra)

	}
	// t.Method (abi.MethodNum) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Method = abi.MethodNum(extra)

	}
	// t.Params ([]uint8) (slice)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.Params: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}
	t.Params = make([]byte, extra)
	if _, err := io.ReadFull(br, t.Params); err != nil {
		return err
	}
	return nil
}

func (m *Message) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := m.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Message) Cid() cid.Cid {
	data, err := m.Serialize()
	if err != nil {
		panic(err)
	}

	pref := cid.NewPrefixV1(cid.DagCBOR, multihash.BLAKE2B_MIN+31)
	c, err := pref.Sum(data)
	if err != nil {
		panic(err)
	}

	blk, err := block.NewBlockWithCid(data, c)
	if err != nil {
		panic(err)
	}
	return blk.Cid()

}

type SignedMessage struct {
	Message   Message
	Signature crypto.Signature
}

func (sm *SignedMessage) MarshalCBOR(w io.Writer) error {
	if sm == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{130}); err != nil {
		return err
	}

	// sm.Message (types.Message) (struct)
	if err := sm.Message.MarshalCBOR(w); err != nil {
		return err
	}

	// sm.Signature (crypto.Signature) (struct)
	if err := sm.Signature.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (sm *SignedMessage) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// sm.Message (types.Message) (struct)

	{

		if err := sm.Message.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// sm.Signature (crypto.Signature) (struct)

	{

		if err := sm.Signature.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	return nil
}

func (sm *SignedMessage) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := sm.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (sm *SignedMessage) Cid() cid.Cid {
	if sm.Signature.Type == crypto.SigTypeBLS {
		return sm.Message.Cid()
	}

	data, err := sm.Serialize()
	if err != nil {
		panic(err)
	}

	pref := cid.NewPrefixV1(cid.DagCBOR, multihash.BLAKE2B_MIN+31)
	c, err := pref.Sum(data)
	if err != nil {
		panic(err)
	}

	blk, err := block.NewBlockWithCid(data, c)
	if err != nil {
		panic(err)
	}
	return blk.Cid()
}

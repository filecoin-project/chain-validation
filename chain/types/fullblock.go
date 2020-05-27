package types

import (
	"bytes"
	"github.com/ipfs/go-cid"
)

type FullBlock struct {
	Header        *BlockHeader
	BlsMessages   []cid.Cid
	SecpkMessages []cid.Cid
}

func DecodeFullBlock(b []byte) (*FullBlock, error) {
	var bm FullBlock
	if err := bm.UnmarshalCBOR(bytes.NewReader(b)); err != nil {
		return nil, err
	}

	return &bm, nil
}

func (fb *FullBlock) Cid() cid.Cid {
	return fb.Header.Cid()
}

func (fb *FullBlock) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := fb.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

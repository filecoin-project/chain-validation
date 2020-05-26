package types

import (
	"bytes"
	"github.com/minio/blake2b-simd"
)

type Ticket struct {
	VRFProof []byte
}

func (t *Ticket) Equals(ot *Ticket) bool {
	return bytes.Equal(t.VRFProof, ot.VRFProof)
}

func (t *Ticket) Less(o *Ticket) bool {
	tDigest := blake2b.Sum256(t.VRFProof)
	oDigest := blake2b.Sum256(o.VRFProof)
	return bytes.Compare(tDigest[:], oDigest[:]) < 0
}

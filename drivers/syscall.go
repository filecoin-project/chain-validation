package drivers

import (
	"bytes"

	"github.com/filecoin-project/go-address"
	commcid "github.com/filecoin-project/go-fil-commcid"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime"
	"github.com/ipfs/go-cid"
	"github.com/minio/blake2b-simd"
)

type ChainValidationSyscalls struct {
}

func NewChainValidationSyscalls() *ChainValidationSyscalls {
	return &ChainValidationSyscalls{}
}

func (c ChainValidationSyscalls) VerifySignature(signature crypto.Signature, signer address.Address, plaintext []byte) error {
	return nil
}

func (c ChainValidationSyscalls) HashBlake2b(data []byte) [32]byte {
	return blake2b.Sum256(data)
}

func (c ChainValidationSyscalls) ComputeUnsealedSectorCID(proof abi.RegisteredProof, pieces []abi.PieceInfo) (cid.Cid, error) {
	// Fake CID computation by hashing the piece info (rather than the real computation over piece commitments).
	buf := bytes.Buffer{}
	for _, p := range pieces {
		err := p.MarshalCBOR(&buf)
		if err != nil {
			panic(err)
		}
	}
	token := blake2b.Sum256(buf.Bytes())
	return commcid.DataCommitmentV1ToCID(token[:]), nil
}

func (c ChainValidationSyscalls) VerifySeal(info abi.SealVerifyInfo) error {
	return nil
}

func (c ChainValidationSyscalls) VerifyPoSt(info abi.PoStVerifyInfo) error {
	return nil
}

func (c ChainValidationSyscalls) VerifyConsensusFault(h1, h2, extra []byte, earliest abi.ChainEpoch) (*runtime.ConsensusFault, error) {
	panic("implement me")
}

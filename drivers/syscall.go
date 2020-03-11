package drivers

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime"
	"github.com/ipfs/go-cid"
	"github.com/minio/blake2b-simd"

	"github.com/filecoin-project/chain-validation/suites/utils"
)

type ChainValidationSyscalls struct {
	unsealedSectorCIDMaker func() cid.Cid
}

func NewChainValidationSyscalls() *ChainValidationSyscalls {
	return &ChainValidationSyscalls{unsealedSectorCIDMaker: utils.NewProofCidForTestGetter()}
}

func (c ChainValidationSyscalls) VerifySignature(signature crypto.Signature, signer address.Address, plaintext []byte) error {
	return nil
}

func (c ChainValidationSyscalls) HashBlake2b(data []byte) [32]byte {
	hasher, err := blake2b.New(&blake2b.Config{Size: 32})
	if err != nil {
		panic(err)
	}
	var something [32]byte
	digest := hasher.Sum(data)
	n := copy(something[:], digest)
	if n != 32 {
		panic("borked")
	}
	return something
}

func (c ChainValidationSyscalls) ComputeUnsealedSectorCID(proof abi.RegisteredProof, pieces []abi.PieceInfo) (cid.Cid, error) {
	return c.unsealedSectorCIDMaker(), nil
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

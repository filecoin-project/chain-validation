package utils

import (
	"math/rand"

	fancypantscidmaker "github.com/filecoin-project/go-fil-commcid"
	"github.com/ipfs/go-cid"
)

// NewProofCidForTestGetter returns a closure that returns a Cid unique to that invocation and has the CommD/R prefix
// The Cid is unique wrt the closure returned, not globally. You can use this function
// in tests.
func NewProofCidForTestGetter() func() cid.Cid {
	rand.Seed(1)
	return func() cid.Cid {
		token := make([]byte, 32)
		_, err := rand.Read(token)
		if err != nil {
			panic(err)
		}
		proofCid, err := fancypantscidmaker.CommitmentToCID(token, fancypantscidmaker.FC_SEALED_V1)
		if err != nil {
			panic(err)
		}
		return proofCid
	}
}

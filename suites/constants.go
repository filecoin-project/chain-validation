package suites

import (
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain"
)

const (
	totalFilecoin     = 2_000_000_000
	filecoinPrecision = 1_000_000_000_000_000_000
)

func init() {
	EmptyReturnValue = chain.MustSerialize(adt_spec.EmptyValue{})

}

var (
	TotalNetworkBalance = big_spec.Mul(big_spec.NewInt(totalFilecoin), big_spec.NewInt(filecoinPrecision))
	EmptyReturnValue    []byte
)

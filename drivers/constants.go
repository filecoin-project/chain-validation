package drivers

import (
	big_spec "github.com/filecoin-project/go-state-types/big"
)

const (
	totalFilecoin     = 2_000_000_000
	filecoinPrecision = 1_000_000_000_000_000_000
)

var (
	TotalNetworkBalance = big_spec.Mul(big_spec.NewInt(totalFilecoin), big_spec.NewInt(filecoinPrecision))
	EmptyReturnValue    = []byte{}
)

package main

import (
	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/pkg/state/actors/paych"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

func main() {

	// General Types
	if err := gen.WriteTupleEncodersToFile("../state/types/cbor_gen.go", "types",
		types.SignedVoucher{},
		types.Merge{},
		types.ModVerifyParams{},
		types.Signature{},
	); err != nil {
		panic(err)
	}

	// Payment Channel Actor
	if err := gen.WriteTupleEncodersToFile("../state/actors/paych/cbor_gen.go", "paych",
		paych.PaymentInfo{},
		paych.PaymentChannelActorState{},
		paych.LaneState{},
	); err != nil {
		panic(err)
	}
}

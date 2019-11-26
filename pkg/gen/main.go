package main

import (
	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/pkg/state/actors"
)

func main() {
	err := gen.WriteTupleEncodersToFile("./state/actors/cbor_gen.go", "types",
		actors.PaymentChannelActorStateType{},
		actors.LaneState{},
	)
	if err != nil {
		panic(err)
	}
}

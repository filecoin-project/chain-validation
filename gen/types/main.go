package main

import (
	"fmt"
	"os"

	types "github.com/filecoin-project/chain-validation/chain/types"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	err := gen.WriteTupleEncodersToFile("./chain/types/cbor_gen.go", "types",
		types.ExpTipSet{},

		types.BlockHeader{},
		types.FullBlock{},

		types.MessageReceipt{},
		types.Message{},

		types.SignedMessage{},

		types.Ticket{},
		types.ElectionProof{},
		types.BeaconEntry{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

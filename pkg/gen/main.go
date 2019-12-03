package main

import (
	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/pkg/state/actors/initialize"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/paych"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgminr"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgpwr"
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

	// Init Actor
	if err := gen.WriteTupleEncodersToFile("../state/actors/initialize/cbor_gen.go", "initialize",
		initialize.ExecParams{},
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

	// Storage Power Actor
	if err := gen.WriteTupleEncodersToFile("../state/actors/strgpwr/cbor_gen.go", "strgpwr",
		strgpwr.CreateStorageMinerParams{},
		strgpwr.UpdateStorageParams{},
	); err != nil {
		panic(err)
	}

	// Storage Miner Actor
	if err := gen.WriteTupleEncodersToFile("../state/actors/strgminr/cbor_gen.go", "strgminr",
		strgminr.StorageMinerActorState{},
		strgminr.MinerInfo{},
		strgminr.PreCommittedSector{},
		strgminr.SectorPreCommitInfo{},
		strgminr.UpdatePeerIDParams{},
	); err != nil {
		panic(err)
	}
}

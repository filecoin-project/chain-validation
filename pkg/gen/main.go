package main

import (
	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/pkg/state/actors/initialize"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/multsig"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/paych"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgminr"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgmrkt"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgpwr"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

func main() {

	// General Types
	if err := gen.WriteTupleEncodersToFile("./pkg/state/types/cbor_gen.go", "types",
		types.SignedVoucher{},
		types.Merge{},
		types.ModVerifyParams{},
	); err != nil {
		panic(err)
	}

	// Init Actor
	if err := gen.WriteTupleEncodersToFile("./pkg/state/actors/initialize/cbor_gen.go", "initialize",
		initialize.ExecParams{},
	); err != nil {
		panic(err)
	}

	// Payment Channel Actor
	if err := gen.WriteTupleEncodersToFile("./pkg/state/actors/paych/cbor_gen.go", "paych",
		paych.PaymentInfo{},
		paych.PaymentChannelActorState{},
		paych.LaneState{},
		paych.PaymentChannelConstructorParams{},
		paych.PaymentChannelUpdateParams{},
		paych.PaymentVerifyParams{},
	); err != nil {
		panic(err)
	}

	// Storage Power Actor
	if err := gen.WriteTupleEncodersToFile("./pkg/state/actors/strgpwr/cbor_gen.go", "strgpwr",
		strgpwr.CreateStorageMinerParams{},
		strgpwr.UpdateStorageParams{},
		strgpwr.StoragePowerState{},
		strgpwr.PledgeCollateralParams{},
		strgpwr.PowerLookupParams{},
	); err != nil {
		panic(err)
	}

	// Storage Market Actor
	if err := gen.WriteTupleEncodersToFile("./pkg/state/actors/strgmrkt/cbor_gen.go", "strgmrkt",
		strgmrkt.OnChainDeal{},
		strgmrkt.StorageDeal{},
		strgmrkt.StorageMarketState{},
		strgmrkt.StorageDealProposal{},
		strgmrkt.WithdrawBalanceParams{},
		strgmrkt.PublishStorageDealsParams{},
		strgmrkt.StorageParticipantBalance{},
		strgmrkt.ActivateStorageDealsParams{},
		strgmrkt.PublishStorageDealResponse{},
		strgmrkt.ComputeDataCommitmentParams{},
		strgmrkt.ProcessStorageDealsPaymentParams{},
	); err != nil {
		panic(err)
	}

	// Storage Miner Actor
	if err := gen.WriteTupleEncodersToFile("./pkg/state/actors/strgminr/cbor_gen.go", "strgminr",
		strgminr.StorageMinerActorState{},
		strgminr.MinerInfo{},
		strgminr.PreCommittedSector{},
		strgminr.SectorPreCommitInfo{},
		strgminr.UpdatePeerIDParams{},
	); err != nil {
		panic(err)
	}

	// Multi Signature Actor
	if err := gen.WriteTupleEncodersToFile("./pkg/state/actors/multsig/cbor_gen.go", "multsig",
		multsig.MultiSigActorState{},
		multsig.MTransaction{},
		multsig.MultiSigTxID{},
		multsig.MultiSigConstructorParams{},
		multsig.MultiSigChangeReqParams{},
		multsig.MultiSigAddSignerParam{},
		multsig.MultiSigProposeParams{},
		multsig.MultiSigRemoveSignerParam{},
		multsig.MultiSigSwapSignerParams{},
	); err != nil {
		panic(err)
	}
}

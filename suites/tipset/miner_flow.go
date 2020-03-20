package tipset

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	market_spec "github.com/filecoin-project/specs-actors/actors/builtin/market"
	miner_spec "github.com/filecoin-project/specs-actors/actors/builtin/miner"
	power_spec "github.com/filecoin-project/specs-actors/actors/builtin/power"
	crypto_spec "github.com/filecoin-project/specs-actors/actors/crypto"
	exitcode_spec "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func CreateMinerWithProvenCommittedSector(td *drivers.TestDriver, minerOwner, minerWorker, minerIDAddr address.Address, sectorType abi_spec.RegisteredProof, collateral abi_spec.TokenAmount, dealStart, dealEnd abi_spec.ChainEpoch) *types.PreSeal {
	td.AssertBalanceCallback(minerOwner, func(actorBal abi_spec.TokenAmount) bool {
		return actorBal.GreaterThanEqual(collateral)
	})
	td.AssertBalanceCallback(minerWorker, func(actorBal abi_spec.TokenAmount) bool {
		return actorBal.GreaterThanEqual(collateral)
	})

	// since we do not know the callseq of this actor when this method is called
	ownerActor, err := td.State().Actor(minerOwner)
	require.NoError(td.T, err)
	ownerCallSeq := ownerActor.CallSeqNum()
	incOwnerCallSeq := func() int64 {
		// fell down the ugly stack and hit every frame on the way down
		defer func() { ownerCallSeq += 1 }()
		return ownerCallSeq
	}

	workerActor, err := td.State().Actor(minerWorker)
	require.NoError(td.T, err)
	workerCallSeq := workerActor.CallSeqNum()
	incWorkerCallSeq := func() int64 {
		defer func() { workerCallSeq += 1 }()
		return workerCallSeq
	}

	sectorBuilder := drivers.NewMockSectorBuilder()
	bb := drivers.NewTipSetMessageBuilder(td)

	sectorSize, err := sectorType.SectorSize()
	require.NoError(td.T, err)

	createMinerRet := td.ComputeInitActorExecReturn(minerOwner, 0, 0, minerIDAddr)

	// Create a miner and add finds to the storage market actor for the miner and a client
	bb.WithTicketCount(1).
		// Step 1: Register the miner with the power actor
		WithBLSMessageAndReceipt(
			td.MessageProducer.PowerCreateMiner(
				builtin_spec.StoragePowerActorAddr, minerOwner,
				power_spec.CreateMinerParams{Owner: minerOwner, Worker: minerWorker, SectorSize: sectorSize, Peer: utils.RequireRandomPeerID(td.T)},
				chain.Nonce(incOwnerCallSeq()), chain.Value(collateral),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: chain.MustSerialize(&createMinerRet), GasUsed: big_spec.Zero()},
		).
		// Step 2.A: Add market funds for client
		WithBLSMessageAndReceipt(
			td.MessageProducer.MarketAddBalance(builtin_spec.StorageMarketActorAddr, minerWorker,
				minerIDAddr,
				chain.Nonce(incWorkerCallSeq()), chain.Value(collateral),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		// Step 2.B: Add market funds for miner
		WithBLSMessageAndReceipt(
			td.MessageProducer.MarketAddBalance(builtin_spec.StorageMarketActorAddr, minerOwner,
				minerWorker,
				chain.Nonce(incOwnerCallSeq()), chain.Value(collateral),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	// The miner preseals a sector
	sectorInfo := sectorBuilder.NewPreSealedSector(minerIDAddr, minerWorker, sectorType, sectorSize, dealStart, dealEnd)

	dealIDs := []abi_spec.DealID{abi_spec.DealID(0)}
	pubRet := chain.MustSerialize(&market_spec.PublishStorageDealsReturn{IDs: dealIDs})

	td.ExeCtx.Epoch++
	// Miner publishes deal to the storage market and precommits its sector
	bb.WithTicketCount(1).
		// Step 3: Publish presealed deals
		WithBLSMessageAndReceipt(
			td.MessageProducer.MarketPublishStorageDeals(builtin_spec.StorageMarketActorAddr, minerWorker,
				market_spec.PublishStorageDealsParams{
					Deals: []market_spec.ClientDealProposal{
						{Proposal: sectorInfo.Deal, ClientSignature: crypto_spec.Signature{Type: crypto_spec.SigTypeBLS, Data: []byte("doesnt matter")}},
					},
				},
				chain.Nonce(incWorkerCallSeq()), chain.Value(big_spec.Zero()),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: pubRet, GasUsed: big_spec.Zero()},
		).
		// Step 4: Pre Committing Sectors
		WithBLSMessageAndReceipt(
			td.MessageProducer.MinerPreCommitSector(minerIDAddr, minerWorker,
				miner_spec.SectorPreCommitInfo{
					RegisteredProof: sectorInfo.ProofType,
					SectorNumber:    sectorInfo.SectorID,
					SealedCID:       sectorInfo.CommR,
					SealRandEpoch:   0,
					DealIDs:         dealIDs,
					Expiration:      sectorInfo.Deal.EndEpoch,
				},
				chain.Nonce(incWorkerCallSeq()), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	// Miner prove commits its sector
	td.ExeCtx.Epoch = dealStart
	bb.WithTicketCount(1).
		WithBLSMessageAndReceipt(
			// Step 5: Prove the committed sector
			td.MessageProducer.MinerProveCommitSector(minerIDAddr, minerWorker,
				miner_spec.ProveCommitSectorParams{SectorNumber: sectorInfo.SectorID, Proof: nil},
				chain.Nonce(incWorkerCallSeq()), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	return sectorInfo
}

func TestMinerMissPoStChallengeWindow(t *testing.T, factory state.Factories) {
	td := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState).
		Build(t)

	bb := drivers.NewTipSetMessageBuilder(td)

	// The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it.
	minerOwner, _ := td.NewAccountActor(address.SECP256K1, abi_spec.NewTokenAmount(1_000_000_000))
	// minerWorker address will be responsible for doing all of the work, submitting proofs, committing new sectors,
	// and all other day to day activities.
	minerWorker, minerWorkerID := td.NewAccountActor(address.BLS, abi_spec.NewTokenAmount(1_000_000_000))
	// The address of the miner actor
	minerActorID := utils.NewIDAddr(t, utils.IdFromAddress(minerWorkerID)+1)

	collateral := abi_spec.NewTokenAmount(1_000_000)
	sectorProofType := abi_spec.RegisteredProof_StackedDRG32GiBSeal
	dealStart := abi_spec.ChainEpoch(15)
	dealEnd := abi_spec.ChainEpoch(1000)

	CreateMinerWithProvenCommittedSector(td, minerOwner, minerWorker, minerActorID, sectorProofType, collateral, dealStart, dealEnd)
	require.Equal(t, dealStart, td.ExeCtx.Epoch)
	// since we do not know the callseq of this actor when this method is called
	ownerActor, err := td.State().Actor(minerOwner)
	require.NoError(td.T, err)
	ownerCallSeq := ownerActor.CallSeqNum()
	incOwnerCallSeq := func() int64 {
		defer func() { ownerCallSeq += 1 }()
		return ownerCallSeq
	}

	// Epoch advances to the end of the proving window. Send a sing message to trigger the cron actor send
	td.ExeCtx.Epoch += power_spec.WindowedPostChallengeDuration + miner_spec.ProvingPeriod
	bb.WithTicketCount(1).
		// Step 6: send a single message that causes the cron actor to trigger
		WithBLSMessageAndReceipt(
			td.MessageProducer.Transfer(minerOwner, minerOwner,
				chain.Nonce(incOwnerCallSeq()), chain.Value(big_spec.Zero()),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	// after the application of the message and moving the epoch past the proving period the miner has had a fault
	var minerSt miner_spec.State
	td.GetActorState(minerActorID, &minerSt)
	require.True(t, minerSt.PoStState.HasFailedPost())
	require.Equal(t, int64(1), minerSt.PoStState.NumConsecutiveFailures)

	// NB: the power actors TotalNetworkPower filed will not change since ConsensusMinerMinPower is larger than
	// what would be a reasonable amount of sectors to seal in a test.
}

func TestMinerSubmitFallbackPoSt(t *testing.T, factory state.Factories) {
	td := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState).
		Build(t)

	bb := drivers.NewTipSetMessageBuilder(td)

	// The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it.
	minerOwner, _ := td.NewAccountActor(address.SECP256K1, abi_spec.NewTokenAmount(1_000_000_000))
	// minerWorker address will be responsible for doing all of the work, submitting proofs, committing new sectors,
	// and all other day to day activities.
	minerWorker, minerWorkerID := td.NewAccountActor(address.BLS, abi_spec.NewTokenAmount(1_000_000_000))
	// The address of the miner actor
	minerActorID := utils.NewIDAddr(t, utils.IdFromAddress(minerWorkerID)+1)

	collateral := abi_spec.NewTokenAmount(1_000_000)
	sectorProofType := abi_spec.RegisteredProof_StackedDRG32GiBSeal
	dealStart := abi_spec.ChainEpoch(15)
	dealEnd := abi_spec.ChainEpoch(1000)

	CreateMinerWithProvenCommittedSector(td, minerOwner, minerWorker, minerActorID, sectorProofType, collateral, dealStart, dealEnd)
	require.Equal(t, dealStart, td.ExeCtx.Epoch)
	// since we do not know the callseq of this actor when this method is called
	workerActor, err := td.State().Actor(minerWorker)
	require.NoError(td.T, err)
	workerCallSeq := workerActor.CallSeqNum()
	incWorkerCallSeq := func() int64 {
		defer func() { workerCallSeq += 1 }()
		return workerCallSeq
	}

	ownerActor, err := td.State().Actor(minerOwner)
	require.NoError(td.T, err)
	ownerCallSeq := ownerActor.CallSeqNum()
	incOwnerCallSeq := func() int64 {
		defer func() { ownerCallSeq += 1 }()
		return ownerCallSeq
	}

	candidates := []abi_spec.PoStCandidate{{
		RegisteredProof: abi_spec.RegisteredProof_StackedDRG32GiBPoSt,
		ChallengeIndex:  0,
	}}
	proofs := []abi_spec.PoStProof{{
		RegisteredProof: abi_spec.RegisteredProof_StackedDRG32GiBPoSt,
		ProofBytes:      []byte("doesn't matter"),
	}}

	// move the epoch forward to be withing the proving period window.
	td.ExeCtx.Epoch += power_spec.WindowedPostChallengeDuration + miner_spec.ProvingPeriod/2
	bb.WithTicketCount(1).
		WithBLSMessageAndReceipt(
			td.MessageProducer.MinerSubmitWindowedPoSt(minerActorID, minerWorker, abi_spec.OnChainPoStVerifyInfo{Candidates: candidates, Proofs: proofs},
				chain.Nonce(incWorkerCallSeq()), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	// move the epoch outside of the proving window and send a message to trigger the cron actor
	td.ExeCtx.Epoch += miner_spec.ProvingPeriod/2 + 1
	bb.WithTicketCount(1).
		// Step 6: send a single message that causes the cron actor to trigger
		WithBLSMessageAndReceipt(
			td.MessageProducer.Transfer(minerOwner, minerOwner,
				chain.Nonce(incOwnerCallSeq()), chain.Value(big_spec.Zero()),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	// after the application of the message and moving the epoch past the proving period the miner has not had a fault
	var minerSt miner_spec.State
	td.GetActorState(minerActorID, &minerSt)
	assert.False(t, minerSt.PoStState.HasFailedPost())

	var powerSt power_spec.State
	td.GetActorState(builtin_spec.StoragePowerActorAddr, &powerSt)
	assert.Equal(t, drivers.InitialTotalNetworkPower+int64(32<<30), powerSt.TotalNetworkPower.Int64())
}

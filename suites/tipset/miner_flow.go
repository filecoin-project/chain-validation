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
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestMinerCreateProveCommitAndMissPoStChallengeWindow(t *testing.T, factory state.Factories) {
	td := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState).Build(t)

	sectorBuilder := drivers.NewMockSectorBuilder()
	blkBuilder := drivers.NewTipSetMessageBuilder(td)

	//
	// Define Parameters
	//

	// The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it.
	minerOwner, _ := td.NewAccountActor(address.SECP256K1, abi_spec.NewTokenAmount(1_000_000_000))
	// minerWorker address will be responsible for doing all of the work, submitting proofs, committing new sectors,
	// and all other day to day activities.
	minerWorker, minerWorkerID := td.NewAccountActor(address.BLS, abi_spec.NewTokenAmount(1_000_000_000))
	// The address of the miner actor
	minerActorID := utils.NewIDAddr(t, utils.IdFromAddress(minerWorkerID)+1)
	createMinerRet := td.ComputeInitActorExecReturn(minerOwner, 0, 1, minerActorID)
	// collaterall the miner will pledge
	collateral := abi_spec.NewTokenAmount(1_000_000)

	sectorProofType := abi_spec.RegisteredProof_StackedDRG32GiBSeal
	sectorSize, err := sectorProofType.SectorSize()
	require.NoError(t, err)

	//
	// Apply messages and assert result
	//

	// Create a miner and add finds to the storage market actor for the miner and a client
	blkBuilder.WithTicketCount(1).
		// Step 1: Register the miner with the power actor
		WithBLSMessageAndReceipt(
			td.MessageProducer.PowerCreateMiner(
				builtin_spec.StoragePowerActorAddr, minerOwner,
				power_spec.CreateMinerParams{Owner: minerOwner, Worker: minerWorker, SectorSize: sectorSize, Peer: utils.RequireRandomPeerID(t)},
				chain.Nonce(0), chain.Value(collateral),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: chain.MustSerialize(&createMinerRet), GasUsed: big_spec.Zero()},
		).
		// Step 2.A: Add market funds for client
		WithBLSMessageAndReceipt(
			td.MessageProducer.MarketAddBalance(builtin_spec.StorageMarketActorAddr, minerWorker,
				minerActorID,
				chain.Nonce(0), chain.Value(collateral),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		// Step 2.B: Add market funds for miner
		WithBLSMessageAndReceipt(
			td.MessageProducer.MarketAddBalance(builtin_spec.StorageMarketActorAddr, minerOwner,
				minerWorker,
				chain.Nonce(1), chain.Value(collateral),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	dealStart := abi_spec.ChainEpoch(15)
	dealEnd := abi_spec.ChainEpoch(1000)
	// The miner preseals a sector
	sectorInfo := sectorBuilder.NewPreSealedSector(minerActorID, minerWorker, sectorProofType, sectorSize, dealStart, dealEnd)

	dealIDs := []abi_spec.DealID{abi_spec.DealID(0)}
	pubRet := chain.MustSerialize(&market_spec.PublishStorageDealsReturn{IDs: dealIDs})

	td.ExeCtx.Epoch++
	// Miner publishes deal to the storage market and precommits its sector
	blkBuilder.WithTicketCount(1).
		// Step 3: Publish presealed deals
		WithBLSMessageAndReceipt(
			td.MessageProducer.MarketPublishStorageDeals(builtin_spec.StorageMarketActorAddr, minerWorker,
				market_spec.PublishStorageDealsParams{
					Deals: []market_spec.ClientDealProposal{
						{Proposal: sectorInfo.Deal, ClientSignature: crypto_spec.Signature{Type: crypto_spec.SigTypeBLS, Data: []byte("doesnt matter")}},
					},
				},
				chain.Nonce(1), chain.Value(big_spec.Zero()),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: pubRet, GasUsed: big_spec.Zero()},
		).
		// Step 4: Pre Committing Sectors
		WithBLSMessageAndReceipt(
			td.MessageProducer.MinerPreCommitSector(minerActorID, minerWorker,
				miner_spec.SectorPreCommitInfo{
					RegisteredProof: sectorInfo.ProofType,
					SectorNumber:    sectorInfo.SectorID,
					SealedCID:       sectorInfo.CommR,
					SealRandEpoch:   0,
					DealIDs:         dealIDs,
					Expiration:      sectorInfo.Deal.EndEpoch,
				},
				chain.Nonce(2), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	// Miner prove commits its sector
	td.ExeCtx.Epoch = dealStart
	blkBuilder.WithTicketCount(1).
		WithBLSMessageAndReceipt(
			// Step 5: Prove the committed sector
			td.MessageProducer.MinerProveCommitSector(minerActorID, minerWorker,
				miner_spec.ProveCommitSectorParams{SectorNumber: sectorInfo.SectorID, Proof: nil},
				chain.Value(big_spec.Zero()), chain.Nonce(3)),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	// the miner has not failed to PoSt
	var minerSt miner_spec.State
	td.GetActorState(minerActorID, &minerSt)
	require.False(t, minerSt.PoStState.HasFailedPost())

	// Epoch advances to the end of the proving window. Send a sing message to trigger the cron actor send
	td.ExeCtx.Epoch += power_spec.WindowedPostChallengeDuration + miner_spec.ProvingPeriod
	blkBuilder.WithTicketCount(1).
		// Step 6: send a single message that causes the cron actor to trigger
		WithBLSMessageAndReceipt(
			td.MessageProducer.Transfer(minerOwner, minerOwner,
				chain.Nonce(2), chain.Value(big_spec.Zero()),
			),
			types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		).
		ApplyAndValidate()

	// after the application of the message and moving the epoch past the proving period the miner has had a fault
	td.GetActorState(minerActorID, &minerSt)
	require.True(t, minerSt.PoStState.HasFailedPost())
	require.Equal(t, int64(1), minerSt.PoStState.NumConsecutiveFailures)

	// NB: the power actors TotalNetworkPower filed will not change since ConsensusMinerMinPower is larger than
	// what would be a reasonable amount of sectors to seal in a test.
}

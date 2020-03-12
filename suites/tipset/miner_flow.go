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
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	sectorBuilder := drivers.NewMockSectorBuilder()
	var dealStart = abi_spec.ChainEpoch(15)
	var dealEnd = abi_spec.ChainEpoch(1000)

	t.Run("create a miner, pre commit then commit a sector, then miss the proving window", func(t *testing.T) {
		td := builder.Build(t)

		// The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it.
		minerOwner, _ := td.NewAccountActor(address.SECP256K1, abi_spec.NewTokenAmount(1_000_000_000))
		// minerWorker address will be responsible for doing all of the work, submitting proofs, committing new sectors,
		// and all other day to day activities.
		minerWorker, minerWorkerID := td.NewAccountActor(address.BLS, abi_spec.NewTokenAmount(1_000_000_000))
		// the next actor to be created will have an ID = previous actor + 1
		minerActorID := utils.NewIDAddr(t, utils.IdFromAddress(minerWorkerID)+1)
		createMinerRet := td.ComputeInitActorExecReturn(builtin_spec.StoragePowerActorAddr, 0, minerActorID)
		// miners need to pledge collateral
		// TODO this is probably the wrong value and should be calculated instead based on sector size
		collateral := abi_spec.NewTokenAmount(1_000_000)

		sectorProofType := abi_spec.RegisteredProof_StackedDRG32GiBSeal
		sectorSize, err := sectorProofType.SectorSize()
		require.NoError(t, err)

		// Step 1: Register teh miner with the power actor
		createMinerMsg := td.MessageProducer.PowerCreateMiner(builtin_spec.StoragePowerActorAddr, minerOwner, power_spec.CreateMinerParams{
			Owner:      minerOwner,
			Worker:     minerWorker,
			SectorSize: sectorSize,
			Peer:       utils.RequireRandomPeerID(t),
		}, chain.Value(collateral), chain.Nonce(0))
		createMinerRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: chain.MustSerialize(&createMinerRet), GasUsed: big_spec.Zero()}

		minerIdAddr := createMinerRet.IDAddress

		// Step 2: Add market funds from miner and client
		minerAddBalMsg := td.MessageProducer.MarketAddBalance(builtin_spec.StorageMarketActorAddr, minerWorker, minerIdAddr, chain.Nonce(0), chain.Value(collateral))
		minerAddBalRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		clientAddBalMsg := td.MessageProducer.MarketAddBalance(builtin_spec.StorageMarketActorAddr, minerOwner, minerWorker, chain.Nonce(1), chain.Value(collateral))
		clientAddBalRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		blkMsgs := chain.NewTipSetMessageBuilder().
			WithTicketCount(1).
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(createMinerMsg).
			WithBLSMessage(minerAddBalMsg).
			WithBLSMessage(clientAddBalMsg).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 3)

		td.AssertReceipt(receipts[0], createMinerRct)
		td.AssertReceipt(receipts[1], minerAddBalRct)
		td.AssertReceipt(receipts[2], clientAddBalRct)

		sectorInfo := sectorBuilder.NewPreSealedSector(minerIdAddr, minerWorker, sectorProofType, sectorSize, dealStart, dealEnd)

		// Step 3: Publish presealed deals
		dealID := abi_spec.DealID(0)
		dealIDs := []abi_spec.DealID{dealID}
		pubRet := chain.MustSerialize(&market_spec.PublishStorageDealsReturn{IDs: dealIDs})

		pubDealMsg := td.MessageProducer.MarketPublishStorageDeals(builtin_spec.StorageMarketActorAddr, minerWorker, market_spec.PublishStorageDealsParams{Deals: []market_spec.ClientDealProposal{{
			Proposal:        sectorInfo.Deal,
			ClientSignature: crypto_spec.Signature{Type: crypto_spec.SigTypeBLS, Data: []byte("doesnt matter")},
		}}}, chain.Nonce(1), chain.Value(big_spec.Zero()))
		pubDealRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: pubRet, GasUsed: big_spec.Zero()}

		// Step 4: Pre Committing Sectors
		preCommitMsg := td.MessageProducer.MinerPreCommitSector(minerIdAddr, minerWorker, miner_spec.SectorPreCommitInfo{
			RegisteredProof: sectorInfo.ProofType,
			SectorNumber:    sectorInfo.SectorID,
			SealedCID:       sectorInfo.CommR,
			SealRandEpoch:   0,
			DealIDs:         dealIDs,
			Expiration:      sectorInfo.Deal.EndEpoch,
		}, chain.Nonce(2), chain.Value(big_spec.Zero()))
		preCommitRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		blkMsgs = chain.NewTipSetMessageBuilder().
			WithTicketCount(1).
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(pubDealMsg).
			WithBLSMessage(preCommitMsg).
			Build()

		td.ExeCtx.Epoch++
		receipts, err = td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 2)

		td.AssertReceipt(receipts[0], pubDealRct)
		td.AssertReceipt(receipts[1], preCommitRct)

		proveCommitMsg := td.MessageProducer.MinerProveCommitSector(minerIdAddr, minerWorker, miner_spec.ProveCommitSectorParams{
			SectorNumber: sectorInfo.SectorID,
			Proof:        nil,
		}, chain.Value(big_spec.Zero()), chain.Nonce(3))
		proveCommitRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		blkMsgs = chain.NewTipSetMessageBuilder().
			WithTicketCount(1).
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(proveCommitMsg).
			Build()

		td.ExeCtx.Epoch = dealStart

		receipts, err = td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)
		td.AssertReceipt(receipts[0], proveCommitRct)

		// the miner has not yet a fault
		var minerSt miner_spec.State
		td.GetActorState(minerIdAddr, &minerSt)
		require.Equal(t, int64(0), minerSt.PoStState.NumConsecutiveFailures)

		// a cron event type: CronEventWindowedPoStExpiration now exists in the power actors state
		// advance the epoch s.t. the miner misses the proving window
		td.ExeCtx.Epoch += power_spec.WindowedPostChallengeDuration + miner_spec.ProvingPeriod

		transferMsg := td.MessageProducer.Transfer(minerOwner, minerOwner, chain.Nonce(2), chain.Value(big_spec.Zero()))
		transferRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		blkMsgs = chain.NewTipSetMessageBuilder().
			WithTicketCount(1).
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(transferMsg).
			Build()

		receipts, err = td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)
		td.AssertReceipt(receipts[0], transferRct)

		// after the application of the message and moving the epoch past the proving period the miner has had a fault
		td.GetActorState(minerIdAddr, &minerSt)
		require.True(t, minerSt.PoStState.HasFailedPost())
		require.Equal(t, int64(1), minerSt.PoStState.NumConsecutiveFailures)

		// NB: the power actors TotalNetworkPower filed will not change since ConsensusMinerMinPower is larger than
		// what would be a reasonable amount of sectors to seal in a test.
	})
}

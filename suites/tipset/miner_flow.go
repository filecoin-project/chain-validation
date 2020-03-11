package tipset

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	sectorbuilder "github.com/filecoin-project/go-sectorbuilder"
	"github.com/filecoin-project/go-sectorbuilder/fs"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	market_spec "github.com/filecoin-project/specs-actors/actors/builtin/market"
	miner_spec "github.com/filecoin-project/specs-actors/actors/builtin/miner"
	power_spec "github.com/filecoin-project/specs-actors/actors/builtin/power"
	crypto_spec "github.com/filecoin-project/specs-actors/actors/crypto"
	exitcode_spec "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

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

		sectorSize, err := abi_spec.RegisteredProof_StackedDRG2KiBSeal.SectorSize()
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

		m1temp, err := ioutil.TempDir("", "preseal")
		require.NoError(t, err)
		preseals, err := PreSealSectors(minerIdAddr, abi_spec.RegisteredProof_StackedDRG2KiBSeal, 0, 1, m1temp, []byte("some randomness"))
		require.NoError(t, err)
		err = createDeals(preseals, minerWorker, minerIdAddr, sectorSize, dealStart, dealEnd)
		require.NoError(t, err)

		// Step 3: Publish presealed deals
		dealID := abi_spec.DealID(0)
		dealIDs := []abi_spec.DealID{dealID}
		pubRet := chain.MustSerialize(&market_spec.PublishStorageDealsReturn{IDs: dealIDs})

		pubDealMsg := td.MessageProducer.MarketPublishStorageDeals(builtin_spec.StorageMarketActorAddr, minerWorker, market_spec.PublishStorageDealsParams{Deals: []market_spec.ClientDealProposal{{
			Proposal:        preseals[0].Deal,
			ClientSignature: crypto_spec.Signature{Type: crypto_spec.SigTypeBLS, Data: []byte("doesnt matter")},
		}}}, chain.Nonce(1), chain.Value(big_spec.Zero()))
		pubDealRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: pubRet, GasUsed: big_spec.Zero()}

		// Step 4: Pre Committing Sectors
		preCommitMsg := td.MessageProducer.MinerPreCommitSector(minerIdAddr, minerWorker, miner_spec.SectorPreCommitInfo{
			RegisteredProof: preseals[0].ProofType,
			SectorNumber:    preseals[0].SectorID,
			SealedCID:       preseals[0].CommR,
			SealRandEpoch:   0,
			DealIDs:         dealIDs,
			Expiration:      preseals[0].Deal.EndEpoch,
		}, chain.Nonce(2), chain.Value(big_spec.Zero()))
		preCommitRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		blkMsgs = chain.NewTipSetMessageBuilder().
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
			SectorNumber: preseals[0].SectorID,
			Proof:        nil,
		}, chain.Value(big_spec.Zero()), chain.Nonce(3))
		proveCommitRct := types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		blkMsgs = chain.NewTipSetMessageBuilder().
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
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(transferMsg).
			Build()

		receipts, err = td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)
		td.AssertReceipt(receipts[0], transferRct)

		// after the application of the message and moving the epoch past the proving period the miner has had a fault
		td.GetActorState(minerIdAddr, &minerSt)
		require.Equal(t, int64(1), minerSt.PoStState.NumConsecutiveFailures)

		// NB: the power actors TotalNetworkPower filed will not change since ConsensusMinerMinPower is larger than
		// what would be a reasonable amount of sectors to seal in a test.
	})
}

type PreSeal struct {
	CommR     cid.Cid
	CommD     cid.Cid
	SectorID  abi_spec.SectorNumber
	Deal      market_spec.DealProposal
	ProofType abi_spec.RegisteredProof
}

func PreSealSectors(maddr address.Address, pt abi_spec.RegisteredProof, offset abi_spec.SectorNumber, sectors int, sbroot string, preimage []byte) ([]*PreSeal, error) {
	ppt, err := pt.RegisteredPoStProof()
	if err != nil {
		return nil, err
	}

	spt, err := pt.RegisteredSealProof()
	if err != nil {
		return nil, err
	}

	cfg := &sectorbuilder.Config{
		Miner:         maddr,
		SealProofType: spt,
		PoStProofType: ppt,
	}

	if err := os.MkdirAll(sbroot, 0775); err != nil {
		return nil, err
	}

	next := offset

	sbfs := &fs.Basic{
		Miner: maddr,
		Root:  sbroot,
	}

	sb, err := sectorbuilder.New(sbfs, cfg)
	if err != nil {
		return nil, err
	}

	ssize, err := pt.SectorSize()
	if err != nil {
		return nil, err
	}

	var sealedSectors []*PreSeal
	for i := 0; i < sectors; i++ {
		sid := next
		next++

		pi, err := sb.AddPiece(context.TODO(), sid, nil, abi_spec.PaddedPieceSize(ssize).Unpadded(), rand.Reader)
		if err != nil {
			return nil, err
		}

		trand := sha256.Sum256(preimage)
		ticket := abi_spec.SealRandomness(trand[:])

		in2, err := sb.SealPreCommit1(context.TODO(), sid, ticket, []abi_spec.PieceInfo{pi})
		if err != nil {
			return nil, xerrors.Errorf("commit: %w", err)
		}

		scid, ucid, err := sb.SealPreCommit2(context.TODO(), sid, in2)
		if err != nil {
			return nil, xerrors.Errorf("commit: %w", err)
		}

		if err := sb.FinalizeSector(context.TODO(), sid); err != nil {
			return nil, xerrors.Errorf("trim cache: %w", err)
		}

		sealedSectors = append(sealedSectors, &PreSeal{
			CommR:     scid,
			CommD:     ucid,
			SectorID:  sid,
			ProofType: pt,
		})
	}

	return sealedSectors, nil
}

func createDeals(preseal []*PreSeal, worker, maddr address.Address, ssize abi_spec.SectorSize, start, end abi_spec.ChainEpoch) error {
	for _, sector := range preseal {
		proposal := &market_spec.DealProposal{
			PieceCID:             sector.CommD,
			PieceSize:            abi_spec.PaddedPieceSize(ssize),
			Client:               worker,
			Provider:             maddr,
			StartEpoch:           start,
			EndEpoch:             end,
			StoragePricePerEpoch: big_spec.Zero(),
			ProviderCollateral:   big_spec.Zero(),
			ClientCollateral:     big_spec.Zero(),
		}

		sector.Deal = *proposal
	}

	return nil
}

package tipset

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	sectorbuilder "github.com/filecoin-project/go-sectorbuilder"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	market_spec "github.com/filecoin-project/specs-actors/actors/builtin/market"
	miner_spec "github.com/filecoin-project/specs-actors/actors/builtin/miner"
	power_spec "github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	badger "github.com/ipfs/go-ds-badger2"
	logging "github.com/ipfs/go-log"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

var log = logging.Logger("chain-validation")

func init() {
	logging.SetAllLoggers(logging.LevelInfo)
}

// https://filecoin-project.github.io/specs/#full-miner-lifecycle

func TestHappyPathMinerStuff(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("happy path miner operation", func(t *testing.T) {
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
		// TODO this is probably the wrong value
		collateral := abi_spec.NewTokenAmount(1_000_000)

		sectorSize, err := abi_spec.RegisteredProof_StackedDRG2KiBSeal.SectorSize()
		require.NoError(t, err)

		// Step 0: Registration and Market participation
		createMinerMsg := td.MessageProducer.PowerCreateMiner(builtin_spec.StoragePowerActorAddr, minerOwner, power_spec.CreateMinerParams{
			Owner:      minerOwner,
			Worker:     minerWorker,
			SectorSize: sectorSize,
			Peer:       utils.RequireRandomPeerID(t),
		}, chain.Value(collateral), chain.Nonce(0))
		createMinerRct := types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createMinerRet), GasUsed: big_spec.Zero()}

		minerIdAddr := createMinerRet.IDAddress

		// Step 0.5: Add market funds from miner and client
		minerAddBalMsg := td.MessageProducer.MarketAddBalance(builtin_spec.StorageMarketActorAddr, minerWorker, minerIdAddr, chain.Nonce(0), chain.Value(collateral))
		minerAddBalRct := types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		clientAddBalMsg := td.MessageProducer.MarketAddBalance(builtin_spec.StorageMarketActorAddr, minerOwner, minerWorker, chain.Nonce(1), chain.Value(collateral))
		clientAddBalRct := types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		blkMsgs := chain.NewTipSetMessageBuilder().
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(createMinerMsg).
			WithBLSMessage(minerAddBalMsg).
			WithBLSMessage(clientAddBalMsg).
			Build()

		td.ExeCtx.Epoch++
		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 3)

		td.AssertReceipt(receipts[0], createMinerRct)
		td.AssertReceipt(receipts[1], minerAddBalRct)
		td.AssertReceipt(receipts[2], clientAddBalRct)

		m1temp, err := ioutil.TempDir("", "preseal")
		require.NoError(t, err)
		preseals, err := PreSealSectors(minerIdAddr, minerWorker, abi_spec.RegisteredProof_StackedDRG2KiBSeal, 0, 1, m1temp, []byte("some randomness"))
		require.NoError(t, err)
		require.NotNil(t, preseals)

		// sign the proposal
		buf, err := cborutil.Dump(&preseals[0].Deal)
		require.NoError(t, err)
		sig, err := td.Wallet().Sign(minerWorker, buf)
		require.NoError(t, err)

		// Step 0.75: Publish presealed deals
		dealID := abi_spec.DealID(0)
		dealIDs := []abi_spec.DealID{dealID}
		pubRet := chain.MustSerialize(&market_spec.PublishStorageDealsReturn{IDs: dealIDs})

		pubDealMsg := td.MessageProducer.MarketPublishStorageDeals(builtin_spec.StorageMarketActorAddr, minerWorker, market_spec.PublishStorageDealsParams{Deals: []market_spec.ClientDealProposal{{
			Proposal:        preseals[0].Deal,
			ClientSignature: sig,
		}}}, chain.Nonce(1), chain.Value(big_spec.Zero()))
		pubDealRct := types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: pubRet, GasUsed: big_spec.Zero()}

		// Step 1: Pre Committing Sectors
		preCommitMsg := td.MessageProducer.MinerPreCommitSector(minerIdAddr, minerWorker, miner_spec.SectorPreCommitInfo{
			RegisteredProof: preseals[0].ProofType,
			SectorNumber:    preseals[0].SectorID,
			SealedCID:       preseals[0].CommR,
			SealRandEpoch:   0,
			DealIDs:         dealIDs,
			Expiration:      preseals[0].Deal.EndEpoch,
		}, chain.Nonce(2), chain.Value(big_spec.Zero()))
		preCommitRct := types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

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
		proveCommitRct := types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		blkMsgs = chain.NewTipSetMessageBuilder().
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(proveCommitMsg).
			Build()

		// TODO don't use magic numbers
		td.ExeCtx.Epoch = 14

		receipts, err = td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)
		td.AssertReceipt(receipts[0], proveCommitRct)

		var pwrSt power_spec.State
		td.GetActorState(builtin_spec.StoragePowerActorAddr, &pwrSt)
		fmt.Println(pwrSt.TotalNetworkPower.String())
		fmt.Println(pwrSt.MinerCount)

		// a cron event type: CronEventWindowedPoStExpiration now exists in the power actors state
		// advance the epoch s.t. the miner misses the proving window
		transferMsg := td.MessageProducer.Transfer(minerOwner, minerOwner, chain.Nonce(2), chain.Value(big_spec.Zero()))
		transferRct := types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: []byte{}, GasUsed: big_spec.Zero()}

		blkMsgs = chain.NewTipSetMessageBuilder().
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(transferMsg).
			Build()

		td.ExeCtx.Epoch = 553
		receipts, err = td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)
		td.AssertReceipt(receipts[0], transferRct)

		transferMsg = td.MessageProducer.Transfer(minerOwner, minerOwner, chain.Nonce(3), chain.Value(big_spec.Zero()))
		blkMsgs = chain.NewTipSetMessageBuilder().
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(transferMsg).
			Build()
		td.ExeCtx.Epoch = 554
		receipts, err = td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)
		td.AssertReceipt(receipts[0], transferRct)

		transferMsg = td.MessageProducer.Transfer(minerOwner, minerOwner, chain.Nonce(4), chain.Value(big_spec.Zero()))
		blkMsgs = chain.NewTipSetMessageBuilder().
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(transferMsg).
			Build()
		td.ExeCtx.Epoch = 555
		receipts, err = td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)
		td.AssertReceipt(receipts[0], transferRct)

		// the miner should have been slashed, lets find out
		td.GetActorState(builtin_spec.StoragePowerActorAddr, &pwrSt)
		fmt.Println(pwrSt.TotalNetworkPower.String())
		fmt.Println(pwrSt.MinerCount)

	})
}

type PreSeal struct {
	CommR     cid.Cid
	CommD     cid.Cid
	SectorID  abi_spec.SectorNumber
	Deal      market_spec.DealProposal
	ProofType abi_spec.RegisteredProof
}

func PreSealSectors(maddr, worker address.Address, pt abi_spec.RegisteredProof, offset abi_spec.SectorNumber, sectors int, sbroot string, preimage []byte) ([]*PreSeal, error) {
	ppt, err := pt.RegisteredPoStProof()
	if err != nil {
		panic(err)
		return nil, err
	}

	spt, err := pt.RegisteredSealProof()
	if err != nil {
		panic(err)
		return nil, err
	}

	log.Infow("PreSealSectors", "maddr", maddr, "worker", worker)
	cfg := &sectorbuilder.Config{
		Miner:           maddr,
		SealProofType:   spt,
		PoStProofType:   ppt,
		FallbackLastNum: offset,
		Paths:           sectorbuilder.SimplePath(sbroot),
		WorkerThreads:   2,
	}

	if err := os.MkdirAll(sbroot, 0775); err != nil {
		panic(err)
		return nil, err
	}

	mds, err := badger.NewDatastore(filepath.Join(sbroot, "badger"), nil)
	if err != nil {
		panic(err)
		return nil, err
	}

	sb, err := sectorbuilder.New(cfg, namespace.Wrap(mds, datastore.NewKey("/sectorbuilder")))
	if err != nil {
		panic(err)
		return nil, err
	}

	ssize, err := pt.SectorSize()
	if err != nil {
		panic(err)
		return nil, err
	}

	var sealedSectors []*PreSeal
	for i := 0; i < sectors; i++ {
		sid, err := sb.AcquireSectorNumber()
		if err != nil {
			panic(err)
			return nil, err
		}

		pi, err := sb.AddPiece(context.TODO(), abi_spec.PaddedPieceSize(ssize).Unpadded(), sid, rand.Reader, nil)
		if err != nil {
			panic(err)
			return nil, err
		}

		trand := sha256.Sum256(preimage)
		ticket := abi_spec.SealRandomness(trand[:])

		log.Infof("sector-id: %d, piece info: %v", sid, pi)

		scid, ucid, err := sb.SealPreCommit(context.TODO(), sid, ticket, []abi_spec.PieceInfo{pi})
		if err != nil {
			panic(err)
			return nil, fmt.Errorf("commit: %w", err)
		}

		if err := sb.TrimCache(context.TODO(), sid); err != nil {
			panic(err)
			return nil, fmt.Errorf("trim cache: %w", err)
		}

		log.Infow("PreCommitOutput: ", "sectorID", sid, "sealedCID", scid, "unsealedCID", ucid)
		sealedSectors = append(sealedSectors, &PreSeal{
			CommR:     scid,
			CommD:     ucid,
			SectorID:  sid,
			ProofType: pt,
		})
	}

	if err := createDeals(sealedSectors, worker, maddr, ssize, 14, 9001); err != nil {
		panic(err)
		return nil, fmt.Errorf("creating deals: %w", err)
	}

	if err := mds.Close(); err != nil {
		panic(err)
		return nil, fmt.Errorf("closing datastore: %w", err)
	}

	return sealedSectors, nil
}

func createDeals(preseal []*PreSeal, worker, maddr address.Address, ssize abi_spec.SectorSize, start, end abi_spec.ChainEpoch) error {
	log.Infow("createDeals", "maddr", maddr, "worker", worker)
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

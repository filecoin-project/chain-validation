package suites

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgpwr"
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type strgpwrWrapper struct {
	T         testing.TB
	Driver    *StateDriver
	Producer  *chain.MessageProducer
	Validator *chain.Validator
	ExeCtx    *chain.ExecutionContext
}

func strgpwrTestSetup(t testing.TB, factory Factories) *strgpwrWrapper {
	drv := NewStateDriver(t, factory.NewState())
	gasPrice := types.NewInt(1)
	gasLimit := types.GasUnit(1000000)

	_, _, err := drv.State().SetSingletonActor(actors.InitAddress, types.NewInt(0))
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.BurntFundsAddress, types.NewInt(0))
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.NetworkAddress, TotalNetworkBalance)
	require.NoError(t, err)

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	testMiner := drv.NewAccountActor(0)
	exeCtx := chain.NewExecutionContext(1, testMiner)

	return &strgpwrWrapper{
		T:         t,
		Driver:    drv,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}

}

func StoragePowerActorConstructor(t testing.TB, factory Factories) {
	w := strgpwrTestSetup(t, factory)
	mustCreateStoragePowerActor(w)
}

func StoragePowerActorCreateStorageMiner(t testing.TB, factory Factories) {
	// 2,000,000,000,000,000,000,000,000
	const initialBal = "2000000000000000000000000"

	w := strgpwrTestSetup(t, factory)
	spAddr := mustCreateStoragePowerActor(w)

	alice := w.Driver.NewAccountActorBigBalance(types.NewIntFromString(initialBal))
	peerID0 := RequireIntPeerID(t, 0)
	minerAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	ms := strgpwr.NewMinerSet(t)
	ms.MinerSetAdd(minerAddr)

	mustCreateStorageMiner(w, 0, strgpwr.SectorSizes[0], types.NewIntFromString("1999999995415053581179420"), minerAddr, alice, alice, alice, peerID0)
	w.Driver.AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     1,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(0),
	})
}

func StoragePowerActorUpdateStorage(t testing.TB, factory Factories) {
	// 2,000,000,000,000,000,000,000,000
	const initialBal = "2000000000000000000000000"
	const updateSize = 100
	const nextPpEnd = uint64(10)
	const prevPpEnd = uint64(0)

	w := strgpwrTestSetup(t, factory)
	spAddr := mustCreateStoragePowerActor(w)

	alice := w.Driver.NewAccountActorBigBalance(types.NewIntFromString(initialBal))
	peerID0 := RequireIntPeerID(t, 0)
	minerAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	ms := strgpwr.NewMinerSet(t)
	// calculate the cid of Miners
	ms.MinerSetAdd(minerAddr)

	mustCreateStorageMiner(w, 0, strgpwr.SectorSizes[0], types.NewIntFromString("1999999995415053581179420"), minerAddr, alice, alice, alice, peerID0)
	w.Driver.AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     1,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(0),
	})

	mustUpdateStoragePower(w, 0, 0, updateSize, nextPpEnd, prevPpEnd, minerAddr)
	// calculate the cid of ProvingBuckets
	ms.CalculateBuckets(minerAddr, prevPpEnd, nextPpEnd)
	w.Driver.AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     1,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(updateSize),
	})
}

func mustUpdateStoragePower(w *strgpwrWrapper, nonce, value, delta, nextPpEnd, prevPpEnd uint64, minerAddr address.Address) {
	msg, err := w.Producer.StoragePowerUpdateStorage(minerAddr, nonce, types.NewInt(delta), nextPpEnd, prevPpEnd, chain.Value(value))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)
	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})

}

func mustCreateStorageMiner(w *strgpwrWrapper, nonce, sectorSize uint64, value types.BigInt, minerAddr, from, owner, worker address.Address, pid peer.ID) {
	msg, err := w.Producer.StoragePowerCreateStorageMiner(from, nonce, owner, worker, sectorSize, pid, chain.BigValue(value))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)

	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: minerAddr.Bytes(),
		GasUsed:     0,
	})
}

func mustCreateStoragePowerActor(w *strgpwrWrapper) address.Address {
	// Storage Power Actor is a singleton actor, requires special setup.
	_, _, err := w.Driver.State().SetSingletonActor(actors.StoragePowerAddress, types.NewInt(0))
	require.NoError(w.T, err)

	ms := strgpwr.NewMinerSet(w.T)
	spAddr := w.Producer.SingletonAddress(actors.StoragePowerAddress)
	w.Driver.AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     0,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(0),
	})
	return spAddr
}

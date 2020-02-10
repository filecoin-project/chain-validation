package suites

// TODO uncomment when ready to implement
/*
import (
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgpwr"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

func StoragePowerActorConstructor(t testing.TB, factory Factories) {
	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:       types.NewInt(0),
		actors.BurntFundsAddress: types.NewInt(0),
		actors.NetworkAddress:    TotalNetworkBalance,
	})
	mustCreateStoragePowerActor(td)
}

func StoragePowerActorCreateStorageMiner(t testing.TB, factory Factories) {
	// 2,000,000,000,000,000,000,000,000
	const initialBal = "2000000000000000000000000"

	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:       types.NewInt(0),
		actors.BurntFundsAddress: types.NewInt(0),
		actors.NetworkAddress:    TotalNetworkBalance,
	})
	spAddr := mustCreateStoragePowerActor(td)

	alice := td.Driver().NewAccountActorBigBalance(types.NewIntFromString(initialBal))
	peerID0 := RequireIntPeerID(t, 0)
	minerAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	ms := strgpwr.NewMinerSet(t)
	ms.MinerSetAdd(minerAddr)

	mustCreateStorageMiner(td, 0, strgpwr.SectorSizes[0], types.NewIntFromString("1999999995415053581179420"), minerAddr, alice, alice, alice, peerID0)
	td.Driver().AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
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

	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:       types.NewInt(0),
		actors.BurntFundsAddress: types.NewInt(0),
		actors.NetworkAddress:    TotalNetworkBalance,
	})
	spAddr := mustCreateStoragePowerActor(td)

	alice := td.Driver().NewAccountActorBigBalance(types.NewIntFromString(initialBal))
	peerID0 := RequireIntPeerID(t, 0)
	minerAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	ms := strgpwr.NewMinerSet(t)
	// calculate the cid of Miners
	ms.MinerSetAdd(minerAddr)

	mustCreateStorageMiner(td, 0, strgpwr.SectorSizes[0], types.NewIntFromString("1999999995415053581179420"), minerAddr, alice, alice, alice, peerID0)
	td.Driver().AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     1,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(0),
	})

	mustUpdateStoragePower(td, 0, 0, updateSize, nextPpEnd, prevPpEnd, minerAddr)
	// calculate the cid of ProvingBuckets
	ms.CalculateBuckets(minerAddr, prevPpEnd, nextPpEnd)
	td.Driver().AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     1,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(updateSize),
	})
}

func mustUpdateStoragePower(td TestDriver, nonce, value, delta, nextPpEnd, prevPpEnd uint64, minerAddr address.Address) {
	msg, err := td.Producer.StoragePowerUpdateStorage(minerAddr, nonce, types.NewInt(delta), nextPpEnd, prevPpEnd, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver().State(), msg)
	require.NoError(td.TB(), err)
	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})

}

func mustCreateStorageMiner(td TestDriver, nonce, sectorSize uint64, value types.BigInt, minerAddr, from, owner, worker address.Address, pid peer.ID) {
	msg, err := td.Producer.StoragePowerCreateStorageMiner(from, nonce, owner, worker, sectorSize, pid, chain.BigValue(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: minerAddr.Bytes(),
		GasUsed:     types.NewInt(0),
	})
}

func mustCreateStoragePowerActor(td TestDriver) address.Address {
	// Storage Power Actor is a singleton actor, requires special setup.
	_, _, err := td.Driver().State().SetSingletonActor(actors.StoragePowerAddress, types.NewInt(0))
	require.NoError(td.TB(), err)

	ms := strgpwr.NewMinerSet(td.TB())
	spAddr := td.Producer.SingletonAddress(actors.StoragePowerAddress)
	td.Driver().AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     0,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(0),
	})
	return spAddr
}
*/

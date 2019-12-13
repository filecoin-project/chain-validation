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

func StoragePowerActorConstructor(t testing.TB, factory Factories) {
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:       types.NewInt(0),
		actors.BurntFundsAddress: types.NewInt(0),
		actors.NetworkAddress:    TotalNetworkBalance,
	})
	mustCreateStoragePowerActor(c)
}

func StoragePowerActorCreateStorageMiner(t testing.TB, factory Factories) {
	// 2,000,000,000,000,000,000,000,000
	const initialBal = "2000000000000000000000000"

	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:       types.NewInt(0),
		actors.BurntFundsAddress: types.NewInt(0),
		actors.NetworkAddress:    TotalNetworkBalance,
	})
	spAddr := mustCreateStoragePowerActor(c)

	alice := c.Driver().NewAccountActorBigBalance(types.NewIntFromString(initialBal))
	peerID0 := RequireIntPeerID(t, 0)
	minerAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	ms := strgpwr.NewMinerSet(t)
	ms.MinerSetAdd(minerAddr)

	mustCreateStorageMiner(c, 0, strgpwr.SectorSizes[0], types.NewIntFromString("1999999995415053581179420"), minerAddr, alice, alice, alice, peerID0)
	c.Driver().AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
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

	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:       types.NewInt(0),
		actors.BurntFundsAddress: types.NewInt(0),
		actors.NetworkAddress:    TotalNetworkBalance,
	})
	spAddr := mustCreateStoragePowerActor(c)

	alice := c.Driver().NewAccountActorBigBalance(types.NewIntFromString(initialBal))
	peerID0 := RequireIntPeerID(t, 0)
	minerAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	ms := strgpwr.NewMinerSet(t)
	// calculate the cid of Miners
	ms.MinerSetAdd(minerAddr)

	mustCreateStorageMiner(c, 0, strgpwr.SectorSizes[0], types.NewIntFromString("1999999995415053581179420"), minerAddr, alice, alice, alice, peerID0)
	c.Driver().AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     1,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(0),
	})

	mustUpdateStoragePower(c, 0, 0, updateSize, nextPpEnd, prevPpEnd, minerAddr)
	// calculate the cid of ProvingBuckets
	ms.CalculateBuckets(minerAddr, prevPpEnd, nextPpEnd)
	c.Driver().AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     1,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(updateSize),
	})
}

func mustUpdateStoragePower(c Candy, nonce, value, delta, nextPpEnd, prevPpEnd uint64, minerAddr address.Address) {
	msg, err := c.Producer().StoragePowerUpdateStorage(minerAddr, nonce, types.NewInt(delta), nextPpEnd, prevPpEnd, chain.Value(value))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})

}

func mustCreateStorageMiner(c Candy, nonce, sectorSize uint64, value types.BigInt, minerAddr, from, owner, worker address.Address, pid peer.ID) {
	msg, err := c.Producer().StoragePowerCreateStorageMiner(from, nonce, owner, worker, sectorSize, pid, chain.BigValue(value))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)

	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: minerAddr.Bytes(),
		GasUsed:     types.NewInt(0),
	})
}

func mustCreateStoragePowerActor(c Candy) address.Address {
	// Storage Power Actor is a singleton actor, requires special setup.
	_, _, err := c.Driver().State().SetSingletonActor(actors.StoragePowerAddress, types.NewInt(0))
	require.NoError(c.TB(), err)

	ms := strgpwr.NewMinerSet(c.TB())
	spAddr := c.Producer().SingletonAddress(actors.StoragePowerAddress)
	c.Driver().AssertStoragePowerState(spAddr, strgpwr.StoragePowerState{
		Miners:         ms.MinerCid,
		ProvingBuckets: ms.ProvingBucketsCid,
		MinerCount:     0,
		LastMinerCheck: 0,
		TotalStorage:   types.NewInt(0),
	})
	return spAddr
}

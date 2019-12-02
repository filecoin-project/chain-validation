package suites

import (
	"encoding/binary"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgminr"
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

const (
	totalFilecoin     = 2000000000
	filecoinPrecision = 1000000000000000000
)

var (
	TotalNetworkBalance = types.NewInt(types.NewInt(1).Mul(types.NewInt(totalFilecoin).Int, types.NewInt(0).SetUint64(filecoinPrecision)).Uint64())

	sectorSizes = []uint64{
		16 << 20,
		256 << 20,
		1 << 30,
	}
)

func testSetup(t testing.TB, factory Factories) (*StateDriver, types.BigInt, types.GasUnit) {
	drv := NewStateDriver(t, factory.NewState())
	gasPrice := types.NewInt(1)
	gasLimit := types.GasUnit(1000000000000)

	_, _, err := drv.State().SetSingletonActor(actors.InitAddress, types.NewInt(0))
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.NetworkAddress, TotalNetworkBalance)
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.StoragePowerAddress, types.NewInt(0))
	require.NoError(t, err)

	return drv, gasPrice, gasLimit

}

func CreateStorageMinerAndUpdatePeerID(t testing.TB, factory Factories) {
	drv, gasPrice, gasLimit := testSetup(t, factory)

	// miner that mines in this test
	testMiner := drv.NewAccountActor(0)
	// account that will own the miner
	minerOwner := drv.NewAccountActorBigBalance(types.NewIntFromString("2000000000000000000000000"))

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)
	exeCtx := chain.NewExecutionContext(1, testMiner)

	//
	// create a storage miner
	//

	// sector size of the miner created
	sectorSize := types.NewInt(sectorSizes[0])
	// peerID of the miner created
	peerID := RequireIntPeerID(t, 1)

	msg, err := producer.StoragePowerCreateStorageMiner(minerOwner, 0, minerOwner, minerOwner, sectorSize.Uint64(), peerID, chain.BigValue(types.NewIntFromString("1999999995415053581179420")))
	require.NoError(t, err)
	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: []byte{0, 102},
		GasUsed:     0,
	})

	//
	// verify storage miners sector size
	//

	// address of the miner created by the above message
	minerAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	msg, err = producer.StorageMinerGetSectorSize(minerAddr, minerOwner, 1, chain.Value(2000000))
	require.NoError(t, err)
	msgReceipt, err = validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: sectorSize.Bytes(),
		GasUsed:     0,
	})

	//
	// verify storage miners owner
	//
	msg, err = producer.StorageMinerGetOwner(minerAddr, minerOwner, 2, chain.Value(2000000))
	require.NoError(t, err)
	msgReceipt, err = validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: minerOwner.Bytes(),
		GasUsed:     0,
	})

	//
	// verify storage miners power
	//
	msg, err = producer.StorageMinerGetPower(minerAddr, minerOwner, 3, chain.Value(2000000))
	require.NoError(t, err)
	msgReceipt, err = validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: []byte{},
		GasUsed:     0,
	})

	//
	// verify storage miner worker address
	//
	msg, err = producer.StorageMinerGetWorkerAddr(minerAddr, minerOwner, 4, chain.Value(2000000))
	require.NoError(t, err)
	msgReceipt, err = validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: minerOwner.Bytes(),
		GasUsed:     0,
	})

	//
	// verify storage miner peerID
	//
	msg, err = producer.StorageMinerGetPeerID(minerAddr, minerOwner, 5, chain.Value(2000000))
	require.NoError(t, err)
	msgReceipt, err = validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: []byte(peerID),
		GasUsed:     0,
	})

	//
	// update peerID
	//

	//peerID to update miner with
	peerID2 := RequireIntPeerID(t, 2)

	msg, err = producer.StorageMinerUpdatePeerID(minerAddr, minerOwner, 6, peerID2, chain.Value(2000000))
	require.NoError(t, err)
	msgReceipt, err = validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: []byte(nil),
		GasUsed:     0,
	})

	//
	// verify storage miner peerID
	//
	msg, err = producer.StorageMinerGetPeerID(minerAddr, minerOwner, 7, chain.Value(2000000))
	require.NoError(t, err)
	msgReceipt, err = validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: []byte(peerID2),
		GasUsed:     0,
	})

	//
	// verify storage miners state
	//
	minerActor, err := drv.State().Actor(minerAddr)
	require.NoError(t, err)

	minerActorStorage, err := drv.State().Storage(minerAddr)
	require.NoError(t, err)

	var minerState strgminr.StorageMinerActorState
	require.NoError(t, minerActorStorage.Get(minerActor.Head(), &minerState))

	var minerInfo strgminr.MinerInfo
	require.NoError(t, minerActorStorage.Get(minerState.Info, &minerInfo))

	drv.AssertMinerInfo(minerInfo, strgminr.MinerInfo{
		Owner:      minerOwner,
		Worker:     minerOwner,
		PeerID:     peerID2,
		SectorSize: sectorSize.Uint64(),
	})

}

// RequireIntPeerID takes in an integer and creates a unique peer id for it.
func RequireIntPeerID(t testing.TB, i int64) peer.ID {
	buf := make([]byte, 16)
	n := binary.PutVarint(buf, i)
	h, err := mh.Sum(buf[:n], mh.ID, -1)
	require.NoError(t, err)
	pid, err := peer.IDFromBytes(h)
	require.NoError(t, err)
	return pid
}

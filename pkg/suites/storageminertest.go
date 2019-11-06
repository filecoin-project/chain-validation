package suites

import (
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
)

const totalFilecoin = 2000000000
const filecoinPrecision = 1000000000000000000

var sectorSizes = []uint64{
	16 << 20,
	256 << 20,
	1 << 30,
}

func CreateStorageMinerAndUpdatePeerIDTest(t testing.TB, factory Factories) {
	drv := NewStateDriver(t, factory.NewState())

	gasPrice := big.NewInt(1)
	// gas prices will be inconsistent for a while, use a big value lotus team suggests using a large value here.
	gasLimit := state.GasUnit(1000000)
	TotalNetworkBalance := big.NewInt(1).Mul(big.NewInt(totalFilecoin), big.NewInt(0).SetUint64(filecoinPrecision))
	_, _, err := drv.State().SetSingletonActor(state.InitAddress, big.NewInt(0))
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(state.NetworkAddress, TotalNetworkBalance)
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(state.StoragePowerAddress, big.NewInt(0))
	require.NoError(t, err)


	// miner that mines in this test
	testMiner := drv.NewAccountActor(0)
	// account that will own the miner
	minerOwner := drv.NewAccountActor(20000000000)

	// address of the miner created
	minerAddr, err := state.NewIDAddress(102)
	require.NoError(t, err)
	// sector size of the miner created
	sectorSize := big.NewInt(int64(sectorSizes[0]))
	// peerID of the miner created
	rawPeerID, err := RequireIntPeerID(t, 1).MarshalBinary()
	require.NoError(t,err)
	peerID := state.PeerID(rawPeerID)
	// peerID of the miner after update
	rawPeerID2, err := RequireIntPeerID(t,2).MarshalBinary()
	require.NoError(t,err)
	peerID2 :=state.PeerID(rawPeerID2)

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)
	exeCtx := chain.NewExecutionContext(1, testMiner)

	//
	// create a storage miner
	//
	msg, err := producer.StoragePowerCreateStorageMiner(minerOwner, 0, minerOwner, minerOwner, sectorSize, peerID, chain.Value(2000000))
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
		ReturnValue: []byte(minerOwner),
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
		ReturnValue: []byte(minerOwner),
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

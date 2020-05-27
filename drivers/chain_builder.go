package drivers

import (
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
	"testing"
)

type ChainBuilder struct {
	t testing.TB

	bs     blockstore.Blockstore
	cstore cbor.IpldStore
}

func NewChainBuilder(t testing.TB, bs blockstore.Blockstore) *ChainBuilder {
	cb := &ChainBuilder{
		t: t,

		bs:     bs,
		cstore: cbor.NewCborStore(bs),
	}
	return cb
}

type FullBlockBuilder struct {
	t testing.TB

	stamper      TimeStamper
	stateBuilder StateBuilder
}

func (fb *FullBlockBuilder) Build(parents *types.TipSet, miner address.Address, height abi.ChainEpoch,
	vrfticket *types.Ticket, eticket *types.ElectionProof, beacons []types.BeaconEntry, wpost []abi.PoStProof,
	msgs []*types.SignedMessage) *types.FullBlock {

}

func (fb *FullBlockBuilder) CreateBlock(parent types.TipSet, nullRounds abi.ChainEpoch, miner address.Address,
	ticket *types.Ticket, wProof []abi.PoStProof, beacons []types.BeaconEntry, smsgs []*types.SignedMessage) *types.FullBlock {

	header := &types.BlockHeader{
		Miner:   miner,
		Parents: parent.Cids(),
		Ticket:  ticket,

		ElectionProof: nil, // FIXME: epost is dead, right?

		Timestamp: fb.stamper.Stamp(),

		BeaconEntries:         beacons,
		WinPoStProof:          wProof,
		Height:                0,
		ParentStateRoot:       cid.Cid{},
		ParentMessageReceipts: cid.Cid{},
		Messages:              cid.Cid{},
		BLSAggregate:          nil,
		BlockSig:              nil,
		ForkSignaling:         0,

		ParentWeight: big.Int{}, // compute below
	}

	block := &types.FullBlock{
		Header:        nil,
		BlsMessages:   nil,
		SecpkMessages: nil,
	}

}

// TimeStamper is an object that timestamps blocks
type TimeStamper interface {
	Stamp(abi.ChainEpoch) uint64
}

// StateBuilder abstracts the computation of state root CIDs from the chain builder.
type StateBuilder interface {
	ComputeState(prev cid.Cid, blsMessages [][]*types.Message, secpMessages [][]*types.SignedMessage) (cid.Cid, []*types.MessageReceipt, error)
	Weigh(tip types.TipSet, state cid.Cid) (big.Int, error)
}

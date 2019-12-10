package strgmrkt

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-amt-ipld"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-hamt-ipld"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type SerializationMode = uint64

const (
	SerializationUnixFSv0 = iota
)

type StorageDeal struct {
	Proposal         StorageDealProposal
	CounterSignature *types.Signature
}

type OnChainDeal struct {
	Deal            StorageDeal
	ActivationEpoch uint64 // 0 = inactive
}

type StorageMarketState struct {
	Balances cid.Cid
	Deals    cid.Cid

	NextDealID uint64
}

type StorageDealProposal struct {
	PieceRef           []byte // cid bytes
	PieceSize          uint64
	PieceSerialization SerializationMode

	Client   address.Address
	Provider address.Address

	ProposalExpiration uint64
	Duration           uint64

	StoragePricePerEpoch types.BigInt
	StorageCollateral    types.BigInt

	ProposerSignature *types.Signature
}

type StorageParticipantBalance struct {
	Locked    types.BigInt
	Available types.BigInt
}

type PublishStorageDealResponse struct {
	DealIDs []uint64
}

//
// Message Method  Params
//

type WithdrawBalanceParams struct {
	Balance types.BigInt
}

type PublishStorageDealsParams struct {
	Deals []StorageDeal
}

type ActivateStorageDealsParams struct {
	Deals []uint64
}

type ComputeDataCommitmentParams struct {
	DealIDs    []uint64
	SectorSize uint64
}

type ProcessStorageDealsPaymentParams struct {
	DealIDs []uint64
}

//
// Helper methods for calculating market deal and balance cid's
//

type MarketTracker struct {
	hamtStore *hamt.CborIpldStore
	Balance   cid.Cid

	amtStore amt.Blocks
	Deals    cid.Cid

	bs blockstore.Blockstore

	T testing.TB
}

func NewMarketTracker(t testing.TB) *MarketTracker {
	mds := ds.NewMapDatastore()
	bs := blockstore.NewBlockstore(mds)

	s := hamt.CSTFromBstore(bs)
	nd := hamt.NewNode(s)
	c, err := s.Put(context.Background(), nd)
	require.NoError(t, err)

	blks := amt.WrapBlockstore(bs)
	emptyamt, err := amt.FromArray(blks, nil)
	require.NoError(t, err)

	return &MarketTracker{
		hamtStore: s,
		Balance:   c,
		amtStore:  blks,
		Deals:     emptyamt,
		bs:        bs,
		T:         t,
	}
}

func (m *MarketTracker) SetMarketBalances(whom map[address.Address]StorageParticipantBalance) {
	ctx := context.Background()

	nd, err := hamt.LoadNode(ctx, m.hamtStore, m.Balance)
	require.NoError(m.T, err)

	for addr, b := range whom {
		balance := b // to stop linter complaining
		err = nd.Set(ctx, string(addr.Bytes()), &balance)
		require.NoError(m.T, err)
	}
	err = nd.Flush(ctx)
	require.NoError(m.T, err)

	c, err := m.hamtStore.Put(ctx, nd)
	require.NoError(m.T, err)

	m.Balance = c
}

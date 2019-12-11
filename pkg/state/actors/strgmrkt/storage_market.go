package strgmrkt

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/filecoin-project/go-amt-ipld"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-hamt-ipld"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/state"
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

func (sd *StorageDeal) Sign(ctx context.Context, from address.Address, s state.Signer) error {
	var buf bytes.Buffer
	if err := sd.Proposal.MarshalCBOR(&buf); err != nil {
		return err
	}
	sig, err := s.Sign(ctx, from, buf.Bytes())
	if err != nil {
		return err
	}
	sd.CounterSignature = sig
	return nil
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

func (sdp *StorageDealProposal) TotalStoragePrice() types.BigInt {
	return types.BigInt{big.NewInt(0).Mul(sdp.StoragePricePerEpoch.Int, big.NewInt(0).SetUint64(sdp.Duration))}
}

func (sdp *StorageDealProposal) Sign(ctx context.Context, from address.Address, s state.Signer) error {
	if sdp.ProposerSignature != nil {
		return errors.New("signature already present in StorageDealProposal")
	}
	var buf bytes.Buffer
	if err := sdp.MarshalCBOR(&buf); err != nil {
		return err
	}
	sig, err := s.Sign(ctx, from, buf.Bytes())
	if err != nil {
		return err
	}
	sdp.ProposerSignature = sig
	return nil
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
// Helper methods for calculating market balance cid's
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

// TODO add support for calculating market deals cid's
/*
func (m *MarketTracker) SetMarketDeals() {
}
*/

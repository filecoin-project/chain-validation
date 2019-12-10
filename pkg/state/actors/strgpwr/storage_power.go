package strgpwr

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/filecoin-project/go-amt-ipld"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-hamt-ipld"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

const SlashablePowerDelay = 200

var (
	SectorSizes = []uint64{
		16 << 20,
		256 << 20,
		1 << 30,
	}
)

type StoragePowerState struct {
	Miners         cid.Cid
	ProvingBuckets cid.Cid // amt[ProvingPeriodBucket]hamt[minerAddress]struct{}
	MinerCount     uint64
	LastMinerCheck uint64

	TotalStorage types.BigInt
}

type CreateStorageMinerParams struct {
	Owner      address.Address
	Worker     address.Address
	SectorSize uint64
	PeerID     peer.ID
}

type UpdateStorageParams struct {
	Delta                    types.BigInt
	NextProvingPeriodEnd     uint64
	PreviousProvingPeriodEnd uint64
}

type PowerLookupParams struct {
	Miner address.Address
}

type PledgeCollateralParams struct {
	Size types.BigInt
}

//
// Helpers for calculating power actor's Miner CID and ProvingBuckets CID
//

// used to calculate storage power state Miner and ProvingBucket Cid's.
func NewMinerSet(tb testing.TB) *minerSet {
	mds := ds.NewMapDatastore()
	bs := blockstore.NewBlockstore(mds)

	s := hamt.CSTFromBstore(bs)
	nd := hamt.NewNode(s)
	c, err := s.Put(context.Background(), nd)
	require.NoError(tb, err)

	blks := amt.WrapBlockstore(bs)
	emptyamt, err := amt.FromArray(blks, nil)
	require.NoError(tb, err)

	return &minerSet{
		hamtStore:         s,
		MinerCid:          c,
		amtStore:          blks,
		ProvingBucketsCid: emptyamt,
		bs:                bs,
		T:                 tb,
	}
}

type minerSet struct {
	hamtStore *hamt.CborIpldStore
	MinerCid  cid.Cid

	amtStore          amt.Blocks
	ProvingBucketsCid cid.Cid

	bs blockstore.Blockstore

	T testing.TB
}

func (m *minerSet) MinerSetAdd(maddr address.Address) {
	ctx := context.Background()
	nd, err := hamt.LoadNode(ctx, m.hamtStore, m.MinerCid)
	require.NoError(m.T, err)

	mkey := string(maddr.Bytes())
	err = nd.Find(ctx, mkey, nil)
	require.NotNil(m.T, err)
	require.Equal(m.T, hamt.ErrNotFound, err)

	err = nd.Set(ctx, mkey, uint64(1))
	require.NoError(m.T, err)

	err = nd.Flush(ctx)
	require.NoError(m.T, err)

	c, err := m.hamtStore.Put(ctx, nd)
	require.NoError(m.T, err)

	m.MinerCid = c
}

func (m *minerSet) MinerSetRemove(maddr address.Address) {
	ctx := context.Background()
	nd, err := hamt.LoadNode(ctx, m.hamtStore, m.MinerCid)
	require.NoError(m.T, err)

	mkey := string(maddr.Bytes())
	err = nd.Delete(ctx, mkey)
	require.NoError(m.T, err)

	err = nd.Flush(ctx)
	require.NoError(m.T, err)

	c, err := m.hamtStore.Put(ctx, nd)
	require.NoError(m.T, err)

	m.MinerCid = c
}

func (m *minerSet) CalculateBuckets(maddr address.Address, pppe, nppe uint64) {
	previousBucket := pppe % SlashablePowerDelay
	nextBucket := nppe % SlashablePowerDelay

	if previousBucket == nextBucket && pppe != 0 {
		return // noop
	}

	buckets, err := amt.LoadAMT(m.amtStore, m.ProvingBucketsCid)
	require.NoError(m.T, err)

	if pppe != 0 { // delete from previous bucket
		m.deleteMinerFromBucket(maddr, buckets, previousBucket)
	}

	m.addMinerToBucket(maddr, buckets, nextBucket)

	newBucketCid, err := buckets.Flush()
	require.NoError(m.T, err)

	m.ProvingBucketsCid = newBucketCid
}

func (m *minerSet) addMinerToBucket(minerAddr address.Address, buckets *amt.Root, nextBucket uint64) {
	ctx := context.Background()
	var bhamt *hamt.Node
	var bucket cid.Cid

	err := buckets.Get(nextBucket, &bucket)
	switch err.(type) {
	case *amt.ErrNotFound:
		bhamt = hamt.NewNode(m.hamtStore)
	case nil:
		bhamt, err = hamt.LoadNode(ctx, m.hamtStore, bucket)
		require.NoError(m.T, err)
	default:
		require.Fail(m.T, "getting proving bucket")
	}

	err = bhamt.Set(ctx, string(minerAddr.Bytes()), cborNull)
	require.NoError(m.T, err)

	err = bhamt.Flush(ctx)
	require.NoError(m.T, err)

	bucket, err = m.hamtStore.Put(ctx, bhamt)
	require.NoError(m.T, err)

	err = buckets.Set(nextBucket, bucket)

	require.NoError(m.T, err)
	m.ProvingBucketsCid = bucket
}

func (m *minerSet) deleteMinerFromBucket(minerAddr address.Address, buckets *amt.Root, previousBucket uint64) cid.Cid {
	ctx := context.Background()
	var bucket cid.Cid
	err := buckets.Get(previousBucket, &bucket)
	switch err.(type) {
	case *amt.ErrNotFound:
		require.NoError(m.T, err)
	case nil: // noop
	default:
		require.Fail(m.T, "failed to get bucket")
	}

	bhamt, err := hamt.LoadNode(ctx, m.hamtStore, bucket)
	require.NoError(m.T, err)

	err = bhamt.Delete(ctx, string(minerAddr.Bytes()))
	require.NoError(m.T, err)

	err = bhamt.Flush(ctx)
	require.NoError(m.T, err)

	bucket, err = m.hamtStore.Put(ctx, bhamt)
	require.NoError(m.T, err)

	err = buckets.Set(previousBucket, bucket)
	require.NoError(m.T, err)

	m.ProvingBucketsCid = bucket
	return bucket
}

type cbgNull struct{}

var cborNull = &cbgNull{}

func (cbgNull) MarshalCBOR(w io.Writer) error {
	n, err := w.Write(cbg.CborNull)
	if err != nil {
		panic(err)
	}
	if n != 1 {
		panic("expected to write 1 byte")
	}
	return nil
}

func (cbgNull) UnmarshalCBOR(r io.Reader) error {
	b := [1]byte{}
	n, err := r.Read(b[:])
	if err != nil {
		panic(err)
	}
	if n != 1 {
		panic("expected 1 byte")
	}
	if !bytes.Equal(b[:], cbg.CborNull) {
		panic("expected cbor null")
	}
	return nil
}

package state_test

import (
	"encoding/binary"
	"math/big"
	"testing"

	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/libp2p/go-libp2p-core/peer"
	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

func TestExampleEncodeValues(t *testing.T) {
	owner, err := state.NewActorAddress([]byte{1,2,3,4,5})
	require.NoError(t, err)

	sectorSize := state.BytesAmount(big.NewInt(10))

	bpid, err := RequireIntPeerID(t, 0).MarshalBinary()
	require.NoError(t, err)

	peerID := state.PeerID(bpid)

	params := []interface{}{owner, sectorSize, peerID}
	data, err := state.EncodeValues(params...)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// since decoding logic is only need for testing in package it is simplest to do it manually here.
	var arr [][]byte
	require.NoError(t, cbor.DecodeInto(data, &arr))

	v := arr[0] // owner address
	expOwner := state.Address(v)
	assert.Equal(t, owner, expOwner)

	v = arr[1] // sectorSize
	expSectorSize := state.BytesAmount(big.NewInt(0).SetBytes(v))
	assert.Equal(t, sectorSize, expSectorSize)

	v = arr[2] // peerID
	bp, err := peer.IDFromBytes(v)
	require.NoError(t, err)
	expPeerID := state.PeerID(bp)
	assert.Equal(t, peerID, expPeerID)

}

func RequireIntPeerID(t *testing.T, i int64) peer.ID {
	buf := make([]byte, 16)
	n := binary.PutVarint(buf, i)
	h, err := mh.Sum(buf[:n], mh.ID, -1)
	require.NoError(t, err)
	pid, err := peer.IDFromBytes(h)
	require.NoError(t, err)
	return pid
}

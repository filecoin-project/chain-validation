package state_test

import (
	"encoding/binary"
	"math/big"
	"testing"

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

	stuff, err := state.EncodeValues(params...)
	assert.NoError(t, err)
	assert.NotEmpty(t, stuff)
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



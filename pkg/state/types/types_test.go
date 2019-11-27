package types_test

import (
	"math/big"
	"testing"

	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/state/address"
)

func TestExampleEncodeValues(t *testing.T) {
	owner, err := address.NewActorAddress([]byte{1, 2, 3, 4, 5})
	require.NoError(t, err)

	sectorSize := BytesAmount(big.NewInt(10))
	peerID := PeerID([]byte{0, 9, 8, 7, 6, 5})

	params := []interface{}{owner, sectorSize, peerID}
	data, err := EncodeValues(params...)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// since decoding logic is only need for testing in package it is simplest to do it manually here.
	var arr []interface{}
	require.NoError(t, cbor.DecodeInto(data, &arr))

	v := arr[0] // owner address
	expOwner := v.([]byte)
	assert.Equal(t, owner.Bytes(), expOwner)

	v = arr[1] // sectorSize
	expSectorSize := v.(uint64)
	actualSectorSize := big.Int(*sectorSize)
	assert.Equal(t, actualSectorSize.Uint64(), expSectorSize)

	v = arr[2] // peerID
	expPeerID := v.(string)
	assert.Equal(t, string(peerID), expPeerID)

}

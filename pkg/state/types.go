package state

import (
	"math/big"
)

// Type aliases for state values and message method parameters.
type (
	BytesAmount *big.Int
	AttoFIL *big.Int
	GasUnit uint64
	PubKey []byte
	PeerID []byte
)


package state

import "math/big"

// MethodID is the type of a VM actor method identifier.
// This is a string for generality at the moment, but should eventually become an integer.
type MethodID string

// Type aliases for message method parameters.
type (
	BytesAmount *big.Int
	AttoFIL *big.Int
	GasUnit uint64
	PeerID string
)

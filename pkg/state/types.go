package state

import "math/big"

// Type aliases for message method parameters.
type (
	BytesAmount *big.Int
	AttoFIL *big.Int
	GasUnit uint64
	PubKey []byte
	PeerID string
)

type MessageReceiept struct {
	Exitcode    uint8
	ReturnValue [][]byte
	GasUsed     AttoFIL
}

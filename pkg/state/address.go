package state

import (
	"github.com/filecoin-project/go-leb128"
	"github.com/minio/blake2b-simd"
)

type Address string

// NewIDAddress builds a new ID-type address.
func NewIDAddress(id uint64) (Address, error) {
	return newAddress(0, leb128.FromUInt64(id))
}

// NewActorAddress builds a new Actor address.
func NewActorAddress(data []byte) (Address, error) {
	digest, err := hash(data, addressHashConfig)
	if err != nil {
		return "", err
	}
	return newAddress(2, digest)
}

func NewSecp256k1Address(pubkey []byte) (Address, error) {
	digest, err := hash(pubkey)
	if err != nil {
		return "", err
	}
	return newAddress(1,digest)
}

func newAddress(protocol byte, payload []byte) (Address, error) {
	buf := make([]byte, 1+len(payload))
	buf[0] = protocol
	copy(buf[1:], payload)

	return Address(buf), nil
}

func hash(ingest []byte, cfg *blake2b.Config) ([]byte, error) {
	hasher, err := blake2b.New(cfg)
	if err != nil {
		return nil, err
	}
	if _, err := hasher.Write(ingest); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

var addressHashConfig = &blake2b.Config{Size: 20}

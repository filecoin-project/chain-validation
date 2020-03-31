package utils

import (
	"encoding/binary"
	"testing"

	addr "github.com/filecoin-project/go-address"
	"github.com/multiformats/go-varint"
)

// If you use this method while writing a test you are more than likely doing something wrong.
func NewIDAddr(t testing.TB, id uint64) addr.Address {
	address, err := addr.NewIDAddress(id)
	if err != nil {
		t.Fatal(err)
	}
	return address
}

func NewSECP256K1Addr(t testing.TB, pubkey string) addr.Address {
	// the pubkey of a secp256k1 address is hashed for consistent length.
	address, err := addr.NewSecp256k1Address([]byte(pubkey))
	if err != nil {
		t.Fatal(err)
	}
	return address
}

func NewBLSAddr(t testing.TB, seed int64) addr.Address {
	buf := make([]byte, addr.BlsPublicKeyBytes)
	binary.PutVarint(buf, seed)

	address, err := addr.NewBLSAddress(buf)
	if err != nil {
		t.Fatal(err)
	}
	return address
}

func NewActorAddr(t testing.TB, data string) addr.Address {
	address, err := addr.NewActorAddress([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	return address
}

func IdFromAddress(a addr.Address) uint64 {
	if a.Protocol() != addr.ID {
		panic("must be ID protocol address")
	}
	id, _, err := varint.FromUvarint(a.Payload())
	if err != nil {
		panic(err)
	}
	return id
}

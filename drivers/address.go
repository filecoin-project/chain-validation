package drivers

import (
	"math/rand"

	addr "github.com/filecoin-project/go-address"
	"github.com/multiformats/go-varint"
)

// If you use this method while writing a test you are more than likely doing something wrong.
func NewIDAddr(h *ValidationHarness, id uint64) addr.Address {
	address, err := addr.NewIDAddress(id)
	if err != nil {
		h.Fatal(err)
	}
	return address
}

func NewSECP256K1Addr(h *ValidationHarness, pubkey string) addr.Address {
	// the pubkey of a secp256k1 address is hashed for consistent length.
	address, err := addr.NewSecp256k1Address([]byte(pubkey))
	if err != nil {
		h.Fatal(err)
	}
	return address
}

func NewBLSAddr(h *ValidationHarness, seed int64) addr.Address {
	// the pubkey of a bls address is not hashed and must be the correct length.
	buf := make([]byte, 48)
	r := rand.New(rand.NewSource(seed))
	r.Read(buf)

	address, err := addr.NewBLSAddress(buf)
	if err != nil {
		h.Fatal(err)
	}
	return address
}

func NewActorAddr(h *ValidationHarness, data string) addr.Address {
	address, err := addr.NewActorAddress([]byte(data))
	if err != nil {
		h.Fatal(err)
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

package state

import (
	"context"

	address "github.com/filecoin-project/go-address"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	crypto_spec "github.com/filecoin-project/specs-actors/actors/crypto"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
)

// Wrapper abstracts the inspection and mutation of an implementation-specific state tree and storage.
// The interface wraps a single, mutable state.
type Wrapper interface {
	// Returns the CID of the root node of the state tree.
	Cid() cid.Cid

	// Returns the actor state at `address` (or an error if there is none).
	Actor(address address.Address) (Actor, error)

	// Returns the actor storage for the actor at `address` (which is empty if there is no such actor).
	Storage() (Storage, error)

	// Creates a new private key and returns the associated address.
	NewAccountAddress() (address.Address, error)

	// Sign data with addr's key.
	Sign(ctx context.Context, addr address.Address, data []byte) (*crypto_spec.Signature, error)

	// Installs a new actor in the state tree.
	// This signature will probably become a little more complex when the actor state is non-empty.
	SetActor(addr address.Address, code cid.Cid, balance big_spec.Int) (Actor, Storage, error)

	// Installs a new singleton actor in the state tree.
	SetSingletonActor(addr address.Address, balance big_spec.Int) (Actor, Storage, error)
}

type Signer interface {
	Sign(ctx context.Context, addr address.Address, data []byte) (*crypto_spec.Signature, error)
}

// Actor is an abstraction over the actor states stored in the root of the state tree.
type Actor interface {
	Code() cid.Cid
	Head() cid.Cid
	CallSeqNum() int64
	Balance() big_spec.Int
}

type Storage interface {
	cbor.IpldStore
	Context() context.Context
}

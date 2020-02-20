package state

import (
	"context"

	address "github.com/filecoin-project/go-address"
	cid "github.com/ipfs/go-cid"

	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	crypto_spec "github.com/filecoin-project/specs-actors/actors/crypto"
	runtime_spec "github.com/filecoin-project/specs-actors/actors/runtime"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"
)

// Wrapper abstracts the inspection and mutation of an implementation-specific state tree and storage.
// The interface wraps a single, mutable state.
type Wrapper interface {
	// Returns the CID of the root node of the state tree.
	Root() cid.Cid

	// Returns the current VM storage
	Store() (adt_spec.Store, error)

	// Returns the actor state at `address` (or an error if there is none).
	Actor(address address.Address) (Actor, error)

	// Installs a new actor in the state tree. The new actors head cid is calculated by marshaling `state` and taking its
	// cid.
	SetActorState(addr address.Address, balance big_spec.Int, code cid.Cid, state runtime_spec.CBORMarshaler) (Actor, error)

	// Creates a new secp private key and returns the associated address.
	NewSecp256k1AccountAddress() address.Address

	// Creates a new BLS private key and returns the associated address.
	NewBLSAccountAddress() address.Address

	// Sign data with addr's key.
	Sign(ctx context.Context, addr address.Address, data []byte) (*crypto_spec.Signature, error)
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

package state

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

// Wrapper abstracts the inspection and mutation of an implementation-specific state tree and storage.
// The interface wraps a single, mutable state.
type Wrapper interface {
	// Returns the CID of the root node of the state tree.
	Cid() cid.Cid

	// Returns the actor state at `address` (or an error if there is none).
	Actor(address address.Address) (Actor, error)

	// Returns the actor storage for the actor at `address` (which is empty if there is no such actor).
	Storage(address address.Address) (Storage, error)

	// Creates a new private key and returns the associated address.
	NewAccountAddress() (address.Address, error)

	// Sign data with addr's key.
	Sign(ctx context.Context, addr address.Address, data []byte) (*types.Signature, error)

	// Installs a new actor in the state tree.
	// This signature will probably become a little more complex when the actor state is non-empty.
	SetActor(address address.Address, code actors.ActorCodeID, balance types.BigInt) (Actor, Storage, error)

	// Installs a new singleton actor in the state tree.
	SetSingletonActor(address actors.SingletonActorID, balance types.BigInt) (Actor, Storage, error)
}

type Signer interface {
	Sign(ctx context.Context, addr address.Address, data []byte) (*types.Signature, error)
}

// Actor is an abstraction over the actor states stored in the root of the state tree.
type Actor interface {
	Code() cid.Cid
	Head() cid.Cid
	Nonce() uint64
	Balance() types.BigInt
}

// Storage provides a key/value store for actor state.
type Storage interface {
	Get(c cid.Cid, out interface{}) error
}

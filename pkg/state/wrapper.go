package state

import (
	"github.com/ipfs/go-cid"
)

// Wrapper abstracts the inspection and mutation of an implementation-specific state tree and storage.
// The interface wraps a single, mutable state.
type Wrapper interface {
	// Returns the CID of the root node of the state tree.
	Cid() cid.Cid

	// Returns the actor state at `address` (or an error if there is none).
	Actor(address Address) (Actor, error)

	// Returns an abstraction over a payment channel actors state.
	PaymentChannelActorState(address Address) (PaymentChannelActorState, error)

	// Returns the actor storage for the actor at `address` (which is empty if there is no such actor).
	Storage(address Address) (Storage, error)

	// Creates a new private key and returns the associated address.
	NewAccountAddress() (Address, error)

	// Installs a new actor in the state tree.
	// This signature will probably become a little more complex when the actor state is non-empty.
	SetActor(address Address, code ActorCodeID, balance AttoFIL) (Actor, Storage, error)

	// Installs a new singleton actor in the state tree.
	SetSingletonActor(address SingletonActorID, balance AttoFIL) (Actor, Storage, error)
}

// Actor is an abstraction over the actor states stored in the root of the state tree.
type Actor interface {
	Code() cid.Cid
	Head() cid.Cid
	Nonce() uint64
	Balance() AttoFIL
}

// Storage provides a key/value store for actor state.
type Storage interface {
	Get(c cid.Cid, out interface{}) error
}

// PaymentChannelActorState is an abstraction over a payment channel actor's state.
type PaymentChannelActorState interface {
	From() Address
	To() Address
	ToSend() AttoFIL
	ClosingAt() uint64
	MinCloseHeight() uint64
}

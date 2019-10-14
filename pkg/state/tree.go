package state

import (
	cid "github.com/ipfs/go-cid"
)

// Tree is an abstraction over a concrete state tree implementation.
type Tree interface {
	// Returns the CID of the root node of the state tree.
	Cid() cid.Cid
	// Returns the actor state at `address` (or an error if there is none).
	Actor(address Address) (Actor, error)
	// Returns the actor storage for the actor at `address` (or an error if there is no such actor).
	ActorStorage(address Address) (Storage, error)
}

// Actor is an abstraction over the actor states stored in the root of the state tree.
type Actor interface {
	Code() cid.Cid
	Head() cid.Cid
	Nonce() uint64
	Balance() AttoFIL
}

type Storage interface {
	Get(cid cid.Cid, out interface{}) error
}
package state

import (
	"github.com/ipfs/go-cid"
)

// Factory abstracts over the construction of concrete implementation-specific state objects.
type Factory interface {
	NewAddress() (Address, error)
	NewActor(code cid.Cid, balance AttoFIL) Actor
	NewState(actors []ActorAndAddress) (Tree, StorageMap, error)
}

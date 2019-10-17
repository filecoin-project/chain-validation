package state

import "github.com/ipfs/go-cid"

// StorageMap is a collection of actor storages indexed by address.
type StorageMap interface {
	NewStorage(addr Address, actor Actor) (Storage, error)
}

// Storage provides a key/value store for actor state.
type Storage interface {
	Get(cid cid.Cid) ([]byte, error)
}

package storage

import (
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"

	vstate "github.com/filecoin-project/chain-validation/pkg/state"
)

type Factory interface {
	NewStorageMap(store blockstore.Blockstore) StorageMap
}

type StorageMap interface {
	NewStorage(addr vstate.Address, actor vstate.Actor) Storage
	Flush() error
}

type Storage interface {
	Get(cid cid.Cid) ([]byte, error)
}

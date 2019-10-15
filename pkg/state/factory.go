package state

import (
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

// StateFactory abstracts over concrete state manipulation methods.
type StateFactory interface {
	NewActor(code cid.Cid, balance AttoFIL) Actor
	NewState(actors map[Address]Actor) (Tree, error)
	ApplyMessage(*VMContext, Tree, interface{}) (Tree, error)
}

type Validator struct {
	factory StateFactory
	store   StorageMap
}


func NewValidator(factory StateFactory, storage StorageFactory) *Validator {
	bs := blockstore.NewBlockstore(datastore.NewMapDatastore())
	storageMap := storage.NewStorageMap(bs)

	return &Validator{factory, storageMap}
}

// TODO probs make an interface
type VMContext struct {
	Store StorageMap
}

func (v *Validator) ApplyMessages(tree Tree, messages []interface{}) (Tree, error) {
	vmctx := &VMContext{v.store}
	var err error
	for _, m := range messages {
		tree, err = v.factory.ApplyMessage(vmctx,tree, m)
		if err != nil {
			return nil, err
		}
	}
	return tree, nil
}
type StorageFactory interface {
	NewStorageMap(store blockstore.Blockstore) StorageMap
}

type StorageMap interface {
	NewStorage(addr Address, actor Actor) Storage
	Flush() error
}

type Storage interface {
	Get(cid cid.Cid) ([]byte, error)
}


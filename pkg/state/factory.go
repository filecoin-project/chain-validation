package state

import (
	"github.com/filecoin-project/chain-validation/pkg/storage"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

// Factory abstracts over concrete state manipulation methods.
type Factory interface {
	NewActor(code cid.Cid, balance AttoFIL) Actor
	NewState(actors map[Address]Actor) (Tree, error)
	ApplyMessage(*VMContext, Tree, interface{}) (Tree, error)
}

type Validator struct {
	factory Factory
	store storage.StorageMap
}


func NewValidator(factory Factory, storage storage.Factory) *Validator {
	bs := blockstore.NewBlockstore(datastore.NewMapDatastore())
	storageMap := storage.NewStorageMap(bs)

	return &Validator{factory, storageMap}
}

// TODO probs make an interface
type VMContext struct {
	Store storage.StorageMap
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


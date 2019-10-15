package state

import (
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

// Factory abstracts over concrete state manipulation methods.
type Factory interface {
	NewActor(code cid.Cid, balance AttoFIL) Actor
	NewState(actors map[Address]Actor) (Tree, error)
	ApplyMessage(VMParams, Tree, interface{}) (Tree, error)
}

type VMParams interface {
	StorageMap() Storage
	BlockHeight() uint64
}

type ValidatorVMParams struct {
	storageMap blockstore.Blockstore
	blockHeight uint64
}

func (vc *ValidatorVMParams) StorageMap() Storage {
	panic("NYI")
}

func (vc *ValidatorVMParams) BlockHeight() uint64 {
	panic("NYI")
}

func NewValidatorContext() VMParams {
	return &ValidatorVMParams{}
}

type Validator struct {
	factory Factory
	context VMParams
}

func NewValidator(factory Factory) *Validator {
	return &Validator{factory, NewValidatorContext()}
}

func (v *Validator) ApplyMessages(tree Tree, messages []interface{}) (Tree, error) {
	var err error
	for _, m := range messages {
		tree, err = v.factory.ApplyMessage(v.context, tree, m)
		if err != nil {
			return nil, err
		}
	}
	return tree, nil
}

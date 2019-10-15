package state

import (
	"github.com/ipfs/go-cid"
)

// Factory abstracts over concrete state manipulation methods.
type Factory interface {
	NewActor(code cid.Cid, balance AttoFIL) Actor
	NewState(actors map[Address]Actor) (Tree, error)
	ApplyMessage(Tree, interface{}) (Tree, error)
}

type Validator struct {
	factory Factory
}

func NewValidator(factory Factory) *Validator {
	return &Validator{factory}
}

func (v *Validator) ApplyMessages(tree Tree, messages []interface{}) (Tree, error) {
	var err error
	for _, m := range messages {
		tree, err = v.factory.ApplyMessage(tree, m)
		if err != nil {
			return nil, err
		}
	}
	return tree, nil
}

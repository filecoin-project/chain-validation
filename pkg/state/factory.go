package state

import (
	"github.com/ipfs/go-cid"
)

// Factory abstracts over concrete state manipulation methods.
type Factory interface {
	NewActor(code cid.Cid, balance AttoFIL) Actor
	NewState(actors map[Address]Actor) Tree
	ApplyMessage(Tree, interface{}) Tree
}

type Validator struct {
	factory Factory
}

func NewValidator(factory Factory) *Validator {
	return &Validator{factory}
}

func (v *Validator) ApplyMessages(tree Tree, messages []interface{}) (Tree, error) {
	for _, m := range messages {
		tree = v.factory.ApplyMessage(tree, m)
	}
	return tree, nil
}

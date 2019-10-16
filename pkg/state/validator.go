package state

import (
	"github.com/ipfs/go-cid"
)

// StateFactory abstracts over concrete state manipulation methods.
type StateFactory interface {
	NewActor(code cid.Cid, balance AttoFIL) Actor
	NewState(actors map[Address]Actor) (Tree, StorageMap, error)
	ApplyMessage(state Tree, storage StorageMap, context *ExecutionContext, msg interface{}) (Tree, error)
}

// ExecutionContext provides the context for execution of a message.
type ExecutionContext struct {
	Epoch      uint64  // The epoch number ("height") during which a message is executed.
	MinerOwner Address // The miner actor which earns gas fees from message execution.
}

func NewExecutionContext(epoch uint64, miner Address) *ExecutionContext {
	return &ExecutionContext{epoch, miner}
}

type Validator struct {
	factory StateFactory
}

func NewValidator(factory StateFactory) *Validator {
	return &Validator{factory}
}

func (v *Validator) ApplyMessages(context *ExecutionContext, tree Tree, storage StorageMap, messages []interface{}) (Tree, error) {
	var err error
	for _, m := range messages {
		tree, err = v.factory.ApplyMessage(tree, storage, context, m)
		if err != nil {
			return nil, err
		}
	}
	return tree, nil
}

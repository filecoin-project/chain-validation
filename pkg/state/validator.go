package state

import (
	"github.com/ipfs/go-cid"
)

// StateFactory abstracts over concrete state manipulation methods.
type StateFactory interface {
	NewAddress() (Address, error)
	NewActor(code cid.Cid, balance AttoFIL) Actor
	NewState(actors map[Address]Actor) (Tree, StorageMap, error)
	ApplyMessage(state Tree, storage StorageMap, context *ExecutionContext, msg interface{}) (Tree, MessageReceiept, error)
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

func (v *Validator) ApplyMessages(context *ExecutionContext, tree Tree, storage StorageMap, messages []interface{}) (Tree, []MessageReceiept, error) {
	var err error
	var receipts []MessageReceiept
	for _, m := range messages {
		var mr MessageReceiept
		tree, mr, err = v.factory.ApplyMessage(tree, storage, context, m)
		receipts = append(receipts, mr)
		if err != nil {
			return nil, nil, err
		}
	}
	return tree, receipts, nil
}

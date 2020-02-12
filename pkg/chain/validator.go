package chain

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"

	"github.com/filecoin-project/chain-validation/pkg/chain/types"
	"github.com/filecoin-project/chain-validation/pkg/state"
)

// Applier applies abstract messages to states.
type Applier interface {
	ApplyMessage(context *ExecutionContext, state state.Wrapper, msg *types.Message) (types.MessageReceipt, error)
}

// ExecutionContext provides the context for execution of a message.
type ExecutionContext struct {
	Epoch      abi.ChainEpoch  // The epoch number ("height") during which a message is executed.
	MinerOwner address.Address // The miner actor which earns gas fees from message execution.
}

// NewExecutionContext builds a new execution context.
func NewExecutionContext(epoch int64, miner address.Address) *ExecutionContext {
	return &ExecutionContext{abi.ChainEpoch(epoch), miner}
}

// Validator arranges the execution of a sequence of messages, returning the resulting receipts and state.
type Validator struct {
	applier Applier
}

// NewValidator builds a new validator.
func NewValidator(executor Applier) *Validator {
	return &Validator{executor}
}

// ApplyMessages applies a message to a state
func (v *Validator) ApplyMessage(context *ExecutionContext, state state.Wrapper, message *types.Message) (types.MessageReceipt, error) {
	return v.applier.ApplyMessage(context, state, message)
}

package chain

import (
	"github.com/filecoin-project/chain-validation/pkg/state"
)

// Applier applies abstract messages to states.
type Applier interface {
	ApplyMessage(context *ExecutionContext, state state.Wrapper, msg interface{}) (MessageReceipt, error)
}

// MessageReceipt is the return value of message application.
type MessageReceipt struct {
	ExitCode    uint8
	ReturnValue []byte
	GasUsed     state.GasUnit
}

// ExecutionContext provides the context for execution of a message.
type ExecutionContext struct {
	Epoch      uint64        // The epoch number ("height") during which a message is executed.
	MinerOwner state.Address // The miner actor which earns gas fees from message execution.
}

// NewExecutionContext builds a new execution context.
func NewExecutionContext(epoch uint64, miner state.Address) *ExecutionContext {
	return &ExecutionContext{epoch, miner}
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
func (v *Validator) ApplyMessage(context *ExecutionContext, state state.Wrapper, message interface{}) (MessageReceipt, error) {
	return v.applier.ApplyMessage(context, state, message)
}

package chain

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi/big"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

// Applier applies abstract messages to states.
type Applier interface {
	ApplyMessage(context *ExecutionContext, state state.Wrapper, msg *Message) (MessageReceipt, error)
}

// MessageReceipt is the return value of message application.
type MessageReceipt struct {
	ExitCode    uint8
	ReturnValue []byte
	GasUsed     big.Int
}

// ExecutionContext provides the context for execution of a message.
type ExecutionContext struct {
	Epoch      int64           // The epoch number ("height") during which a message is executed.
	MinerOwner address.Address // The miner actor which earns gas fees from message execution.
}

// NewExecutionContext builds a new execution context.
func NewExecutionContext(epoch int64, miner address.Address) *ExecutionContext {
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
func (v *Validator) ApplyMessage(context *ExecutionContext, state state.Wrapper, message *Message) (MessageReceipt, error) {
	return v.applier.ApplyMessage(context, state, message)
}

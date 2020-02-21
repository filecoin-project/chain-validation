package chain

import (
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/state"
)

// Validator arranges the execution of a sequence of messages, returning the resulting receipts and state.
type Validator struct {
	applier state.Applier
}

// NewValidator builds a new validator.
func NewValidator(executor state.Applier) *Validator {
	return &Validator{executor}
}

// ApplyMessages applies a message to a state
func (v *Validator) ApplyMessage(context *types.ExecutionContext, state state.VMWrapper, message *types.Message) (types.MessageReceipt, error) {
	return v.applier.ApplyMessage(context, state, message)
}

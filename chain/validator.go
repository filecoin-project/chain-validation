package chain

import (
	"github.com/filecoin-project/specs-actors/actors/abi"

	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/state"
)

// validator arranges the execution of a sequence of messages, returning the resulting receipts and state.
type Validator struct {
	applier state.Applier
}

// NewValidator builds a new validator.
func NewValidator(executor state.Applier) *Validator {
	return &Validator{executor}
}

// ApplyMessages applies a message to a state
func (v *Validator) ApplyMessage(epoch abi.ChainEpoch, message *types.Message) (types.ApplyMessageResult, error) {
	return v.applier.ApplyMessage(epoch, message)
}

func (v *Validator) ApplySignedMessage(epoch abi.ChainEpoch, message *types.SignedMessage) (types.ApplyMessageResult, error) {
	return v.applier.ApplySignedMessage(epoch, message)
}

func (v *Validator) ApplyTipSetMessages(epoch abi.ChainEpoch, blocks []types.BlockMessagesInfo, rnd state.RandomnessSource) (types.ApplyTipSetResult, error) {
	return v.applier.ApplyTipSetMessages(epoch, blocks, rnd)
}

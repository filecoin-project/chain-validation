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

type ApplyMessageResult struct {
	Receipt types.MessageReceipt
	Penalty abi.TokenAmount
	Reward  abi.TokenAmount
	Root    string
}

type ApplyTipSetMessagesResult struct {
	Receipts []types.MessageReceipt
	Root     string
}

// NewValidator builds a new validator.
func NewValidator(executor state.Applier) *Validator {
	return &Validator{executor}
}

// ApplyMessages applies a message to a state
func (v *Validator) ApplyMessage(context *types.ExecutionContext, state state.VMWrapper, message *types.Message) (ApplyMessageResult, error) {
	receipt, penalty, reward, err := v.applier.ApplyMessage(context, state, message)
	return ApplyMessageResult{receipt, penalty, reward, state.Root().String()}, err
}

func (v *Validator) ApplySignedMessage(context *types.ExecutionContext, state state.VMWrapper, message *types.SignedMessage) (ApplyMessageResult, error) {
	receipt, penalty, reward, err := v.applier.ApplySignedMessage(context, state, message)
	return ApplyMessageResult{receipt, penalty, reward, state.Root().String()}, err
}

func (v *Validator) ApplyTipSetMessages(epoch abi.ChainEpoch, state state.VMWrapper, blocks []types.BlockMessagesInfo, rnd state.RandomnessSource) (ApplyTipSetMessagesResult, error) {
	receipts, err := v.applier.ApplyTipSetMessages(state, blocks, epoch, rnd)
	return ApplyTipSetMessagesResult{receipts, state.Root().String()}, err
}

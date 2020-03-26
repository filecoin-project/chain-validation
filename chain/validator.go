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

type ApplyResult struct {
	Receipt types.MessageReceipt
	Penalty abi.TokenAmount
	Reward  abi.TokenAmount
}

// NewValidator builds a new validator.
func NewValidator(executor state.Applier) *Validator {
	return &Validator{executor}
}

// ApplyMessages applies a message to a state
func (v *Validator) ApplyMessage(context *types.ExecutionContext, state state.VMWrapper, message *types.Message) (ApplyResult, error) {
	receipt, penalty, reward, err := v.applier.ApplyMessage(context, state, message)
	return ApplyResult{
		Receipt: receipt,
		Penalty: penalty,
		Reward:  reward,
	}, err
}

func (v *Validator) ApplySignedMessage(context *types.ExecutionContext, state state.VMWrapper, message *types.SignedMessage) (ApplyResult, error) {
	receipt, penalty, reward, err := v.applier.ApplySignedMessage(context, state, message)
	return ApplyResult{
		Receipt: receipt,
		Penalty: penalty,
		Reward:  reward,
	}, err
}

func (v *Validator) ApplyTipSetMessages(epoch abi.ChainEpoch, state state.VMWrapper, blocks []types.BlockMessagesInfo, rnd state.RandomnessSource) ([]types.MessageReceipt, error) {
	return v.applier.ApplyTipSetMessages(state, blocks, epoch, rnd)
}

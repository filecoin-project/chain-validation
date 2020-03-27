package chain

import (
	"fmt"

	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/ipfs/go-cid"

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
func (v *Validator) ApplyMessage(context *types.ExecutionContext, state state.VMWrapper, message *types.Message) (ApplyMessageResult, error) {
	receipt, penalty, reward, err := v.applier.ApplyMessage(context, state, message)
	return ApplyMessageResult{receipt, penalty, reward, state.Root().String()}, err
}

func (v *Validator) ApplySignedMessage(context *types.ExecutionContext, state state.VMWrapper, message *types.SignedMessage) (ApplyMessageResult, error) {
	receipt, penalty, reward, err := v.applier.ApplySignedMessage(context, state, message)
	return ApplyMessageResult{receipt, penalty, reward, state.Root().String()}, err
}

func (v *Validator) ApplyTipSetMessages(epoch abi.ChainEpoch, state state.VMWrapper, blocks []types.BlockMessagesInfo, rnd state.RandomnessSource) (ApplyTipSetResult, error) {
	receipts, err := v.applier.ApplyTipSetMessages(state, blocks, epoch, rnd)
	return ApplyTipSetResult{receipts, state.Root().String()}, err
}

type Trackable interface {
	GoSyntax() string
	GoContainer() string
}

var _ Trackable = (*ApplyMessageResult)(nil)
var _ Trackable = (*ApplyTipSetResult)(nil)

type ApplyMessageResult struct {
	Receipt types.MessageReceipt
	Penalty abi.TokenAmount
	Reward  abi.TokenAmount
	Root    string
}

func (mr ApplyMessageResult) GoSyntax() string {
	return fmt.Sprintf("chain.ApplyMessageResult{Receipt: %#v, Penalty: abi.NewTokenAmount(%d), Reward: abi.NewTokenAmount(%d), Root: \"%s\"}", mr.Receipt, mr.Penalty, mr.Reward, mr.Root)
}

func (mr ApplyMessageResult) GoContainer() string {
	return "[]chain.ApplyMessageResult"
}

func (mr ApplyMessageResult) StateRoot() cid.Cid {
	root, err := cid.Decode(mr.Root)
	if err != nil {
		panic(err)
	}
	return root
}

func (mr ApplyMessageResult) GasUsed() types.GasUnits {
	return mr.Receipt.GasUsed
}

type ApplyTipSetResult struct {
	Receipts []types.MessageReceipt
	Root     string
}

func (tr ApplyTipSetResult) GoSyntax() string {
	return fmt.Sprintf("%#v", tr)
}

func (tr ApplyTipSetResult) GoContainer() string {
	return "[]chain.ApplyTipSetResult"
}

func (mr ApplyTipSetResult) StateRoot() cid.Cid {
	root, err := cid.Decode(mr.Root)
	if err != nil {
		panic(err)
	}
	return root
}

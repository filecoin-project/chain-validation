package state

import (
	"context"

	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/crypto"

	"github.com/filecoin-project/chain-validation/chain/types"
)

// Applier applies abstract messages to states.
type Applier interface {
	ApplyMessage(context *types.ExecutionContext, state VMWrapper, msg *types.Message) (types.MessageReceipt, abi.TokenAmount, abi.TokenAmount, error)
	ApplySignedMessage(context *types.ExecutionContext, state VMWrapper, msg *types.SignedMessage) (types.MessageReceipt, abi.TokenAmount, abi.TokenAmount, error)
	ApplyTipSetMessages(state VMWrapper, blocks []types.BlockMessagesInfo, epoch abi.ChainEpoch, rnd RandomnessSource) ([]types.MessageReceipt, error)
}

// RandomnessSource provides randomness to actors.
type RandomnessSource interface {
	Randomness(ctx context.Context, tag crypto.DomainSeparationTag, epoch abi.ChainEpoch, entropy []byte) (abi.Randomness, error)
}

// Specifies a domain for randomness generation.
type RandomnessType int

package state

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/actors/crypto"

	"github.com/filecoin-project/chain-validation/chain/types"
)

// Applier applies abstract messages to states.
type Applier interface {
	ApplyMessage(epoch abi.ChainEpoch, msg *types.Message) (types.ApplyMessageResult, error)
	ApplySignedMessage(epoch abi.ChainEpoch, msg *types.SignedMessage) (types.ApplyMessageResult, error)
	ApplyTipSetMessages(epoch abi.ChainEpoch, blocks []types.BlockMessagesInfo, rnd RandomnessSource) (types.ApplyTipSetResult, error)
}

// RandomnessSource provides randomness to actors.
type RandomnessSource interface {
	Randomness(ctx context.Context, tag crypto.DomainSeparationTag, epoch abi.ChainEpoch, entropy []byte) (abi.Randomness, error)
}

// Specifies a domain for randomness generation.
type RandomnessType int

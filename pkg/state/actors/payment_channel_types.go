package actors

import (
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type PaymentChannelActorStateType struct {
	From address.Address
	To   address.Address

	ToSend types.BigInt

	ClosingAt      uint64
	MinCloseHeight uint64

	// TODO: needs to be map[uint64]*laneState
	// waiting on refmt#35 to be fixed
	LaneStates map[string]*LaneState
}

type LaneState struct {
	Closed   bool
	Redeemed types.BigInt
	Nonce    uint64
}

package paych

import (
	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type PaymentInfo struct {
	PayChActor     address.Address
	Payer          address.Address
	ChannelMessage *cid.Cid

	Vouchers []*types.SignedVoucher
}

type LaneState struct {
	Closed   bool
	Redeemed types.BigInt
	Nonce    uint64
}

type PaymentChannelActorState struct {
	From address.Address
	To   address.Address

	ToSend types.BigInt

	ClosingAt      uint64
	MinCloseHeight uint64

	// TODO: needs to be map[uint64]*laneState
	// waiting on refmt#35 to be fixed
	LaneStates map[string]*LaneState
}

type PaymentChannelConstructorParams struct {
	To address.Address
}

type PaymentChannelUpdateParams struct {
	Sv     types.SignedVoucher
	Secret []byte
	Proof  []byte
}

type PaymentVerifyParams struct {
	Extra []byte
	Proof []byte
}

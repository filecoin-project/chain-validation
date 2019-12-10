package strgmrkt

import (
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
	"github.com/ipfs/go-cid"
)

type SerializationMode = uint64

const (
	SerializationUnixFSv0 = iota
)

type StorageDeal struct {
	Proposal         StorageDealProposal
	CounterSignature *types.Signature
}

type OnChainDeal struct {
	Deal            StorageDeal
	ActivationEpoch uint64 // 0 = inactive
}

type StorageMarketState struct {
	Balances cid.Cid
	Deals    cid.Cid

	NextDealID uint64
}

type StorageDealProposal struct {
	PieceRef           []byte // cid bytes
	PieceSize          uint64
	PieceSerialization SerializationMode

	Client   address.Address
	Provider address.Address

	ProposalExpiration uint64
	Duration           uint64

	StoragePricePerEpoch types.BigInt
	StorageCollateral    types.BigInt

	ProposerSignature *types.Signature
}

type StorageParticipantBalance struct {
	Locked    types.BigInt
	Available types.BigInt
}

type PublishStorageDealResponse struct {
	DealIDs []uint64
}

//
// Message Method  Params
//

type WithdrawBalanceParams struct {
	Balance types.BigInt
}

type PublishStorageDealsParams struct {
	Deals []StorageDeal
}

type ActivateStorageDealsParams struct {
	Deals []uint64
}

type ComputeDataCommitmentParams struct {
	DealIDs    []uint64
	SectorSize uint64
}

type ProcessStorageDealsPaymentParams struct {
	DealIDs []uint64
}

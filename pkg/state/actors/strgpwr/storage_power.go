package strgpwr

import (
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type StoragePowerState struct {
	Miners         cid.Cid
	ProvingBuckets cid.Cid // amt[ProvingPeriodBucket]hamt[minerAddress]struct{}
	MinerCount     uint64
	LastMinerCheck uint64

	TotalStorage types.BigInt
}

type CreateStorageMinerParams struct {
	Owner      address.Address
	Worker     address.Address
	SectorSize uint64
	PeerID     peer.ID
}

type UpdateStorageParams struct {
	Delta                    types.BigInt
	NextProvingPeriodEnd     uint64
	PreviousProvingPeriodEnd uint64
}

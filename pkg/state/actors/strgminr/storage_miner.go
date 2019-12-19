package strgminr

import (
	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type StorageMinerActorState struct {
	// PreCommittedSectors is the set of sectors that have been committed to but not
	// yet had their proofs submitted
	PreCommittedSectors map[string]*PreCommittedSector

	// All sectors this miner has committed.
	Sectors cid.Cid

	// TODO: Spec says 'StagedCommittedSectors', which one is it?

	// Sectors this miner is currently mining. It is only updated
	// when a PoSt is submitted (not as each new sector commitment is added).
	ProvingSet cid.Cid

	// TODO: these:
	//    SectorTable
	//    SectorExpirationQueue
	//    ChallengeStatus

	// Contains mostly static info about this miner
	Info cid.Cid

	// Faulty sectors reported since last SubmitPost,
	// up to the current proving period's challenge time.
	CurrentFaultSet types.BitField

	// Faults submitted after the current proving period's challenge time,
	// but before the PoSt for that period is submitted.
	// These become the currentFaultSet when a PoSt is submitted.
	NextFaultSet types.BitField

	// Sectors reported during the last PoSt submission as being 'done'.
	// The collateral for them is still being held until
	// the next PoSt submission in case early sector
	// removal penalization is needed.
	NextDoneSet types.BitField

	// Amount of power this miner has.
	Power types.BigInt

	// Active is set to true after the miner has submitted their first PoSt
	Active bool

	// The height at which this miner was slashed at.
	SlashedAt uint64

	ProvingPeriodEnd uint64
}

type MinerInfo struct {
	// Account that owns this miner.
	// - Income and returned collateral are paid to this address.
	// - This address is also allowed to change the worker address for the miner.
	Owner address.Address

	// Worker account for this miner.
	// This will be the key that is used to sign blocks created by this miner, and
	// sign messages sent on behalf of this miner to commit sectors, submit PoSts, and
	// other day to day miner activities.
	Worker address.Address

	// Libp2p identity that should be used when connecting to this miner.
	PeerID peer.ID

	// Amount of space in each sector committed to the network by this miner.
	SectorSize uint64

	// SubsectorCount
}

type PreCommittedSector struct {
	Info          SectorPreCommitInfo
	ReceivedEpoch uint64
}

type SectorPreCommitInfo struct {
	SectorNumber uint64

	CommR     []byte // TODO: Spec says CID
	SealEpoch uint64
	DealIDs   []uint64
}

type UpdatePeerIDParams struct {
	PeerID peer.ID
}

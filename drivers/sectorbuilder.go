package drivers

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

type MockSectorBuilder struct {
	// PreSeal is intexted by sectorID
	MinerSectors map[address.Address][]*types.PreSeal

	cidGetter func() cid.Cid
}

func NewMockSectorBuilder() *MockSectorBuilder {
	return &MockSectorBuilder{
		MinerSectors: make(map[address.Address][]*types.PreSeal),
		cidGetter:    utils.NewProofCidForTestGetter(),
	}
}

func (msb *MockSectorBuilder) NewPreSealedSector(miner, client address.Address, pt abi.RegisteredProof, ssize abi.SectorSize, start, end abi.ChainEpoch) *types.PreSeal {
	minerSectors := msb.MinerSectors[miner]
	sectorID := len(minerSectors)

	R := msb.cidGetter()
	D := msb.cidGetter()
	preseal := &types.PreSeal{
		CommR:    R,
		CommD:    D,
		SectorID: abi.SectorNumber(sectorID),
		Deal: market.DealProposal{
			PieceCID:   D,
			PieceSize:  abi.PaddedPieceSize(ssize),
			Client:     client,
			Provider:   miner,
			StartEpoch: start,
			EndEpoch:   end,
			// TODO how do we want to interact with these values?
			StoragePricePerEpoch: big.Zero(),
			ProviderCollateral:   big.Zero(),
			ClientCollateral:     big.Zero(),
		},
		ProofType: pt,
	}

	msb.MinerSectors[miner] = append(msb.MinerSectors[miner], preseal)
	return preseal
}

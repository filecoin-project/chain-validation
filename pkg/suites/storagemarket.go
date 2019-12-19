package suites

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgmrkt"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgpwr"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

func StorageMarketActorConstructor(t testing.TB, factory Factories) {
	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})
	mustCreateStorageMarketActor(c)
}

func StorageMarketBalanceUpdates(t testing.TB, factory Factories) {
	const initialBal = 2000000000
	const balAddAmount = 100
	const balWithdrawAmount = 10

	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})
	smaddr := mustCreateStorageMarketActor(c)

	alice := c.Driver().NewAccountActor(initialBal)
	bob := c.Driver().NewAccountActor(initialBal)

	mustAddBalance(c, alice, 0, balAddAmount)
	c.Driver().AssertStorageMarketParticipantAvailableBalance(smaddr, alice, types.NewInt(balAddAmount))

	mustAddBalance(c, bob, 0, balAddAmount)
	c.Driver().AssertStorageMarketParticipantAvailableBalance(smaddr, bob, types.NewInt(balAddAmount))

	mustWithdrawBalance(c, bob, 1, balWithdrawAmount)
	c.Driver().AssertStorageMarketParticipantAvailableBalance(smaddr, bob, types.NewInt(balAddAmount-balWithdrawAmount))
}

func StorageMarketStoragePublishDeal(t testing.TB, factory Factories) {
	const initialBal = "2000000000000000000000000"
	const dealID = 0
	const dealCost = 10
	const dealDuration = 10
	const dealExpiration = 20
	const pricePerEpoch = 1
	const collateral = 5
	ctx := context.Background()

	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})

	// create the storage market actor
	smaddr := mustCreateStorageMarketActor(c)

	// create an account to own a miner.
	minerOwner := c.Driver().NewAccountActorBigBalance(types.NewIntFromString(initialBal))

	mustAddBalance(c, minerOwner, 0, dealCost)
	c.Driver().AssertStorageMarketParticipantAvailableBalance(smaddr, minerOwner, types.NewInt(dealCost))
	c.Driver().AssertStorageMarketParticipantLockedBalance(smaddr, minerOwner, types.NewInt(0))

	peerID0 := RequireIntPeerID(t, 0)
	minerAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)
	mustCreateStorageMiner(c, 1, strgpwr.SectorSizes[0], types.NewIntFromString("1999999995415053581179420"), minerAddr, minerOwner, minerOwner, minerOwner, peerID0)

	// create a client and sign the deal
	client := c.Driver().NewAccountActorBigBalance(types.NewIntFromString(initialBal))
	mustAddBalance(c, client, 0, dealCost)
	c.Driver().AssertStorageMarketParticipantAvailableBalance(smaddr, client, types.NewInt(dealCost))
	c.Driver().AssertStorageMarketParticipantLockedBalance(smaddr, client, types.NewInt(0))

	dealProposal := strgmrkt.StorageDealProposal{
		PieceRef:             []byte{1},
		PieceSize:            1,
		PieceSerialization:   strgmrkt.SerializationUnixFSv0,
		Client:               client,
		Provider:             minerAddr,
		ProposalExpiration:   dealExpiration,
		Duration:             dealDuration,
		StoragePricePerEpoch: types.NewInt(pricePerEpoch),
		StorageCollateral:    types.NewInt(collateral),
	}
	err = dealProposal.Sign(ctx, client, c.Driver().State())
	require.NoError(t, err)

	// miner signs the deal
	storageDeal := strgmrkt.StorageDeal{
		Proposal: dealProposal,
	}
	err = storageDeal.Sign(ctx, minerOwner, c.Driver().State())
	require.NoError(t, err)

	mustPublishStorageDeal(c, 1, client, dealID, storageDeal)

	c.Driver().AssertStorageMarketHasOnChainDeal(smaddr, dealID, strgmrkt.OnChainDeal{
		Deal:            storageDeal,
		ActivationEpoch: 0,
	})
	c.Driver().AssertStorageMarketParticipantAvailableBalance(smaddr, client, types.NewInt(dealCost-dealProposal.TotalStoragePrice().Uint64()))
	c.Driver().AssertStorageMarketParticipantAvailableBalance(smaddr, minerOwner, types.NewInt(dealCost-collateral))

	c.Driver().AssertStorageMarketParticipantLockedBalance(smaddr, client, types.NewInt(dealCost))
	c.Driver().AssertStorageMarketParticipantLockedBalance(smaddr, minerOwner, types.NewInt(collateral))
}

func mustPublishStorageDeal(c Candy, nonce uint64, from address.Address, dealID uint64, storageDeal strgmrkt.StorageDeal) {
	// expected response
	pubDealResp := strgmrkt.PublishStorageDealResponse{
		DealIDs: []uint64{dealID},
	}
	respBytes, err := types.Serialize(&pubDealResp)
	require.NoError(c.TB(), err)

	msg, err := c.Producer().StorageMarketPublishStorageDeals(from, nonce, []strgmrkt.StorageDeal{storageDeal}, chain.Value(0))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: respBytes,
		GasUsed:     types.NewInt(0),
	})

}

func mustWithdrawBalance(c Candy, from address.Address, nonce, amount uint64) {
	msg, err := c.Producer().StorageMarketWithdrawBalance(from, nonce, types.NewInt(amount), chain.Value(0))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)

	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

func mustAddBalance(c Candy, from address.Address, nonce, amount uint64) {
	msg, err := c.Producer().StorageMarketAddBalance(from, nonce, chain.Value(amount))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)

	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

func mustCreateStorageMarketActor(c Candy) address.Address {
	_, _, err := c.Driver().State().SetSingletonActor(actors.StorageMarketAddress, types.NewInt(0))
	require.NoError(c.TB(), err)

	mt := strgmrkt.NewMarketTracker(c.TB())
	smaddr := c.Producer().SingletonAddress(actors.StorageMarketAddress)
	c.Driver().AssertStorageMarketState(smaddr, strgmrkt.StorageMarketState{
		Balances:   mt.Balance,
		Deals:      mt.Deals,
		NextDealID: 0,
	})
	return smaddr
}

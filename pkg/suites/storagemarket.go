package suites

// TODO uncomment when ready to implement
/*
import (
	"context"
	"github.com/filecoin-project/chain-validation/pkg/state"
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
	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})
	mustCreateStorageMarketActor(td)
}

func StorageMarketBalanceUpdates(t testing.TB, factory Factories) {
	const initialBal = 2000000000
	const balAddAmount = 100
	const balWithdrawAmount = 10

	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})
	smaddr := mustCreateStorageMarketActor(td)

	alice := td.State.NewAccountActor(initialBal)
	bob := td.State.NewAccountActor(initialBal)

	mustAddBalance(td, alice, 0, balAddAmount)
	td.State.AssertStorageMarketParticipantAvailableBalance(smaddr, alice, types.NewInt(balAddAmount))

	mustAddBalance(td, bob, 0, balAddAmount)
	td.State.AssertStorageMarketParticipantAvailableBalance(smaddr, bob, types.NewInt(balAddAmount))

	mustWithdrawBalance(td, bob, 1, balWithdrawAmount)
	td.State.AssertStorageMarketParticipantAvailableBalance(smaddr, bob, types.NewInt(balAddAmount-balWithdrawAmount))
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

	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})

	// create the storage market actor
	smaddr := mustCreateStorageMarketActor(td)

	// create an account to own a miner.
	minerOwner := td.State.NewAccountActorBigBalance(types.NewIntFromString(initialBal))

	mustAddBalance(td, minerOwner, 0, dealCost)
	td.State.AssertStorageMarketParticipantAvailableBalance(smaddr, minerOwner, types.NewInt(dealCost))
	td.State.AssertStorageMarketParticipantLockedBalance(smaddr, minerOwner, types.NewInt(0))

	peerID0 := RequireIntPeerID(t, 0)
	minerAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)
	mustCreateStorageMiner(td, 1, strgpwr.SectorSizes[0], types.NewIntFromString("1999999995415053581179420"), minerAddr, minerOwner, minerOwner, minerOwner, peerID0)

	// create a client and sign the deal
	client := td.State.NewAccountActorBigBalance(types.NewIntFromString(initialBal))
	mustAddBalance(td, client, 0, dealCost)
	td.State.AssertStorageMarketParticipantAvailableBalance(smaddr, client, types.NewInt(dealCost))
	td.State.AssertStorageMarketParticipantLockedBalance(smaddr, client, types.NewInt(0))

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
	err = dealProposal.Sign(ctx, client, td.State.State())
	require.NoError(t, err)

	// miner signs the deal
	storageDeal := strgmrkt.StorageDeal{
		Proposal: dealProposal,
	}
	err = storageDeal.Sign(ctx, minerOwner, td.State.State())
	require.NoError(t, err)

	mustPublishStorageDeal(td, 1, client, dealID, storageDeal)

	td.State.AssertStorageMarketHasOnChainDeal(smaddr, dealID, strgmrkt.OnChainDeal{
		Deal:            storageDeal,
		ActivationEpoch: 0,
	})
	td.State.AssertStorageMarketParticipantAvailableBalance(smaddr, client, types.NewInt(dealCost-dealProposal.TotalStoragePrice().Uint64()))
	td.State.AssertStorageMarketParticipantAvailableBalance(smaddr, minerOwner, types.NewInt(dealCost-collateral))

	td.State.AssertStorageMarketParticipantLockedBalance(smaddr, client, types.NewInt(dealCost))
	td.State.AssertStorageMarketParticipantLockedBalance(smaddr, minerOwner, types.NewInt(collateral))
}

func mustPublishStorageDeal(td TestDriver, nonce uint64, from address.Address, dealID uint64, storageDeal strgmrkt.StorageDeal) {
	// expected response
	pubDealResp := strgmrkt.PublishStorageDealResponse{
		DealIDs: []uint64{dealID},
	}
	respBytes, err := state.Serialize(&pubDealResp)
	require.NoError(td.TB(), err)

	msg, err := td.Producer.StorageMarketPublishStorageDeals(from, nonce, []strgmrkt.StorageDeal{storageDeal}, chain.Value(0))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State().State(), msg)
	require.NoError(td.TB(), err)
	td.State().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: respBytes,
		GasUsed:     types.NewInt(0),
	})

}

func mustWithdrawBalance(td TestDriver, from address.Address, nonce, amount uint64) {
	msg, err := td.Producer.StorageMarketWithdrawBalance(from, nonce, types.NewInt(amount), chain.Value(0))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State().State(), msg)
	require.NoError(td.TB(), err)

	td.State().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

func mustAddBalance(td TestDriver, from address.Address, nonce, amount uint64) {
	msg, err := td.Producer.StorageMarketAddBalance(from, nonce, chain.Value(amount))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State().State(), msg)
	require.NoError(td.TB(), err)

	td.State().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

func mustCreateStorageMarketActor(td TestDriver) address.Address {
	_, _, err := td.State().State().SetSingletonActor(actors.StorageMarketAddress, types.NewInt(0))
	require.NoError(td.TB(), err)

	mt := strgmrkt.NewMarketTracker(td.TB())
	smaddr := td.Producer.SingletonAddress(actors.StorageMarketAddress)
	td.State().AssertStorageMarketState(smaddr, strgmrkt.StorageMarketState{
		Balances:   mt.Balance,
		Deals:      mt.Deals,
		NextDealID: 0,
	})
	return smaddr
}
*/

package suites

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgmrkt"
	"github.com/filecoin-project/chain-validation/pkg/state/address"
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

	mt := strgmrkt.NewMarketTracker(c.TB())
	mt.SetMarketBalances(map[address.Address]strgmrkt.StorageParticipantBalance{
		alice: {
			Locked:    types.NewInt(0),
			Available: types.NewInt(balAddAmount),
		},
	})
	mustAddBalance(c, alice, 0, balAddAmount)
	c.Driver().AssertStorageMarketState(smaddr, strgmrkt.StorageMarketState{
		Balances:   mt.Balance,
		Deals:      mt.Deals,
		NextDealID: 0,
	})

	bob := c.Driver().NewAccountActor(initialBal)
	mt.SetMarketBalances(map[address.Address]strgmrkt.StorageParticipantBalance{
		bob: {
			Locked:    types.NewInt(0),
			Available: types.NewInt(balAddAmount),
		},
	})
	mustAddBalance(c, bob, 0, balAddAmount)
	c.Driver().AssertStorageMarketState(smaddr, strgmrkt.StorageMarketState{
		Balances:   mt.Balance,
		Deals:      mt.Deals,
		NextDealID: 0,
	})

	mustWithdrawBalance(c, bob, 1, balWithdrawAmount)
	mt.SetMarketBalances(map[address.Address]strgmrkt.StorageParticipantBalance{
		bob: {
			Locked:    types.NewInt(0),
			Available: types.NewInt(balAddAmount - balWithdrawAmount),
		},
	})
	c.Driver().AssertStorageMarketState(smaddr, strgmrkt.StorageMarketState{
		Balances:   mt.Balance,
		Deals:      mt.Deals,
		NextDealID: 0,
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
		GasUsed:     0,
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
		GasUsed:     0,
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

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

type strgmrktWrapper struct {
	T         testing.TB
	Driver    *StateDriver
	Producer  *chain.MessageProducer
	Validator *chain.Validator
	ExeCtx    *chain.ExecutionContext
}

func strgmrktTestSetup(t testing.TB, factory Factories) *strgmrktWrapper {
	drv := NewStateDriver(t, factory.NewState())
	gasPrice := types.NewInt(1)
	gasLimit := types.GasUnit(1000000)

	_, _, err := drv.State().SetSingletonActor(actors.InitAddress, types.NewInt(0))
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.BurntFundsAddress, types.NewInt(0))
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.NetworkAddress, TotalNetworkBalance)
	require.NoError(t, err)
	_, _, err = drv.State().SetSingletonActor(actors.StoragePowerAddress, types.NewInt(0))
	require.NoError(t, err)

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	testMiner := drv.NewAccountActor(0)
	exeCtx := chain.NewExecutionContext(1, testMiner)

	return &strgmrktWrapper{
		T:         t,
		Driver:    drv,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}
}

func StorageMarketActorConstructor(t testing.TB, factory Factories) {
	w := strgmrktTestSetup(t, factory)
	mustCreateStorageMarketActor(w)
}

func StorageMarketBalanceUpdates(t testing.TB, factory Factories) {
	const initialBal = 2000000000
	const balAddAmount = 100
	const balWithdrawAmount = 10

	w := strgmrktTestSetup(t, factory)
	smaddr := mustCreateStorageMarketActor(w)

	alice := w.Driver.NewAccountActor(initialBal)

	mt := strgmrkt.NewMarketTracker(w.T)
	mt.SetMarketBalances(map[address.Address]strgmrkt.StorageParticipantBalance{
		alice: {
			Locked:    types.NewInt(0),
			Available: types.NewInt(balAddAmount),
		},
	})
	mustAddBalance(w, alice, 0, balAddAmount)
	w.Driver.AssertStorageMarketState(smaddr, strgmrkt.StorageMarketState{
		Balances:   mt.Balance,
		Deals:      mt.Deals,
		NextDealID: 0,
	})

	bob := w.Driver.NewAccountActor(initialBal)
	mt.SetMarketBalances(map[address.Address]strgmrkt.StorageParticipantBalance{
		bob: {
			Locked:    types.NewInt(0),
			Available: types.NewInt(balAddAmount),
		},
	})
	mustAddBalance(w, bob, 0, balAddAmount)
	w.Driver.AssertStorageMarketState(smaddr, strgmrkt.StorageMarketState{
		Balances:   mt.Balance,
		Deals:      mt.Deals,
		NextDealID: 0,
	})

	mustWithdrawBalance(w, bob, 1, balWithdrawAmount)
	mt.SetMarketBalances(map[address.Address]strgmrkt.StorageParticipantBalance{
		bob: {
			Locked:    types.NewInt(0),
			Available: types.NewInt(balAddAmount - balWithdrawAmount),
		},
	})
	w.Driver.AssertStorageMarketState(smaddr, strgmrkt.StorageMarketState{
		Balances:   mt.Balance,
		Deals:      mt.Deals,
		NextDealID: 0,
	})
}

func mustWithdrawBalance(w *strgmrktWrapper, from address.Address, nonce, amount uint64) {
	msg, err := w.Producer.StorageMarketWithdrawBalance(from, nonce, types.NewInt(amount), chain.Value(0))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)

	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})
}

func mustAddBalance(w *strgmrktWrapper, from address.Address, nonce, amount uint64) {
	msg, err := w.Producer.StorageMarketAddBalance(from, nonce, chain.Value(amount))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)

	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})
}

func mustCreateStorageMarketActor(w *strgmrktWrapper) address.Address {
	_, _, err := w.Driver.st.SetSingletonActor(actors.StorageMarketAddress, types.NewInt(0))
	require.NoError(w.T, err)

	mt := strgmrkt.NewMarketTracker(w.T)
	smaddr := w.Producer.SingletonAddress(actors.StorageMarketAddress)
	w.Driver.AssertStorageMarketState(smaddr, strgmrkt.StorageMarketState{
		Balances:   mt.Balance,
		Deals:      mt.Deals,
		NextDealID: 0,
	})
	return smaddr
}

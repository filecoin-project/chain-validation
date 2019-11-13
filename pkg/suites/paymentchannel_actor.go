package suites

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
)

func testSetup(t *testing.T, factory Factories) (*StateDriver, *chain.MessageProducer, *chain.Validator) {
	drv := NewStateDriver(t, factory.NewState())

	_, _, err := drv.State().SetSingletonActor(state.InitAddress, big.NewInt(0))
	require.NoError(t, err)

	gasPrice := big.NewInt(1)
	gasLimit := state.GasUnit(10000)

	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	return drv, producer, validator
}

func PaymentChannelCreateSuccess(t *testing.T, factory Factories, expGasUsed uint64) {
	drv, producer, validator := testSetup(t, factory)

	expPayChAddress, err := state.NewIDAddress(103)

	alice := drv.NewAccountActor(30000)
	bob := drv.NewAccountActor(0)
	miner := drv.NewAccountActor(0) // Miner owner

	exeCtx := chain.NewExecutionContext(1, miner)

	msg, err := producer.PaymentChannelCreate(bob, alice, 0, 50)
	require.NoError(t, err)

	msgReceipt, err := validator.ApplyMessage(exeCtx, drv.State(), msg)
	require.NoError(t, err)
	drv.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: []byte(expPayChAddress),
		GasUsed:     state.GasUnit(expGasUsed),
	})

	// initial balance - paych amount - gas
	drv.AssertBalance(alice, 30000-50-expGasUsed)
	drv.AssertBalance(bob, 0)
	drv.AssertBalance(miner, expGasUsed)
	drv.AssertBalance(expPayChAddress, 50)

	// TODO make this state inspection work
	/*
		pca, err := drv.State().Actor(expPayChAddress)
		require.NoError(t, err)

		pcaStorage, err := drv.State().Storage(expPayChAddress)
		require.NoError(t, err)

		var pcs state.PaymentChannelActorState
		require.NoError(t,pcaStorage.Get(pca.Head(), &pcs))
		assert.Equal(t, []byte(alice), pcs.From)
		assert.Equal(t, []byte(bob), pcs.To)
	*/
}

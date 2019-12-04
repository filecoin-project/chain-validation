package suites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/paych"
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type paychTestingWrapper struct {
	T         testing.TB
	Driver    *StateDriver
	Producer  *chain.MessageProducer
	Validator *chain.Validator
	ExeCtx    *chain.ExecutionContext
}

func paychTestSetup(t testing.TB, factory Factories) *paychTestingWrapper {
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

	return &paychTestingWrapper{
		T:         t,
		Driver:    drv,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}

}

func PayChActorConstructor(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const valueSend = 10

	w := paychTestSetup(t, factory)

	alice := w.Driver.NewAccountActor(initialBal)
	bob := w.Driver.NewAccountActor(initialBal)

	paychAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// alice creates a payment channel with bob.
	mustCreatePaychActor(w, 0, valueSend, paychAddr, alice, bob)
	w.Driver.AssertBalance(paychAddr, valueSend)
	w.Driver.AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(0),
		ClosingAt:      0,
		MinCloseHeight: 0,
		LaneStates:     nil,
	})
}

func mustCreatePaychActor(w *paychTestingWrapper, nonce, value uint64, paychAddr, creator, paychTo address.Address) {
	paychConstructParams, err := types.Serialize(&paych.PaymentChannelConstructorParams{To: paychTo})
	require.NoError(w.T, err)

	msg, err := w.Producer.InitExec(creator, nonce, actors.PaymentChannelActorCodeCid, paychConstructParams, chain.Value(value))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)

	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: paychAddr.Bytes(),
		GasUsed:     0,
	})

}

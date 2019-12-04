package suites

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/paych"
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

// taken from lotus /build/params_shared.go
const PaymentChannelClosingDelay = 6 * 60 * 2 // six hours

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

func PayChActorUpdate(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const paychInitBal = 0
	const paychVoucherAmount = 50
	const paychUpdateBal = 100

	w := paychTestSetup(t, factory)

	alice := w.Driver.NewAccountActor(initialBal)
	bob := w.Driver.NewAccountActor(initialBal)

	paychAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// alice creates a payment channel with bob.
	mustCreatePaychActor(w, 0, paychInitBal, paychAddr, alice, bob)
	w.Driver.AssertBalance(paychAddr, paychInitBal)
	w.Driver.AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(0),
		ClosingAt:      0,
		MinCloseHeight: 0,
		LaneStates:     nil,
	})

	sv := &types.SignedVoucher{
		Nonce:  0,
		Amount: types.NewInt(paychVoucherAmount),
	}
	signMe, err := sv.SigningBytes()
	require.NoError(t, err)
	sig, err := w.Driver.State().Sign(context.TODO(), alice, signMe)
	require.NoError(t, err)
	sv.Signature = sig

	// unused but required
	proof, secret := []byte{}, []byte{}
	// alice updates the payment channel
	mustUpdatePaychActor(w, 1, paychUpdateBal, paychAddr, alice, proof, secret, *sv)
	w.Driver.AssertBalance(paychAddr, paychUpdateBal)
	w.Driver.AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(paychVoucherAmount),
		ClosingAt:      0,
		MinCloseHeight: 0,
		LaneStates: map[string]*paych.LaneState{
			"0": {
				Redeemed: types.NewInt(paychVoucherAmount),
				Closed:   false,
				Nonce:    0,
			},
		},
	})

	// alice asserts that they own the channel and the amount to send is correct.
	assertPaychOwner(w, 2, 0, paychAddr, alice, alice)
	assertPaychToSend(w, 3, 0, paychAddr, alice, types.NewInt(paychVoucherAmount))

	// alice closes the channel
	mustClosePaych(w, 4, 0, paychAddr, alice)
	w.Driver.AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(paychVoucherAmount),
		ClosingAt:      PaymentChannelClosingDelay + w.ExeCtx.Epoch,
		MinCloseHeight: 0,
		LaneStates: map[string]*paych.LaneState{
			"0": {
				Redeemed: types.NewInt(paychVoucherAmount),
				Closed:   false,
				Nonce:    0,
			},
		},
	})

	// advance the ChainEpoch to cause the channel to close
	w.ExeCtx.Epoch += PaymentChannelClosingDelay
	// bob collects the payment from alice
	const gasPaiedByBob = 360
	mustCollectPaych(w, 0, 0, paychAddr, bob)
	w.Driver.AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(0), // the funds must have moved to bob.
		ClosingAt:      PaymentChannelClosingDelay + 1,
		MinCloseHeight: 0,
		LaneStates: map[string]*paych.LaneState{
			"0": {
				Redeemed: types.NewInt(paychVoucherAmount),
				Closed:   false,
				Nonce:    0,
			},
		},
	})

	// This will break if gas ever changes..and it will..
	w.Driver.AssertBalance(bob, initialBal+paychVoucherAmount-gasPaiedByBob)

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

func mustUpdatePaychActor(w *paychTestingWrapper, nonce, value uint64, to, from address.Address, secret, proof []byte, sv types.SignedVoucher) {
	msg, err := w.Producer.PaychUpdateChannelState(to, from, nonce, sv, secret, proof, chain.Value(value))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)

	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})
}

func mustClosePaych(w *paychTestingWrapper, nonce, value uint64, paychAddr, from address.Address) {
	msg, err := w.Producer.PaychClose(paychAddr, from, nonce, chain.Value(value))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)

	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})
}

func mustCollectPaych(w *paychTestingWrapper, nonce, value uint64, paychAddr, from address.Address) {
	msg, err := w.Producer.PaychCollect(paychAddr, from, nonce, chain.Value(value))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)

	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     0,
	})

}

func assertPaychOwner(w *paychTestingWrapper, nonce, value uint64, paychAddr, from, owner address.Address) {
	msg, err := w.Producer.PaychGetOwner(paychAddr, from, nonce, chain.Value(value))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)

	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: owner.Bytes(),
		GasUsed:     0,
	})
}

func assertPaychToSend(w *paychTestingWrapper, nonce, value uint64, paychAddr, from address.Address, toSend types.BigInt) {
	msg, err := w.Producer.PaychGetToSend(paychAddr, from, nonce, chain.Value(value))
	require.NoError(w.T, err)

	msgReceipt, err := w.Validator.ApplyMessage(w.ExeCtx, w.Driver.State(), msg)
	require.NoError(w.T, err)

	w.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: toSend.Bytes(),
		GasUsed:     0,
	})
}

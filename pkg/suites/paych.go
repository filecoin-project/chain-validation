package suites

// TODO uncomment when spec settles.
/*

import (
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/paych"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

// taken from lotus /build/params_shared.go
const PaymentChannelClosingDelay = 6 * 60 * 2 // six hours

func PayChActorConstructor(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const valueSend = 10

	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})

	alice := td.Driver().NewAccountActor(initialBal)
	bob := td.Driver().NewAccountActor(initialBal)

	paychAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// alice creates a payment channel with bob.
	mustCreatePaychActor(td, 0, valueSend, paychAddr, alice, bob)
	td.Driver().AssertBalance(paychAddr, valueSend)
	td.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
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

	td := NewTestDriver(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})

	alice := td.Driver().NewAccountActor(initialBal)
	bob := td.Driver().NewAccountActor(initialBal)

	paychAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// alice creates a payment channel with bob.
	mustCreatePaychActor(td, 0, paychInitBal, paychAddr, alice, bob)
	td.Driver().AssertBalance(paychAddr, paychInitBal)
	td.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(0),
		ClosingAt:      0,
		MinCloseHeight: 0,
		LaneStates:     nil,
	})

	sv := &types.SignedVoucher{
		CallSeqNum:  0,
		Amount: types.NewInt(paychVoucherAmount),
	}
	signMe, err := sv.SigningBytes()
	require.NoError(t, err)
	sig, err := td.Driver().State().Sign(context.TODO(), alice, signMe)
	require.NoError(t, err)
	sv.Signature = sig

	// unused but required
	proof, secret := []byte{}, []byte{}
	// alice updates the payment channel
	mustUpdatePaychActor(td, 1, paychUpdateBal, paychAddr, alice, proof, secret, *sv)
	td.Driver().AssertBalance(paychAddr, paychUpdateBal)
	td.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(paychVoucherAmount),
		ClosingAt:      0,
		MinCloseHeight: 0,
		LaneStates: map[string]*paych.LaneState{
			"0": {
				Redeemed: types.NewInt(paychVoucherAmount),
				Closed:   false,
				CallSeqNum:    0,
			},
		},
	})

	// alice asserts that they own the channel and the amount to send is correct.
	assertPaychOwner(td, 2, 0, paychAddr, alice, alice)
	assertPaychToSend(td, 3, 0, paychAddr, alice, types.NewInt(paychVoucherAmount))

	// alice closes the channel
	mustClosePaych(td, 4, 0, paychAddr, alice)
	td.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(paychVoucherAmount),
		ClosingAt:      PaymentChannelClosingDelay + td.ExeCtx().Epoch,
		MinCloseHeight: 0,
		LaneStates: map[string]*paych.LaneState{
			"0": {
				Redeemed: types.NewInt(paychVoucherAmount),
				Closed:   false,
				CallSeqNum:    0,
			},
		},
	})

	// advance the ChainEpoch to cause the channel to close
	td.ExeCtx().Epoch += PaymentChannelClosingDelay
	// bob collects the payment from alice
	const gasPaiedByBob = 360
	mustCollectPaych(td, 0, 0, paychAddr, bob)
	td.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(0), // the funds must have moved to bob.
		ClosingAt:      PaymentChannelClosingDelay + 1,
		MinCloseHeight: 0,
		LaneStates: map[string]*paych.LaneState{
			"0": {
				Redeemed: types.NewInt(paychVoucherAmount),
				Closed:   false,
				CallSeqNum:    0,
			},
		},
	})

	// This will break if gas ever changes..and it will..
	td.Driver().AssertBalance(bob, initialBal+paychVoucherAmount-gasPaiedByBob)

}

func mustCreatePaychActor(td TestDriver, nonce, value uint64, paychAddr, creator, paychTo address.Address) {
	paychConstructParams, err := state.Serialize(&paych.PaymentChannelConstructorParams{To: paychTo})
	require.NoError(td.TB(), err)

	msg, err := td.Producer().InitExec(creator, nonce, actors.PaymentChannelActorCodeCid, paychConstructParams, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: paychAddr.Bytes(),
		GasUsed:     types.NewInt(0),
	})
}

func mustUpdatePaychActor(td TestDriver, nonce, value uint64, to, from address.Address, secret, proof []byte, sv types.SignedVoucher) {
	msg, err := td.Producer().PaychUpdateChannelState(to, from, nonce, sv, secret, proof, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

func mustClosePaych(td TestDriver, nonce, value uint64, paychAddr, from address.Address) {
	msg, err := td.Producer().PaychClose(paychAddr, from, nonce, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

func mustCollectPaych(td TestDriver, nonce, value uint64, paychAddr, from address.Address) {
	msg, err := td.Producer().PaychCollect(paychAddr, from, nonce, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})

}

func assertPaychOwner(td TestDriver, nonce, value uint64, paychAddr, from, owner address.Address) {
	msg, err := td.Producer().PaychGetOwner(paychAddr, from, nonce, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: owner.Bytes(),
		GasUsed:     types.NewInt(0),
	})
}

func assertPaychToSend(td TestDriver, nonce, value uint64, paychAddr, from address.Address, toSend types.BigInt) {
	msg, err := td.Producer().PaychGetToSend(paychAddr, from, nonce, chain.Value(value))
	require.NoError(td.TB(), err)

	msgReceipt, err := td.Validator().ApplyMessage(td.ExeCtx(), td.Driver().State(), msg)
	require.NoError(td.TB(), err)

	td.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: toSend.Bytes(),
		GasUsed:     types.NewInt(0),
	})
}
*/

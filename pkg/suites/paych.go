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

func PayChActorConstructor(t testing.TB, factory Factories) {
	const initialBal = 200000000000
	const valueSend = 10

	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})

	alice := c.Driver().NewAccountActor(initialBal)
	bob := c.Driver().NewAccountActor(initialBal)

	paychAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// alice creates a payment channel with bob.
	mustCreatePaychActor(c, 0, valueSend, paychAddr, alice, bob)
	c.Driver().AssertBalance(paychAddr, valueSend)
	c.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
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

	c := NewCandy(t, factory, map[actors.SingletonActorID]types.BigInt{
		actors.InitAddress:         types.NewInt(0),
		actors.BurntFundsAddress:   types.NewInt(0),
		actors.StoragePowerAddress: types.NewInt(0),
		actors.NetworkAddress:      TotalNetworkBalance,
	})

	alice := c.Driver().NewAccountActor(initialBal)
	bob := c.Driver().NewAccountActor(initialBal)

	paychAddr, err := address.NewIDAddress(103)
	require.NoError(t, err)

	// alice creates a payment channel with bob.
	mustCreatePaychActor(c, 0, paychInitBal, paychAddr, alice, bob)
	c.Driver().AssertBalance(paychAddr, paychInitBal)
	c.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
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
	sig, err := c.Driver().State().Sign(context.TODO(), alice, signMe)
	require.NoError(t, err)
	sv.Signature = sig

	// unused but required
	proof, secret := []byte{}, []byte{}
	// alice updates the payment channel
	mustUpdatePaychActor(c, 1, paychUpdateBal, paychAddr, alice, proof, secret, *sv)
	c.Driver().AssertBalance(paychAddr, paychUpdateBal)
	c.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
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
	assertPaychOwner(c, 2, 0, paychAddr, alice, alice)
	assertPaychToSend(c, 3, 0, paychAddr, alice, types.NewInt(paychVoucherAmount))

	// alice closes the channel
	mustClosePaych(c, 4, 0, paychAddr, alice)
	c.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
		From:           alice,
		To:             bob,
		ToSend:         types.NewInt(paychVoucherAmount),
		ClosingAt:      PaymentChannelClosingDelay + c.ExeCtx().Epoch,
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
	c.ExeCtx().Epoch += PaymentChannelClosingDelay
	// bob collects the payment from alice
	const gasPaiedByBob = 360
	mustCollectPaych(c, 0, 0, paychAddr, bob)
	c.Driver().AssertPayChState(paychAddr, paych.PaymentChannelActorState{
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
	c.Driver().AssertBalance(bob, initialBal+paychVoucherAmount-gasPaiedByBob)

}

func mustCreatePaychActor(c Candy, nonce, value uint64, paychAddr, creator, paychTo address.Address) {
	paychConstructParams, err := types.Serialize(&paych.PaymentChannelConstructorParams{To: paychTo})
	require.NoError(c.TB(), err)

	msg, err := c.Producer().InitExec(creator, nonce, actors.PaymentChannelActorCodeCid, paychConstructParams, chain.Value(value))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)

	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: paychAddr.Bytes(),
		GasUsed:     types.NewInt(0),
	})
}

func mustUpdatePaychActor(c Candy, nonce, value uint64, to, from address.Address, secret, proof []byte, sv types.SignedVoucher) {
	msg, err := c.Producer().PaychUpdateChannelState(to, from, nonce, sv, secret, proof, chain.Value(value))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)

	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

func mustClosePaych(c Candy, nonce, value uint64, paychAddr, from address.Address) {
	msg, err := c.Producer().PaychClose(paychAddr, from, nonce, chain.Value(value))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)

	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})
}

func mustCollectPaych(c Candy, nonce, value uint64, paychAddr, from address.Address) {
	msg, err := c.Producer().PaychCollect(paychAddr, from, nonce, chain.Value(value))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)

	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: nil,
		GasUsed:     types.NewInt(0),
	})

}

func assertPaychOwner(c Candy, nonce, value uint64, paychAddr, from, owner address.Address) {
	msg, err := c.Producer().PaychGetOwner(paychAddr, from, nonce, chain.Value(value))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)

	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: owner.Bytes(),
		GasUsed:     types.NewInt(0),
	})
}

func assertPaychToSend(c Candy, nonce, value uint64, paychAddr, from address.Address, toSend types.BigInt) {
	msg, err := c.Producer().PaychGetToSend(paychAddr, from, nonce, chain.Value(value))
	require.NoError(c.TB(), err)

	msgReceipt, err := c.Validator().ApplyMessage(c.ExeCtx(), c.Driver().State(), msg)
	require.NoError(c.TB(), err)

	c.Driver().AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: toSend.Bytes(),
		GasUsed:     types.NewInt(0),
	})
}

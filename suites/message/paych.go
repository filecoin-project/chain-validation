package message

import (
	"context"
	"testing"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	paych_spec "github.com/filecoin-project/specs-actors/actors/builtin/paych"
	crypto_spec "github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"
	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestPaych(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	var initialBal = abi_spec.NewTokenAmount(200_000_000_000)
	var toSend = abi_spec.NewTokenAmount(10_000)
	var createPaychGasCost = abi_spec.NewTokenAmount(1416)
	t.Run("happy path constructor", func(t *testing.T) {
		td := builder.Build(t)

		// will create and send on payment channel
		sender, senderID := td.NewAccountActor(drivers.SECP, initialBal)

		// will be receiver on paych
		receiver, receiverID := td.NewAccountActor(drivers.SECP, initialBal)

		// the _expected_ address of the payment channel
		paychAddr := utils.NewIDAddr(t, utils.IdFromAddress(receiverID)+1)
		createRet := td.ComputeInitActorExecReturn(sender, 0, 0, paychAddr)

		// init actor creates the payment channel
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: createPaychGasCost},
		)

		var pcState paych_spec.State
		td.GetActorState(paychAddr, &pcState)
		assert.Equal(t, senderID, pcState.From)
		assert.Equal(t, receiverID, pcState.To)
		td.AssertBalance(paychAddr, toSend)
	})

	t.Run("happy path update", func(t *testing.T) {
		td := builder.Build(t)

		const pcTimeLock = abi_spec.ChainEpoch(1)
		const pcLane = uint64(123)
		const pcNonce = uint64(1)
		var pcAmount = big_spec.NewInt(10)
		var pcSig = &crypto_spec.Signature{
			Type: crypto_spec.SigTypeBLS,
			Data: []byte("signature goes here"), // TODO may need to generate an actual signature
		}

		// will create and send on payment channel
		sender, _ := td.NewAccountActor(drivers.SECP, initialBal)

		// will be receiver on paych
		receiver, receiverID := td.NewAccountActor(drivers.SECP, initialBal)

		// the _expected_ address of the payment channel
		paychAddr := utils.NewIDAddr(t, utils.IdFromAddress(receiverID)+1)
		createRet := td.ComputeInitActorExecReturn(sender, 0, 0, paychAddr)
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: createPaychGasCost},
		)

		td.ApplyMessageExpectReceipt(
			td.MessageProducer.PaychUpdateChannelState(paychAddr, sender, paych_spec.UpdateChannelStateParams{
				Sv: paych_spec.SignedVoucher{
					TimeLockMin: pcTimeLock,
					TimeLockMax: 0, // TimeLockMax set to 0 means no timeout
					Lane:        pcLane,
					Nonce:       pcNonce,
					Amount:      pcAmount,
					Signature:   pcSig,
				},
			}, chain.Nonce(1), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: abi_spec.NewTokenAmount(319)},
		)
		var pcState paych_spec.State
		td.GetActorState(paychAddr, &pcState)
		assert.Equal(t, 1, len(pcState.LaneStates))
		ls := pcState.LaneStates[0]
		assert.Equal(t, pcAmount, ls.Redeemed)
		assert.Equal(t, pcNonce, ls.Nonce)
		assert.Equal(t, pcLane, ls.ID)
	})

	t.Run("happy path collect", func(t *testing.T) {
		td := builder.Build(t)

		// create the payment channel
		sender, _ := td.NewAccountActor(drivers.SECP, initialBal)
		receiver, receiverID := td.NewAccountActor(drivers.SECP, initialBal)
		paychAddr := utils.NewIDAddr(t, utils.IdFromAddress(receiverID)+1)
		initRet := td.ComputeInitActorExecReturn(sender, 0, 0, paychAddr)
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&initRet), GasUsed: createPaychGasCost},
		)
		td.AssertBalance(paychAddr, toSend)

		updatePaychGasCost := abi_spec.NewTokenAmount(321)
		settlePayChGasCost := abi_spec.NewTokenAmount(234)
		collectPaychGasCost := abi_spec.NewTokenAmount(271)

		td.ApplyMessageExpectReceipt(
			td.MessageProducer.PaychUpdateChannelState(paychAddr, sender, paych_spec.UpdateChannelStateParams{
				Sv: paych_spec.SignedVoucher{
					TimeLockMin: abi_spec.ChainEpoch(1),
					TimeLockMax: 0, // TimeLockMax set to 0 means no timeout
					Lane:        1,
					Nonce:       1,
					Amount:      toSend, // the amount that can be redeemed by receiver
					Signature: &crypto_spec.Signature{
						Type: crypto_spec.SigTypeBLS,
						Data: []byte("signature goes here"),
					},
				},
			}, chain.Nonce(1), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: updatePaychGasCost},
		)

		// settle the payment channel so it may be collected
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.PaychSettle(paychAddr, receiver, adt_spec.EmptyValue{}, chain.Value(big_spec.Zero()), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: settlePayChGasCost},
		)

		// advance the epoch so the funds may be redeemed.
		td.ExeCtx.Epoch++

		td.ApplyMessageExpectReceipt(
			td.MessageProducer.PaychCollect(paychAddr, receiver, adt_spec.EmptyValue{}, chain.Nonce(1), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: collectPaychGasCost},
		)

		// receiver_balance = initial_balance + paych_send - settle_paych_msg_gas - collect_paych_msg_gas
		td.AssertBalance(receiver, big_spec.Sub(big_spec.Sub(big_spec.Add(toSend, initialBal), settlePayChGasCost), collectPaychGasCost))
		td.AssertBalance(paychAddr, big_spec.Zero())
		var pcState paych_spec.State
		td.GetActorState(paychAddr, &pcState)
		assert.Equal(t, big_spec.Zero(), pcState.ToSend)
	})

}

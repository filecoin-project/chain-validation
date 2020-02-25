package suites

import (
	"context"
	"testing"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	account_spec "github.com/filecoin-project/specs-actors/actors/builtin/account"
	init_spec "github.com/filecoin-project/specs-actors/actors/builtin/init"
	paych_spec "github.com/filecoin-project/specs-actors/actors/builtin/paych"
	reward_spec "github.com/filecoin-project/specs-actors/actors/builtin/reward"
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
		WithActorState([]drivers.ActorState{
			{
				Addr:    builtin_spec.InitActorAddr,
				Balance: big_spec.Zero(),
				Code:    builtin_spec.InitActorCodeID,
				State:   init_spec.ConstructState(drivers.EmptyMapCid, "chain-validation"),
			},
			{
				Addr:    builtin_spec.RewardActorAddr,
				Balance: TotalNetworkBalance,
				Code:    builtin_spec.RewardActorCodeID,
				State:   reward_spec.ConstructState(drivers.EmptyMultiMapCid),
			},
			{
				Addr:    builtin_spec.BurntFundsActorAddr,
				Balance: big_spec.Zero(),
				Code:    builtin_spec.AccountActorCodeID,
				State:   &account_spec.State{Address: builtin_spec.BurntFundsActorAddr},
			},
		})

	var initialBal = abi_spec.NewTokenAmount(200_000_000_000)
	var toSend = abi_spec.NewTokenAmount(10_000)
	t.Run("happy path constructor", func(t *testing.T) {
		td := builder.Build(t)

		// will create and send on payment channel
		sender, senderID := td.NewAccountActor(drivers.SECP, initialBal) // 100

		// will be receiver on paych
		receiver, receiverID := td.NewAccountActor(drivers.SECP, initialBal) // 101

		// the _expected_ address of the payment channel
		paychAddr := utils.NewIDAddr(t, 102) // 102
		createRet := td.ComputeInitActorExecReturn(senderID, 0, paychAddr)

		// init actor creates the payment channel
		td.ApplyMessageExpectReceipt(
			td.Producer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: big_spec.Zero()},
		)

		var pcState paych_spec.State
		td.GetActorState(paychAddr, &pcState)
		assert.Equal(t, senderID, pcState.From)
		assert.Equal(t, receiverID, pcState.To)
		td.AssertBalance(paychAddr, toSend)
	})

	t.Run("happy path update", func(t *testing.T) {
		td := builder.Build(t)

		const pcTimeLock = abi_spec.ChainEpoch(10)
		const pcLane = uint64(123)
		const pcNonce = uint64(1)
		var pcAmount = big_spec.NewInt(10)
		var pcSig = &crypto_spec.Signature{
			Type: crypto_spec.SigTypeBLS,
			Data: []byte("signature goes here"), // TODO may need to generate an actual signature
		}

		// will create and send on payment channel
		sender, senderID := td.NewAccountActor(drivers.SECP, initialBal) // 100

		// will be receiver on paych
		receiver, _ := td.NewAccountActor(drivers.SECP, initialBal) // 101

		// the _expected_ address of the payment channel
		paychAddr := utils.NewIDAddr(t, 102) // 102
		createRet := td.ComputeInitActorExecReturn(senderID, 0, paychAddr)
		td.ApplyMessageExpectReceipt(
			td.Producer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: big_spec.Zero()},
		)

		td.ApplyMessageExpectReceipt(
			td.Producer.PaychUpdateChannelState(paychAddr, sender, paych_spec.UpdateChannelStateParams{
				Sv: paych_spec.SignedVoucher{
					TimeLock:  pcTimeLock,
					Lane:      pcLane,
					Nonce:     pcNonce,
					Amount:    pcAmount,
					Signature: pcSig,
				},
			}, chain.Nonce(1), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: EmptyReturnValue, GasUsed: big_spec.Zero()},
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
		sender, _ := td.NewAccountActor(drivers.SECP, initialBal)   // 101
		receiver, _ := td.NewAccountActor(drivers.SECP, initialBal) // 102
		paychAddr := utils.NewIDAddr(t, 103)                        // 103
		td.ApplyMessageExpectReceipt(
			td.Producer.CreatePaymentChannelActor(sender, receiver, chain.Value(toSend), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: paychAddr.Bytes(), GasUsed: big_spec.Zero()},
		)
		td.AssertBalance(paychAddr, toSend)

		td.ApplyMessageExpectReceipt(
			td.Producer.PaychCollect(paychAddr, receiver, adt_spec.EmptyValue{}, chain.Nonce(1)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: EmptyReturnValue, GasUsed: big_spec.Zero()},
		)

		td.AssertBalance(receiver, toSend)
		td.AssertBalance(paychAddr, big_spec.Zero())
		var pcState paych_spec.State
		td.GetActorState(paychAddr, &pcState)
		assert.Equal(t, big_spec.Zero(), pcState.ToSend)
	})

}

package suites

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	paych_spec "github.com/filecoin-project/specs-actors/actors/builtin/paych"
	crypto_spec "github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestPaych(t *testing.T, factory state.Factories) {
	defaultMiner := utils.NewSECP256K1Addr(t, "miner")

	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(big_spec.NewInt(1000000)).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithSingletonActors(map[address.Address]big_spec.Int{
			builtin_spec.InitActorAddr:         big_spec.NewInt(0),
			builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
			builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
			builtin_spec.RewardActorAddr:       TotalNetworkBalance,
		}).
		WithDefaultMiner(defaultMiner)

	var initialBal = abi_spec.NewTokenAmount(200_000_000_000)
	toSend := abi_spec.NewTokenAmount(10_000)
	t.Run("happy path constructor", func(t *testing.T) {
		td := builder.Build(t)

		// will create and send on payment channel
		sender := td.NewAccountActor(drivers.SECP, initialBal) // 101
		// will be receiver on paych
		receiver := td.NewAccountActor(drivers.SECP, initialBal) // 102
		// the _expected_ address of the payment channel
		paychAddr := utils.NewIDAddr(t, 103) // 103

		// init actor creates the payment channel
		td.ApplyMessageExpectReceipt(
			td.Producer.CreatePaymentChannel(sender, receiver, chain.Value(toSend), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: 0, ReturnValue: paychAddr.Bytes(), GasUsed: big_spec.Zero()},
		)

		var pcState paych_spec.State
		td.GetActorState(paychAddr, &pcState)
		assert.Equal(t, sender, pcState.To)
		assert.Equal(t, receiver, pcState.From)
		assert.Equal(t, toSend, pcState.ToSend)
	})

	t.Run("happy path update", func(t *testing.T) {
		td := builder.Build(t)

		const pcTimeLock = abi_spec.ChainEpoch(10)
		const pcLane = int64(123)
		const pcNonce = int64(1)
		var pcAmount = big_spec.NewInt(10)
		var pcSig = &crypto_spec.Signature{
			Type: crypto_spec.SigTypeBLS,
			Data: []byte("does this matter??!!"),
		}

		// create the payment channel
		sender := td.NewAccountActor(drivers.SECP, initialBal)   // 101
		receiver := td.NewAccountActor(drivers.SECP, initialBal) // 102
		paychAddr := utils.NewIDAddr(t, 103)                     // 103
		td.ApplyMessageExpectReceipt(
			td.Producer.CreatePaymentChannel(sender, receiver, chain.Value(toSend), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: 0, ReturnValue: paychAddr.Bytes(), GasUsed: big_spec.Zero()},
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
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: nil, GasUsed: big_spec.Zero()},
		)
		var pcState paych_spec.State
		td.GetActorState(paychAddr, &pcState)
		assert.Equal(t, 1, len(pcState.LaneStates))
		ls := pcState.LaneStates[0]
		assert.Equal(t, pcAmount, ls.Redeemed)
		assert.Equal(t, pcNonce, ls.Nonce)
		assert.Equal(t, pcLane, ls.ID)
	})
}

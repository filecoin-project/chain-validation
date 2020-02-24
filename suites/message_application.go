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
	crypto_spec "github.com/filecoin-project/specs-actors/actors/crypto"

	//paych_spec "github.com/filecoin-project/specs-actors/actors/builtin/paych"
	reward_spec "github.com/filecoin-project/specs-actors/actors/builtin/reward"
	//crypto_spec "github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestMessageApplicationEdgecases(t *testing.T, factory state.Factories) {
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
	t.Run("fail to cover gas cost for message receipt on chain", func(t *testing.T) {
		td := builder.Build(t)

		var aliceBal = abi_spec.NewTokenAmount(1_000_000_000)
		var transferAmnt = abi_spec.NewTokenAmount(10)
		var gasCost = big_spec.Zero()

		alice := td.NewAccountActor(drivers.SECP, aliceBal)

		td.ApplyMessageExpectReceipt(
			td.Producer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(0), chain.GasPrice(1), chain.GasLimit(8)),
			types.MessageReceipt{ExitCode: exitcode.SysErrOutOfGas, ReturnValue: EmptyReturnValue, GasUsed: gasCost},
		)
	})
	t.Run("not enough gas to pay message on-chain-size cost", func(t *testing.T) {
		td := builder.Build(t)

		var aliceBal = abi_spec.NewTokenAmount(1_000_000_000)
		var transferAmnt = abi_spec.NewTokenAmount(10)
		var gasCost = big_spec.Zero()

		alice := td.NewAccountActor(drivers.SECP, aliceBal)

		// Expect Message application to fail due to lack of gas
		td.ApplyMessageExpectReceipt(
			td.Producer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(0), chain.GasPrice(10), chain.GasLimit(1)),
			types.MessageReceipt{ExitCode: exitcode.SysErrOutOfGas, ReturnValue: EmptyReturnValue, GasUsed: gasCost},
		)

		// Expect Message application to fail due to lack of gas when sender address is unknown
		unknown := utils.NewIDAddr(t, 10000000)
		td.ApplyMessageExpectReceipt(
			td.Producer.Transfer(alice, unknown, chain.Value(transferAmnt), chain.Nonce(0), chain.GasPrice(10), chain.GasLimit(1)),
			types.MessageReceipt{ExitCode: exitcode.SysErrOutOfGas, ReturnValue: EmptyReturnValue, GasUsed: gasCost},
		)

	})

	t.Run("invalid actor CallSeqNum", func(t *testing.T) {
		td := builder.Build(t)

		var aliceBal = abi_spec.NewTokenAmount(1_000_000_000)
		var transferAmnt = abi_spec.NewTokenAmount(10)
		var gasCost = big_spec.Zero()

		alice := td.NewAccountActor(drivers.SECP, aliceBal)

		// Expect Message application to fail due to callseqnum being invalid: 1 instead of 0
		td.ApplyMessageExpectReceipt(
			td.Producer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(1)),
			types.MessageReceipt{ExitCode: exitcode.SysErrInvalidCallSeqNum, ReturnValue: EmptyReturnValue, GasUsed: gasCost},
		)

		// Expect message application to fail due to unknow actor when call seq num is also incorrect
		unknown := utils.NewIDAddr(t, 10000000)
		td.ApplyMessageExpectReceipt(
			td.Producer.Transfer(alice, unknown, chain.Value(transferAmnt), chain.Nonce(1)),
			types.MessageReceipt{ExitCode: exitcode.SysErrActorNotFound, ReturnValue: EmptyReturnValue, GasUsed: gasCost},
		)
	})

	t.Run("abort during actor execution", func(t *testing.T) {
		td := builder.Build(t)

		const pcTimeLock = abi_spec.ChainEpoch(10)
		const pcLane = uint64(123)
		const pcNonce = uint64(1)
		var pcAmount = big_spec.NewInt(10)
		var initialBal = abi_spec.NewTokenAmount(200_000_000_000)
		var toSend = abi_spec.NewTokenAmount(10_000)
		var pcSig = &crypto_spec.Signature{
			Type: crypto_spec.SigTypeBLS,
			Data: []byte("Grrr im an invalid signature, I cause panics in the payment channel actor"),
		}

		// will create and send on payment channel
		sender := td.NewAccountActor(drivers.SECP, initialBal) // 100
		senderID := utils.NewIDAddr(t, 100)

		// will be receiver on paych
		receiver := td.NewAccountActor(drivers.SECP, initialBal) // 101
		//receiverID := utils.NewIDAddr(t, 101)

		// the _expected_ address of the payment channel
		paychAddr := utils.NewIDAddr(t, 102) // 103
		createRet := td.ComputeInitActorExecReturn(senderID, 0, paychAddr)
		td.ApplyMessageExpectReceipt(
			td.Producer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
			types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: big_spec.Zero()},
		)

		// message application fails due to invalid argument (signature).
		td.ApplyMessageExpectReceipt(
			td.Producer.PaychUpdateChannelState(paychAddr, sender, paych_spec.UpdateChannelStateParams{
				Sv: paych_spec.SignedVoucher{
					TimeLock:  pcTimeLock,
					Lane:      pcLane,
					Nonce:     pcNonce,
					Amount:    pcAmount,
					Signature: pcSig, // construct with invalid signature
				},
			}, chain.Nonce(1), chain.Value(big_spec.Zero())),
			types.MessageReceipt{ExitCode: exitcode.ErrIllegalArgument, ReturnValue: EmptyReturnValue, GasUsed: big_spec.Zero()},
		)
	})
}

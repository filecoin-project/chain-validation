package message

import (
	"context"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	exitcode_spec "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
)

func InitActorCreatesActorsWithSequentialIDAddresses(h *drivers.ValidationHarness, factory state.Factories) {
	td := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState([]drivers.ActorState{
			drivers.DefaultInitActorState,
			drivers.DefaultRewardActorState,
			drivers.DefaultBurntFundsActorState,
			drivers.DefaultStoragePowerActorState,
		}).Build(h)

	var initialBal = abi_spec.NewTokenAmount(200_000_000_000)
	var toSend = abi_spec.NewTokenAmount(10_000)

	sender, senderID := td.NewAccountActor(drivers.SECP, initialBal)

	receiver, receiverID := td.NewAccountActor(drivers.SECP, initialBal)

	firstPaychAddr := drivers.NewIDAddr(h, drivers.IdFromAddress(receiverID)+1)
	secondPaychAddr := drivers.NewIDAddr(h, drivers.IdFromAddress(receiverID)+2)

	firstInitRet := td.ComputeInitActorExecReturn(senderID, 0, firstPaychAddr)
	secondInitRet := td.ComputeInitActorExecReturn(senderID, 1, secondPaychAddr)

	td.ApplyMessageExpectReceipt(
		td.MessageProducer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
		types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: chain.MustSerialize(&firstInitRet), GasUsed: big_spec.Zero()},
	)

	td.ApplyMessageExpectReceipt(
		td.MessageProducer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(1)),
		types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: chain.MustSerialize(&secondInitRet), GasUsed: big_spec.Zero()},
	)
}

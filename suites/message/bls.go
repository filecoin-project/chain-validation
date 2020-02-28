package message

import (
	"context"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	exitcode_spec "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
)

func SuccessfullyCreateBLSAccountActor(h *drivers.ValidationHarness, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState([]drivers.ActorState{
			drivers.DefaultInitActorState,
			drivers.DefaultRewardActorState,
			drivers.DefaultBurntFundsActorState,
			drivers.DefaultStoragePowerActorState,
		})
	td := builder.Build(h)

	balance := abi_spec.NewTokenAmount(10_000_000)
	send := abi_spec.NewTokenAmount(10_000)

	existingAccountAddr, _ := td.NewAccountActor(address.BLS, balance)
	td.ApplyMessageExpectReceipt(
		td.MessageProducer.Transfer(drivers.NewBLSAddr(h, 1), existingAccountAddr, chain.Value(send), chain.Nonce(0)),
		types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.Zero()},
	)
}

func FailToCreateBLSAccountActorWithInsufficientBalance(h *drivers.ValidationHarness, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState([]drivers.ActorState{
			drivers.DefaultInitActorState,
			drivers.DefaultRewardActorState,
			drivers.DefaultBurntFundsActorState,
			drivers.DefaultStoragePowerActorState,
		})
	td := builder.Build(h)

	balance := abi_spec.NewTokenAmount(9_999)
	send := abi_spec.NewTokenAmount(10_000)

	existingAccountAddr, _ := td.NewAccountActor(address.BLS, balance)
	td.ApplyMessageExpectReceipt(
		td.MessageProducer.Transfer(drivers.NewBLSAddr(h, 1), existingAccountAddr, chain.Value(send), chain.Nonce(0)),
		types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: big_spec.Zero()},
	)

}

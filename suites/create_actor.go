package suites

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	exitcode_spec "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestAccountActorCreation(t *testing.T, factory state.Factories) {
	defaultMiner := utils.NewBLSAddr(t, 123)

	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(big_spec.NewInt(1_000_000)).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithSingletonActors(map[address.Address]big_spec.Int{
			builtin_spec.InitActorAddr:         big_spec.NewInt(0),
			builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
			builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
			builtin_spec.RewardActorAddr:       TotalNetworkBalance,
		}).
		WithDefaultMiner(defaultMiner)

	testCases := []struct {
		desc string

		existingActorType address.Protocol
		existingActorBal  abi_spec.TokenAmount

		newActorAddr    address.Address
		newActorInitBal abi_spec.TokenAmount

		expGasCost  abi_spec.TokenAmount
		expExitCode exitcode_spec.ExitCode
	}{
		{
			"success create SECP256K1 account actor",
			address.SECP256K1,
			abi_spec.NewTokenAmount(10_000_000),

			utils.NewSECP256K1Addr(t, "publickeyfoo"),
			abi_spec.NewTokenAmount(10_000),

			abi_spec.NewTokenAmount(0),
			exitcode_spec.Ok,
		},
		{
			"success create BLS account actor",
			address.SECP256K1,
			abi_spec.NewTokenAmount(10_000_000),

			utils.NewBLSAddr(t, 1),
			abi_spec.NewTokenAmount(10_000),

			abi_spec.NewTokenAmount(0),
			exitcode_spec.Ok,
		},
		{
			"fail create SECP256K1 account actor insufficient balance",
			address.SECP256K1,
			abi_spec.NewTokenAmount(9_999),

			utils.NewSECP256K1Addr(t, "publickeybar"),
			abi_spec.NewTokenAmount(10_000),

			abi_spec.NewTokenAmount(0),
			exitcode_spec.Ok,
		},
		{
			"fail create BLS account actor insufficient balance",
			address.BLS,
			abi_spec.NewTokenAmount(9_999),

			utils.NewSECP256K1Addr(t, "publickeybaz"),
			abi_spec.NewTokenAmount(10_000),

			abi_spec.NewTokenAmount(0),
			exitcode_spec.Ok,
		},
		// TODO add edge case tests that have insufficient balance after gas fees
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			td := builder.Build(t)

			existingAccountAddr := td.NewAccountActor(tc.existingActorType, tc.existingActorBal)
			td.ApplyMessageExpectReceipt(
				td.Producer.Transfer(tc.newActorAddr, existingAccountAddr, chain.Value(tc.newActorInitBal), chain.Nonce(0)),
				types.MessageReceipt{ExitCode: tc.expExitCode, ReturnValue: drivers.EmptyRetrunValueBytes, GasUsed: tc.expGasCost},
			)

			td.AssertBalance(tc.newActorAddr, tc.newActorInitBal)
			td.AssertBalance(existingAccountAddr, big_spec.Sub(tc.existingActorBal, tc.expGasCost))
		})
	}
}

func TestInitActorSequentialIDAddressCreate(t *testing.T, factory state.Factories) {
	defaultMiner := utils.NewBLSAddr(t, 123)

	td := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(big_spec.NewInt(1_000_000)).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithSingletonActors(map[address.Address]big_spec.Int{
			builtin_spec.InitActorAddr:         big_spec.NewInt(0),
			builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
			builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
			builtin_spec.RewardActorAddr:       TotalNetworkBalance,
		}).
		WithDefaultMiner(defaultMiner).
		Build(t)

	var initialBal = abi_spec.NewTokenAmount(200_000_000_000)
	var toSend = abi_spec.NewTokenAmount(10_000)

	sender := td.NewAccountActor(drivers.SECP, initialBal)   // 101
	receiver := td.NewAccountActor(drivers.SECP, initialBal) // 102
	firstPaychAddr := utils.NewIDAddr(t, 103)                // 103
	secondPaychAddr := utils.NewIDAddr(t, 104)               // 104

	td.ApplyMessageExpectReceipt(
		td.Producer.CreatePaymentChannelActor(sender, receiver, chain.Value(toSend), chain.Nonce(0)),
		types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: firstPaychAddr.Bytes(), GasUsed: big_spec.Zero()},
	)

	td.ApplyMessageExpectReceipt(
		td.Producer.CreatePaymentChannelActor(sender, receiver, chain.Value(toSend), chain.Nonce(1)),
		types.MessageReceipt{ExitCode: exitcode_spec.Ok, ReturnValue: secondPaychAddr.Bytes(), GasUsed: big_spec.Zero()},
	)
}

package message

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	exitcode_spec "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestAccountActorCreation(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

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

			abi_spec.NewTokenAmount(131),
			exitcode_spec.Ok,
		},
		{
			"success create BLS account actor",
			address.SECP256K1,
			abi_spec.NewTokenAmount(10_000_000),

			utils.NewBLSAddr(t, 1),
			abi_spec.NewTokenAmount(10_000),

			abi_spec.NewTokenAmount(188),
			exitcode_spec.Ok,
		},
		{
			"fail create SECP256K1 account actor insufficient balance",
			address.SECP256K1,
			abi_spec.NewTokenAmount(9_999),

			utils.NewSECP256K1Addr(t, "publickeybar"),
			abi_spec.NewTokenAmount(10_000),

			abi_spec.NewTokenAmount(1_000_000),
			exitcode_spec.SysErrInsufficientFunds,
		},
		{
			"fail create BLS account actor insufficient balance",
			address.SECP256K1,
			abi_spec.NewTokenAmount(9_999),

			utils.NewBLSAddr(t, 1),
			abi_spec.NewTokenAmount(10_000),

			abi_spec.NewTokenAmount(1_000_000),
			exitcode_spec.SysErrInsufficientFunds,
		},
		// TODO add edge case tests that have insufficient balance after gas fees
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			td := builder.Build(t)
			defer td.Complete()

			existingAccountAddr, _ := td.NewAccountActor(tc.existingActorType, tc.existingActorBal)
			gasUsed := td.ApplyMessageExpectCode(
				td.MessageProducer.Transfer(tc.newActorAddr, existingAccountAddr, chain.Value(tc.newActorInitBal), chain.Nonce(0)),
				tc.expExitCode,
			)

			// new actor balance will only exist if message was applied successfully.
			if tc.expExitCode.IsSuccess() {
				td.AssertBalance(tc.newActorAddr, tc.newActorInitBal)
				td.AssertBalance(existingAccountAddr, big_spec.Sub(big_spec.Sub(tc.existingActorBal, big_spec.NewInt(gasUsed)), tc.newActorInitBal))
			}
		})
	}
}

func TestInitActorSequentialIDAddressCreate(t *testing.T, factory state.Factories) {
	td := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState).Build(t)
	defer td.Complete()

	var initialBal = abi_spec.NewTokenAmount(200_000_000_000)
	var toSend = abi_spec.NewTokenAmount(10_000)

	sender, _ := td.NewAccountActor(drivers.SECP, initialBal)

	receiver, receiverID := td.NewAccountActor(drivers.SECP, initialBal)

	firstPaychAddr := utils.NewIDAddr(t, utils.IdFromAddress(receiverID)+1)
	secondPaychAddr := utils.NewIDAddr(t, utils.IdFromAddress(receiverID)+2)

	firstInitRet := td.ComputeInitActorExecReturn(sender, 0, 0, firstPaychAddr)
	secondInitRet := td.ComputeInitActorExecReturn(sender, 1, 0, secondPaychAddr)

	td.ApplyMessageExpectSuccessAndReturn(
		td.MessageProducer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
		chain.MustSerialize(&firstInitRet),
	)

	td.ApplyMessageExpectSuccessAndReturn(
		td.MessageProducer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
		chain.MustSerialize(&secondInitRet),
	)
}

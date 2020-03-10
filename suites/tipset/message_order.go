package tipset

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
)

func TestBlockMessageSendOrdering(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("BLS Message application happens before SECP message application", func(t *testing.T) {
		td := builder.Build(t)

		senderInitialBal := abi_spec.NewTokenAmount(10_000_000)
		senderTransferAmnt := abi_spec.NewTokenAmount(500_000)
		middleTransferAmnt := abi_spec.NewTokenAmount(10_000)
		sender, _ := td.NewAccountActor(address.BLS, senderInitialBal)
		middle, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())
		final, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkMsgs := chain.NewTipSetMessageBuilder().
			// miner addresses are required to use ID protocol.
			WithMiner(td.ExeCtx.Miner).
			// send value from sender to middle
			WithBLSMessage(td.MessageProducer.Transfer(middle, sender, chain.Nonce(0), chain.Value(senderTransferAmnt))).
			WithSECPMessage(signMessage(td.MessageProducer.Transfer(final, middle, chain.Nonce(0), chain.Value(middleTransferAmnt)), td.Wallet())).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 2)

		td.AssertReceipt(types.MessageReceipt{
			ExitCode:    exitcode.Ok,
			ReturnValue: drivers.EmptyReturnValue,
			GasUsed:     abi_spec.NewTokenAmount(128),
		}, receipts[0])

		td.AssertReceipt(types.MessageReceipt{
			ExitCode:    exitcode.Ok,
			ReturnValue: drivers.EmptyReturnValue,
			GasUsed:     abi_spec.NewTokenAmount(128),
		}, receipts[1])

		td.AssertBalance(sender, big_spec.Sub(senderInitialBal, senderTransferAmnt))
		td.AssertBalance(middle, big_spec.Sub(senderTransferAmnt, middleTransferAmnt))
		td.AssertBalance(final, middleTransferAmnt)
	})

}

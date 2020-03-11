package tipset

import (
	"context"
	"testing"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestInvalidSenderAddress(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("message sender address is defined but not in state tree", func(t *testing.T) {
		td := builder.Build(t)

		badSender := utils.NewSECP256K1Addr(t, "123")

		receiver, _ := td.NewAccountActor(drivers.SECP, abi_spec.NewTokenAmount(1_000_000))

		trnsferMsg := td.MessageProducer.Transfer(receiver, badSender, chain.Value(big_spec.NewInt(10)), chain.Nonce(0))

		blkMsgs := chain.NewTipSetMessageBuilder().
			WithMiner(td.ExeCtx.Miner).
			WithBLSMessage(trnsferMsg).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)

		td.AssertReceipt(receipts[0], types.MessageReceipt{
			ExitCode:    exitcode.SysErrActorNotFound,
			ReturnValue: drivers.EmptyReturnValue,
		})

		td.AssertBalance(receiver, abi_spec.NewTokenAmount(1_000_000))
	})
}

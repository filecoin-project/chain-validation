package suites

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	power_spec "github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestSomeStuff(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("idk dude", func(t *testing.T) {
		td := builder.Build(t)

		// TODO all this miner creation boiler plate should be moved somewhere else
		minerOwner, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(1_000_000_000))
		minerWorker, _ := td.NewAccountActor(address.BLS, big_spec.NewInt(0))
		ret := td.ComputeInitActorExecReturn(builtin_spec.StoragePowerActorAddr, 0, utils.NewIDAddr(t, 102))
		td.ApplyMessageExpectReceipt(
			td.MessageProducer.PowerCreateMiner(builtin_spec.StoragePowerActorAddr, minerOwner, power_spec.CreateMinerParams{
				Worker:     minerWorker,
				SectorSize: 1,
				Peer:       "peerid",
			}, chain.Value(big_spec.NewInt(1_000_000)), chain.Nonce(0)),
			types.MessageReceipt{
				ExitCode:    exitcode.Ok,
				ReturnValue: chain.MustSerialize(&ret),
				GasUsed:     big_spec.Zero(),
			},
		)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkMsgs := chain.NewTipSetMessageBuilder().
			WithMiner(ret.IDAddress).
			WithBLSMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.State(), []types.BlockMessagesInfo{blkMsgs}, td.ExeCtx.Epoch)
		require.NoError(t, err)
		require.Len(t, receipts, 1)

		require.Equal(t, exitcode.Ok, receipts[0].ExitCode)
		require.Equal(t, drivers.EmptyReturnValue, receipts[0].ReturnValue)

		td.AssertBalance(receiver, big_spec.NewInt(100))

	})

}

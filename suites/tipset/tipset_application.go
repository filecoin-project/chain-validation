package tipset

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
)

func TestBlockMessageInfoApplication(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("apply a single BLS message", func(t *testing.T) {
		td := builder.Build(t)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkMsgs := chain.NewTipSetMessageBuilder().
			// miner addresses are required to use ID protocol.
			WithMiner(td.ExeCtx.Miner).
			// send value from sender to receiver
			WithBLSMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)

		require.Equal(t, exitcode.Ok, receipts[0].ExitCode)
		require.Equal(t, drivers.EmptyReturnValue, receipts[0].ReturnValue)

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply a duplicated BLS message", func(t *testing.T) {
		td := builder.Build(t)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkMsgs := chain.NewTipSetMessageBuilder().
			// miner addresses are required to use ID protocol.
			WithMiner(td.ExeCtx.Miner).
			// duplicate the message
			WithBLSMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
			WithBLSMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		// despite there being 2 messages there is only one receipt.
		require.Len(t, receipts, 1)

		require.Equal(t, exitcode.Ok, receipts[0].ExitCode)
		require.Equal(t, drivers.EmptyReturnValue, receipts[0].ReturnValue)

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply a single SECP message", func(t *testing.T) {
		td := builder.Build(t)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkMsgs := chain.NewTipSetMessageBuilder().
			// miner addresses are required to use ID protocol.
			WithMiner(td.ExeCtx.Miner).
			// send value from sender to receiver
			WithSECPMessage(signMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))), td.Wallet())).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)

		require.Equal(t, exitcode.Ok, receipts[0].ExitCode)
		require.Equal(t, drivers.EmptyReturnValue, receipts[0].ReturnValue)

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply duplicate SECP message", func(t *testing.T) {
		td := builder.Build(t)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkMsgs := chain.NewTipSetMessageBuilder().
			// miner addresses are required to use ID protocol.
			WithMiner(td.ExeCtx.Miner).
			// send value from sender to receiver
			WithSECPMessage(signMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))), td.Wallet())).
			WithSECPMessage(signMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))), td.Wallet())).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)

		require.Equal(t, exitcode.Ok, receipts[0].ExitCode)
		require.Equal(t, drivers.EmptyReturnValue, receipts[0].ReturnValue)

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply duplicate BLS and SECP message", func(t *testing.T) {
		td := builder.Build(t)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkMsgs := chain.NewTipSetMessageBuilder().
			// miner addresses are required to use ID protocol.
			WithMiner(td.ExeCtx.Miner).
			// send value from sender to receiver
			WithSECPMessage(signMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))), td.Wallet())).
			WithBLSMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 1)

		require.Equal(t, exitcode.Ok, receipts[0].ExitCode)
		require.Equal(t, drivers.EmptyReturnValue, receipts[0].ReturnValue)

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

}

// TODO produce a valid signature
func signMessage(msg *types.Message, km state.KeyManager) *types.SignedMessage {
	return &types.SignedMessage{
		Message:   *msg,
		Signature: crypto.Signature{},
	}
}
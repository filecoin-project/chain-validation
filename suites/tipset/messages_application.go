package tipset

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
)

func TipSetTest_BlockMessageApplication(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState...)

	t.Run("SECP and BLS messages cost different amounts of gas", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		senderBLS, _ := td.NewAccountActor(address.BLS, big_spec.NewInt(10_000_000))
		receiverBLS, _ := td.NewAccountActor(address.BLS, big_spec.Zero())
		senderSECP, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiverSECP, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())
		transferAmnt := abi.NewTokenAmount(100)

		results := tipB.WithBlockBuilder(
			blkB.
				WithBLSMessageOk(
					td.MessageProducer.Transfer(receiverBLS, senderBLS, chain.Nonce(0), chain.Value(transferAmnt))).
				WithSECPMessageOk(
					td.MessageProducer.Transfer(receiverSECP, senderSECP, chain.Nonce(0), chain.Value(transferAmnt))),
		).ApplyAndValidate()

		require.Equal(t, 2, len(results.Receipts))

		blsGasUsed := int64(results.Receipts[0].GasUsed)
		secpGasUsed := int64(results.Receipts[1].GasUsed)
		assert.Greater(t, secpGasUsed, blsGasUsed)
	})

}

func TipSetTest_BlockMessageDeduplication(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState...)

	t.Run("apply a single BLS message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		tipB.WithBlockBuilder(
			// send value from sender to receiver
			blkB.WithBLSMessageOk(
				td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))),
			),
		).ApplyAndValidate()

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply a duplicated BLS message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		tipB.WithBlockBuilder(
			// duplicate the message
			blkB.WithBLSMessageOk(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
				// only should have a single result
				WithBLSMessageDropped(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))),
		).ApplyAndValidate()
		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply a single SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		tipB.WithBlockBuilder(
			// send value from sender to receiver
			blkB.WithSECPMessageOk(
				td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))),
			),
		).ApplyAndValidate()

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply duplicate SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		tipB.WithBlockBuilder(
			// send value from sender to receiver
			blkB.WithSECPMessageOk(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
				WithSECPMessageDropped(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))),
		).ApplyAndValidate()

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply duplicate BLS and SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		senderInitialBal := big_spec.NewInt(10_000_000)
		_, senderID := td.NewAccountActor(address.SECP256K1, senderInitialBal)
		_, receiverID := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		amountSent := big_spec.NewInt(100)
		result := tipB.WithBlockBuilder(
			// using ID addresses will ensure the BLS message and the unsigned message encapsulated in the SECP message
			// have the same CID.
			blkB.WithBLSMessageOk(td.MessageProducer.Transfer(receiverID, senderID, chain.Nonce(0), chain.Value(amountSent))).
				WithSECPMessageDropped(td.MessageProducer.Transfer(receiverID, senderID, chain.Nonce(0), chain.Value(amountSent))),
		).ApplyAndValidate()

		assert.Equal(t, 1, len(result.Receipts))

		td.AssertBalance(receiverID, amountSent)
		td.AssertBalance(senderID, big_spec.Sub(big_spec.Sub(senderInitialBal, amountSent), result.Receipts[0].GasUsed.Big()))
	})
}

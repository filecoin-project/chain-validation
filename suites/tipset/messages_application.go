package tipset

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	big_spec "github.com/filecoin-project/go-state-types/big"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
)

func TipSetTest_BlockMessageApplication(t *testing.T, factory state.Factories) {
	const gasLimit = 1_000_000_000
	const gasFeeCap = 200
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(gasLimit).
		WithDefaultGasFeeCap(gasFeeCap).
		WithDefaultGasPremium(1).
		WithActorState(drivers.DefaultBuiltinActorsState...)

	t.Run("SECP and BLS messages cost different amounts of gas", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		senderBLS, _ := td.NewAccountActor(address.BLS, big_spec.NewInt(10*gasFeeCap*gasLimit))
		receiverBLS, _ := td.NewAccountActor(address.BLS, big_spec.Zero())
		senderSECP, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10*gasFeeCap*gasLimit))
		receiverSECP, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())
		transferAmnt := abi.NewTokenAmount(100)

		results := tipB.WithBlockBuilder(
			blkB.
				WithBLSMessageOk(
					td.MessageProducer.Transfer(senderBLS, receiverBLS, chain.Nonce(0), chain.Value(transferAmnt))).
				WithSECPMessageOk(
					td.MessageProducer.Transfer(senderSECP, receiverSECP, chain.Nonce(0), chain.Value(transferAmnt))),
		).ApplyAndValidate()

		require.Equal(t, 2, len(results.Receipts))

		blsGasUsed := int64(results.Receipts[0].GasUsed)
		secpGasUsed := int64(results.Receipts[1].GasUsed)
		assert.Greater(t, secpGasUsed, blsGasUsed)
	})

}

func TipSetTest_BlockMessageDeduplication(t *testing.T, factory state.Factories) {
	const gasLimit = 1_000_000_000
	const gasFeeCap = 200
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(gasLimit).
		WithDefaultGasFeeCap(200).
		WithDefaultGasPremium(1).
		WithActorState(drivers.DefaultBuiltinActorsState...)

	t.Run("apply a single BLS message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10*gasFeeCap*gasLimit))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		tipB.WithBlockBuilder(
			// send value from sender to receiver
			blkB.WithBLSMessageOk(
				td.MessageProducer.Transfer(sender, receiver, chain.Nonce(0), chain.Value(big_spec.NewInt(100))),
			),
		).ApplyAndValidate()

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply a duplicated BLS message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10*gasFeeCap*gasLimit))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		tipB.WithBlockBuilder(
			// duplicate the message
			blkB.WithBLSMessageOk(td.MessageProducer.Transfer(sender, receiver, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
				// only should have a single result
				WithBLSMessageDropped(td.MessageProducer.Transfer(sender, receiver, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))),
		).ApplyAndValidate()
		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply a single SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10*gasFeeCap*gasLimit))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		tipB.WithBlockBuilder(
			// send value from sender to receiver
			blkB.WithSECPMessageOk(
				td.MessageProducer.Transfer(sender, receiver, chain.Nonce(0), chain.Value(big_spec.NewInt(100))),
			),
		).ApplyAndValidate()

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply duplicate SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10*gasFeeCap*gasLimit))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		tipB.WithBlockBuilder(
			// send value from sender to receiver
			blkB.WithSECPMessageOk(td.MessageProducer.Transfer(sender, receiver, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
				WithSECPMessageDropped(td.MessageProducer.Transfer(sender, receiver, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))),
		).ApplyAndValidate()

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply duplicate BLS and SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		tipB := drivers.NewTipSetMessageBuilder(td)
		blkB := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		senderInitialBal := big_spec.NewInt(10 * gasFeeCap * gasLimit)
		_, senderID := td.NewAccountActor(address.SECP256K1, senderInitialBal)
		_, receiverID := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		amountSent := big_spec.NewInt(100)
		msgOriginal := td.MessageProducer.Transfer(senderID, receiverID, chain.Nonce(0), chain.Value(amountSent))
		msgDup := td.MessageProducer.Transfer(senderID, receiverID, chain.Nonce(0), chain.Value(amountSent))
		result := tipB.WithBlockBuilder(
			// using ID addresses will ensure the BLS message and the unsigned message encapsulated in the SECP message
			// have the same CID.
			blkB.WithBLSMessageOk(msgOriginal).WithSECPMessageDropped(msgDup),
		).ApplyAndValidate()

		assert.Equal(t, 1, len(result.Receipts))

		td.AssertBalance(receiverID, amountSent)
		td.AssertActorChange(senderID, senderInitialBal, msgOriginal.GasLimit, msgOriginal.GasPremium, msgOriginal.Value, result.Receipts[0], msgOriginal.CallSeqNum+1)
	})
}

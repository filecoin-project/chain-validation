package tipset

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
)

func TestBlockMessageDeduplication(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

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

}

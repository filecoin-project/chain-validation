package tipset

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
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
		blkBuilder := drivers.NewTipSetMessageBuilder(td)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkBuilder.WithTicketCount(1).
			// send value from sender to receiver
			WithBLSMessageAndReceipt(
				td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))),
				types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: abi_spec.NewTokenAmount(128)},
			).ApplyAndValidate()

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply a duplicated BLS message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		blkBuilder := drivers.NewTipSetMessageBuilder(td)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkBuilder.WithTicketCount(1).
			// duplicate the message
			WithBLSMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
			WithBLSMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100)))).
			WithMessageReceipt(types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: nil, GasUsed: abi_spec.NewTokenAmount(128)}).
			ApplyAndValidate()
		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply a single SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		blkBuilder := drivers.NewTipSetMessageBuilder(td)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkBuilder.WithTicketCount(1).
			// send value from sender to receiver
			WithSECPMessageAndReceipt(
				signMessage(
					td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))),
					td.Wallet()),
				types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: abi_spec.NewTokenAmount(128)},
			).ApplyAndValidate()

		td.AssertBalance(receiver, big_spec.NewInt(100))
	})

	t.Run("apply duplicate SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		blkBuilder := drivers.NewTipSetMessageBuilder(td)

		sender, _ := td.NewAccountActor(address.SECP256K1, big_spec.NewInt(10_000_000))
		receiver, _ := td.NewAccountActor(address.SECP256K1, big_spec.Zero())

		blkBuilder.WithTicketCount(1).
			// send value from sender to receiver
			WithSECPMessage(signMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))), td.Wallet())).
			WithSECPMessage(signMessage(td.MessageProducer.Transfer(receiver, sender, chain.Nonce(0), chain.Value(big_spec.NewInt(100))), td.Wallet())).
			WithMessageReceipt(types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: nil, GasUsed: abi_spec.NewTokenAmount(128)}).
			ApplyAndValidate()

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

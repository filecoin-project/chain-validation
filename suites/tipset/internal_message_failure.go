package tipset

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	init_ "github.com/filecoin-project/specs-actors/actors/builtin/init"
	"github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	typegen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestInternalMessageApplicationFailure(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("multisig internal message send fails and receiver of message does not exist in state", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		blkBuilder := drivers.NewTipSetMessageBuilder(td)

		alice, aliceId := td.NewAccountActor(drivers.SECP, abi.NewTokenAmount(2_000_000_000_000))

		multisigAddr := utils.NewIDAddr(t, utils.IdFromAddress(aliceId)+1)
		createRet := td.ComputeInitActorExecReturn(alice, 0, 0, multisigAddr)
		txID0 := typegen.CborInt(multisig.TxnID(0))

		target := utils.NewSECP256K1Addr(t, "pubbub")
		proposeParams := multisig.ProposeParams{To: target, Value: big.Zero(), Method: 99, Params: nil} // Note invalid method ID

		// Create the multisig actor and propose the send
		blkBuilder.WithTicketCount(1).
			WithSECPMessageAndReceipt(
				signMessage(td.MessageProducer.CreateMultisigActor(aliceId, []address.Address{aliceId}, 0, 1, chain.Nonce(0)), td.Wallet()),
				types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: abi.NewTokenAmount(1282)},
			).
			WithSECPMessageAndReceipt(
				signMessage(td.MessageProducer.MultisigPropose(multisigAddr, aliceId, proposeParams, chain.Nonce(1)), td.Wallet()),
				types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&txID0), GasUsed: abi.NewTokenAmount(1235)},
			).
			ApplyAndValidate()

		// assert we get an error when attempting to get the actor that should not exist
		_, err := td.State().Actor(target)
		require.Error(t, err)
		require.Contains(t, err.Error(), init_.ErrAddressNotFound.Error())
	})

	t.Run("multisig internal message send fails when receiver is invalid address protocol", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		blkBuilder := drivers.NewTipSetMessageBuilder(td)

		// create the multisig actor, set number of approvals to 1 so propose goes through on first send.
		alice, aliceId := td.NewAccountActor(drivers.SECP, abi.NewTokenAmount(2_000_000_000_000))
		multisigAddr := utils.NewIDAddr(t, utils.IdFromAddress(aliceId)+1)

		invalidAddress := utils.NewSECP256K1Addr(t, "pubbbub")
		invalidAddress.Bytes()[0] = address.Unknown // Exploit imperfect immutability of address to mutate its internals.
		params := multisig.ProposeParams{To: invalidAddress, Value: big.Zero(), Method: 0, Params: nil}

		createRet := td.ComputeInitActorExecReturn(alice, 0, 0, multisigAddr)
		blkBuilder.WithTicketCount(1).
			WithSECPMessageAndReceipt(
				signMessage(
					td.MessageProducer.CreateMultisigActor(aliceId, []address.Address{aliceId}, 0, 1, chain.Nonce(0)),
					td.Wallet()), types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: abi.NewTokenAmount(1282)},
			).
			WithSECPMessageAndReceipt(
				signMessage(
					td.MessageProducer.MultisigPropose(multisigAddr, aliceId, params, chain.Nonce(1)), td.Wallet()),
				types.MessageReceipt{ExitCode: exitcode.SysErrInvalidParameters, ReturnValue: drivers.EmptyReturnValue, GasUsed: abi.NewTokenAmount(1235)},
			).
			ApplyAndValidate()

		// the multisig txnid should increment and the 0th transaction should have been removed
		var msa multisig.State
		td.GetActorState(multisigAddr, &msa)
		td.AssertMultisigContainsTransaction(multisigAddr, 0, false)
		assert.Equal(t, multisig.TxnID(1), msa.NextTxnID)
	})
}

package tipset

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	init_spec "github.com/filecoin-project/specs-actors/actors/builtin/init"

	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func TestInternalMessageApplicationFailure(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("multisig internal message send fails and receiver of message does not exist in state", func(t *testing.T) {
		td := builder.Build(t)
		blkBuilder := drivers.NewTipSetMessageBuilder(td)

		_, aliceId := td.NewAccountActor(drivers.SECP, abi_spec.NewTokenAmount(2_000_000_000_000))

		// this address should not exist since the propose is sent to the wrong account actor method
		nobody := utils.NewSECP256K1Addr(t, "pubbub")

		multisigAddr := utils.NewIDAddr(t, utils.IdFromAddress(aliceId)+1)

		createRet := td.ComputeInitActorExecReturn(aliceId, 0, multisigAddr)
		// Create the multisig actor and propose the send
		blkBuilder.WithTicketCount(1).
			WithSECPMessageAndReceipt(
				signMessage(
					td.MessageProducer.CreateMultisigActor(aliceId,
						[]address.Address{aliceId}, 10, 1,
						chain.Nonce(0), chain.Value(big_spec.Zero())),
					td.Wallet()),
				types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: abi_spec.NewTokenAmount(1282)},
			).
			WithSECPMessageAndReceipt(
				signMessage(
					td.MessageProducer.MultisigPropose(multisigAddr, aliceId,
						multisig_spec.ProposeParams{To: nobody, Value: big_spec.Zero(), Method: 99, Params: nil},
						chain.Nonce(1), chain.Value(big_spec.Zero())),
					td.Wallet()),
				types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: abi_spec.NewTokenAmount(1235)},
			).
			ApplyAndValidate()

		// assert we get an error when attempting to get the actor that should not exist
		_, err := td.State().Actor(nobody)
		require.Error(t, err)
		require.Contains(t, err.Error(), init_spec.ErrAddressNotFound.Error())
	})

	t.Run("multisig internal message send fails when receiver is invalid address protocol", func(t *testing.T) {
		td := builder.Build(t)
		blkBuilder := drivers.NewTipSetMessageBuilder(td)

		// create the multisig actor, set number of approvals to 1 so propose goes through on first send.
		_, aliceId := td.NewAccountActor(drivers.SECP, abi_spec.NewTokenAmount(2_000_000_000_000))
		multisigAddr := utils.NewIDAddr(t, utils.IdFromAddress(aliceId)+1)

		// this address will be changed to protocol Unknown after its been serialized into bytes for a parameter.
		nobody := utils.NewSECP256K1Addr(t, "pubbbub")
		params := multisig_spec.ProposeParams{
			To:     nobody,
			Value:  big_spec.Zero(),
			Method: 0,
			Params: nil,
		}
		proposeParams := chain.MustSerialize(&params)
		proposeParams[2] = address.Unknown // the 3rd byte in the slice is the address protocol identifier, set to invalid protocol

		createRet := td.ComputeInitActorExecReturn(aliceId, 0, multisigAddr)
		blkBuilder.WithTicketCount(1).
			WithSECPMessageAndReceipt(
				signMessage(
					td.MessageProducer.CreateMultisigActor(aliceId,
						[]address.Address{aliceId}, 10, 1,
						chain.Nonce(0), chain.Value(big_spec.Zero()),
					),
					td.Wallet()), types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: abi_spec.NewTokenAmount(1282)},
			).
			WithSECPMessageAndReceipt(
				signMessage(
					td.MessageProducer.Build(multisigAddr, aliceId, builtin_spec.MethodsMultisig.Propose, proposeParams, chain.Nonce(1), chain.Value(big_spec.Zero())),
					td.Wallet()),
				types.MessageReceipt{ExitCode: exitcode.SysErrInvalidParameters, ReturnValue: drivers.EmptyReturnValue, GasUsed: abi_spec.NewTokenAmount(1235)},
			).
			ApplyAndValidate()

		// the multisig txnid should increment and the 0th transaction should have been removed
		var msa multisig_spec.State
		td.GetActorState(multisigAddr, &msa)
		td.AssertMultisigContainsTransaction(multisigAddr, 0, false)
		assert.Equal(t, multisig_spec.TxnID(1), msa.NextTxnID)
	})
}

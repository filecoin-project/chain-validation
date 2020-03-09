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

	t.Run("multisig internal message send fails and receiver of message is not created", func(t *testing.T) {
		td := builder.Build(t)

		_, aliceId := td.NewAccountActor(drivers.SECP, abi_spec.NewTokenAmount(2_000_000_000_000))

		// this address should not exist since the propose is sent to the wrong account actor method
		nobody := utils.NewSECP256K1Addr(t, "pubbub")

		multisigAddr := utils.NewIDAddr(t, utils.IdFromAddress(aliceId)+1)

		createMsMsg := td.MessageProducer.CreateMultisigActor(aliceId, []address.Address{aliceId}, 10, 1,
			chain.Nonce(0), chain.Value(big_spec.Zero()))

		proposeMsMsg := td.MessageProducer.MultisigPropose(multisigAddr, aliceId, multisig_spec.ProposeParams{
			To:     nobody,
			Value:  big_spec.Zero(),
			Method: 99, // this will cause the internal send to fail as this is not a method on an account actor
			Params: nil,
		}, chain.Nonce(1), chain.Value(big_spec.Zero()))

		// Create the multisig actor and propose the send
		blkMsgs := chain.NewTipSetMessageBuilder().
			WithMiner(td.ExeCtx.Miner).
			WithSECPMessage(signMessage(createMsMsg, td.Wallet())).
			WithSECPMessage(signMessage(proposeMsMsg, td.Wallet())).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 2)

		createMsRet := receipts[0]
		proposeMsRet := receipts[1]

		createRet := td.ComputeInitActorExecReturn(aliceId, 0, multisigAddr)
		td.AssertReceipt(createMsRet, types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: abi_spec.NewTokenAmount(1282)})
		td.AssertReceipt(proposeMsRet, types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: drivers.EmptyReturnValue, GasUsed: abi_spec.NewTokenAmount(1235)})

		var msa multisig_spec.State
		td.GetActorState(multisigAddr, &msa)

		// assert we get an error when attempting to get the actor that should not exist
		_, err = td.State().Actor(nobody)
		require.Error(t, err)
		require.Contains(t, err.Error(), init_spec.ErrAddressNotFound.Error())
	})

	t.Run("multisig internal message send fails when receiver is invalid address protocol", func(t *testing.T) {
		td := builder.Build(t)

		// create the multisig actor, set number of approvals to 1 so propose goes through on first send.
		_, aliceId := td.NewAccountActor(drivers.SECP, abi_spec.NewTokenAmount(2_000_000_000_000))
		multisigAddr := utils.NewIDAddr(t, utils.IdFromAddress(aliceId)+1)
		createMsMsg := td.MessageProducer.CreateMultisigActor(aliceId, []address.Address{aliceId}, 10, 1,
			chain.Nonce(0), chain.Value(big_spec.Zero()))

		// this address will be changed to protocol Unknown after its been serialized into bytes for a parameter.
		nobody := utils.NewSECP256K1Addr(t, "pubbbub")
		params := multisig_spec.ProposeParams{
			To:     nobody,
			Value:  big_spec.Zero(),
			Method: 0,
			Params: nil,
		}
		ser := chain.MustSerialize(&params)
		ser[2] = address.Unknown // the 3rd byte in the slice is the address protocol identifier, set to invalid protocol
		proposeMsMsg := td.MessageProducer.Build(multisigAddr, aliceId, builtin_spec.MethodsMultisig.Propose, ser, chain.Nonce(1), chain.Value(big_spec.Zero()))

		// Create the multisig actor and propose the send
		blkMsgs := chain.NewTipSetMessageBuilder().
			WithMiner(td.ExeCtx.Miner).
			WithSECPMessage(signMessage(createMsMsg, td.Wallet())).
			WithSECPMessage(signMessage(proposeMsMsg, td.Wallet())).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 2)

		createMsRet := receipts[0]
		proposeMsRet := receipts[1]

		createRet := td.ComputeInitActorExecReturn(aliceId, 0, multisigAddr)
		td.AssertReceipt(createMsRet, types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: abi_spec.NewTokenAmount(1282)})
		td.AssertReceipt(proposeMsRet, types.MessageReceipt{ExitCode: exitcode.SysErrInvalidParameters, ReturnValue: drivers.EmptyReturnValue, GasUsed: abi_spec.NewTokenAmount(1235)})

		// the multisig txnid should increment and the 0th transaction should have been removed
		var msa multisig_spec.State
		td.GetActorState(multisigAddr, &msa)
		td.AssertMultisigContainsTransaction(multisigAddr, 0, false)
		assert.Equal(t, multisig_spec.TxnID(1), msa.NextTxnID)
	})

	t.Run("multisig internal message send fails when receiver is undefined address", func(t *testing.T) {
		td := builder.Build(t)

		// create the multisig actor, set number of approvals to 1 so propose goes through on first send.
		_, aliceId := td.NewAccountActor(drivers.SECP, abi_spec.NewTokenAmount(2_000_000_000_000))
		multisigAddr := utils.NewIDAddr(t, utils.IdFromAddress(aliceId)+1)
		createMsMsg := td.MessageProducer.CreateMultisigActor(aliceId, []address.Address{aliceId}, 10, 1,
			chain.Nonce(0), chain.Value(big_spec.Zero()))

		// this address is the receiver with value undefined
		nobody := address.Undef
		params := multisig_spec.ProposeParams{
			To:     nobody,
			Value:  big_spec.Zero(),
			Method: 0,
			Params: nil,
		}
		ser := chain.MustSerialize(&params)
		ser[2] = address.Unknown // the 3rd byte in the slice is the address protocol identifier, set to invalid protocol
		proposeMsMsg := td.MessageProducer.Build(multisigAddr, aliceId, builtin_spec.MethodsMultisig.Propose, ser, chain.Nonce(1), chain.Value(big_spec.Zero()))

		// Create the multisig actor and propose the send
		blkMsgs := chain.NewTipSetMessageBuilder().
			WithMiner(td.ExeCtx.Miner).
			WithSECPMessage(signMessage(createMsMsg, td.Wallet())).
			WithSECPMessage(signMessage(proposeMsMsg, td.Wallet())).
			Build()

		receipts, err := td.Validator.ApplyTipSetMessages(td.ExeCtx, td.State(), []types.BlockMessagesInfo{blkMsgs}, td.Randomness())
		require.NoError(t, err)
		require.Len(t, receipts, 2)

		createMsRet := receipts[0]
		proposeMsRet := receipts[1]

		createRet := td.ComputeInitActorExecReturn(aliceId, 0, multisigAddr)
		td.AssertReceipt(createMsRet, types.MessageReceipt{ExitCode: exitcode.Ok, ReturnValue: chain.MustSerialize(&createRet), GasUsed: abi_spec.NewTokenAmount(1282)})
		td.AssertReceipt(proposeMsRet, types.MessageReceipt{ExitCode: exitcode.SysErrInvalidParameters, ReturnValue: drivers.EmptyReturnValue, GasUsed: abi_spec.NewTokenAmount(1235)})

		// the multisig txnid should increment and the 0th transaction should have been removed
		var msa multisig_spec.State
		td.GetActorState(multisigAddr, &msa)
		td.AssertMultisigContainsTransaction(multisigAddr, 0, false)
		assert.Equal(t, multisig_spec.TxnID(1), msa.NextTxnID)
	})
}

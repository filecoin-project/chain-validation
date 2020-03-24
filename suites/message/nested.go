package message

import (
	"bytes"
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin "github.com/filecoin-project/specs-actors/actors/builtin"
	init_ "github.com/filecoin-project/specs-actors/actors/builtin/init"
	"github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	"github.com/filecoin-project/specs-actors/actors/runtime"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	typegen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

// Tests exercising messages sent internally from one actor to another.
// These use a multisig actor with approvers=1 as a convenient staging ground for arbitrary internal messages.
func TestNestedSends(t *testing.T, factory state.Factories) {
	var acctDefaultBalance = abi.NewTokenAmount(1_000_000_000)
	var multisigBalance = abi.NewTokenAmount(1_000_000)

	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	t.Run("ok basic", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		stage := prepareStage(td, acctDefaultBalance, multisigBalance)
		balanceBefore := td.GetBalance(stage.creator)

		// Multisig sends back to the creator.
		amtSent := abi.NewTokenAmount(1)
		result := stage.send(stage.creator, amtSent, builtin.MethodSend, nil, 1)
		assert.Equal(t, exitcode.Ok, result.Receipt.ExitCode)

		td.AssertActor(stage.creator, big.Sub(big.Add(balanceBefore, amtSent), result.Receipt.GasUsed.Big()), 2)
	})

	t.Run("ok to new actor", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		stage := prepareStage(td, acctDefaultBalance, multisigBalance)
		balanceBefore := td.GetBalance(stage.creator)

		// Multisig sends to new address.
		newAddr := td.Wallet().NewSECP256k1AccountAddress()
		amtSent := abi.NewTokenAmount(1)
		result := stage.send(newAddr, amtSent, builtin.MethodSend, nil, 1)
		assert.Equal(t, exitcode.Ok, result.Receipt.ExitCode)

		td.AssertBalance(stage.msAddr, big.Sub(multisigBalance, amtSent))
		td.AssertBalance(stage.creator, big.Sub(balanceBefore, result.Receipt.GasUsed.Big()))
		td.AssertBalance(newAddr, amtSent)
	})

	t.Run("ok to new actor with invoke", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		stage := prepareStage(td, acctDefaultBalance, multisigBalance)
		balanceBefore := td.GetBalance(stage.creator)

		// Multisig sends to new address and invokes pubkey method at the same time.
		newAddr := td.Wallet().NewSECP256k1AccountAddress()
		amtSent := abi.NewTokenAmount(1)
		result := stage.send(newAddr, amtSent, builtin.MethodsAccount.PubkeyAddress, nil, 1)
		assert.Equal(t, exitcode.Ok, result.Receipt.ExitCode)
		// TODO: use an explicit Approve() and check the return value is the correct pubkey address
		// when the multisig Approve() method plumbs through the inner exit code and value.
		// https://github.com/filecoin-project/specs-actors/issues/113
		//expected := bytes.Buffer{}
		//require.NoError(t, newAddr.MarshalCBOR(&expected))
		//assert.Equal(t, expected.Bytes(), result.Receipt.ReturnValue)

		td.AssertBalance(stage.msAddr, big.Sub(multisigBalance, amtSent))
		td.AssertBalance(stage.creator, big.Sub(balanceBefore, result.Receipt.GasUsed.Big()))
		td.AssertBalance(newAddr, amtSent)
	})

	t.Run("ok recursive", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		stage := prepareStage(td, acctDefaultBalance, multisigBalance)
		_, anotherId := td.NewAccountActor(drivers.SECP, big.Zero())
		balanceBefore := td.GetBalance(stage.creator)

		// Multisig sends to itself.
		params := multisig.AddSignerParams{
			Signer:   anotherId,
			Increase: false,
		}
		result := stage.send(stage.msAddr, big.Zero(), builtin.MethodsMultisig.AddSigner, &params, 1)
		assert.Equal(t, exitcode.Ok, result.Receipt.ExitCode)

		td.AssertBalance(stage.msAddr, multisigBalance)
		assert.Equal(t, big.Sub(balanceBefore, result.Receipt.GasUsed.Big()), td.GetBalance(stage.creator))
		var st multisig.State
		td.GetActorState(stage.msAddr, &st)
		assert.Equal(t, []address.Address{stage.creator, anotherId}, st.Signers)
	})

	t.Run("ok non-CBOR params with transfer", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		stage := prepareStage(td, acctDefaultBalance, multisigBalance)

		newAddr := td.Wallet().NewSECP256k1AccountAddress()
		amtSent := abi.NewTokenAmount(1)
		// So long as the parameters are not actually used by the method, a message can carry arbitrary bytes.
		params := typegen.Deferred{Raw: []byte{1, 2, 3, 4}}
		result := stage.send(newAddr, amtSent, builtin.MethodSend, &params, 1)
		assert.Equal(t, exitcode.Ok, result.Receipt.ExitCode)

		td.AssertBalance(stage.msAddr, big.Sub(multisigBalance, amtSent))
		td.AssertBalance(newAddr, amtSent)
	})

	//
	// TODO: Tests to exercise invalid "syntax" of the inner message.
	// These would fail message syntax validation if the message were top-level.
	//
	// Some of these require handcrafting the proposal params serialization.
	// - malformed address: zero-length, one-length, too-short pubkeys, invalid UVarints, ...
	// - negative method num
	//
	// Unfortunately the multisig actor can't be used to trigger a negative-value internal transfer because
	// it checks just before sending.
	// We need a custom actor for staging whackier messages.

	//
	// The following tests exercise invalid semantics of the inner message
	//

	t.Run("fail nonexistent ID address", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		stage := prepareStage(td, acctDefaultBalance, multisigBalance)

		newAddr := utils.NewIDAddr(t, 1234)
		amtSent := abi.NewTokenAmount(1)
		result := stage.send(newAddr, amtSent, builtin.MethodSend, nil, 1)
		assert.Equal(t, exitcode.Ok, result.Receipt.ExitCode)

		td.AssertBalance(stage.msAddr, multisigBalance) // No change.
		_, err := td.State().Actor(newAddr)
		assert.Error(t, err)
	})

	t.Run("fail nonexistent actor address", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		stage := prepareStage(td, acctDefaultBalance, multisigBalance)

		newAddr := utils.NewActorAddr(t, "1234")
		amtSent := abi.NewTokenAmount(1)
		result := stage.send(newAddr, amtSent, builtin.MethodSend, nil, 1)
		assert.Equal(t, exitcode.Ok, result.Receipt.ExitCode)

		td.AssertBalance(stage.msAddr, multisigBalance) // No change.
		_, err := td.State().Actor(newAddr)
		assert.Error(t, err)
	})

	t.Run("fail invalid methodnum new actor", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		stage := prepareStage(td, acctDefaultBalance, multisigBalance)

		newAddr := td.Wallet().NewSECP256k1AccountAddress()
		amtSent := abi.NewTokenAmount(1)
		result := stage.send(newAddr, amtSent, abi.MethodNum(99), nil, 1)
		assert.Equal(t, exitcode.Ok, result.Receipt.ExitCode)

		td.AssertBalance(stage.msAddr, multisigBalance) // No change.
		_, err := td.State().Actor(newAddr)
		assert.Error(t, err)
	})

	t.Run("fail invalid methodnum for actor", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		stage := prepareStage(td, acctDefaultBalance, multisigBalance)
		balanceBefore := td.GetBalance(stage.creator)

		amtSent := abi.NewTokenAmount(1)
		result := stage.send(stage.creator, amtSent, abi.MethodNum(99), nil, 1)
		assert.Equal(t, exitcode.Ok, result.Receipt.ExitCode)

		td.AssertBalance(stage.msAddr, multisigBalance)                                       // No change.
		td.AssertBalance(stage.creator, big.Sub(balanceBefore, result.Receipt.GasUsed.Big())) // Pay gas, don't receive funds.
	})

	// TODO more tests:
	// fail send empty params to method with params
	// fail send wrong shape params to method
	// fail send running out of gas on inner method
	// fail send when target method aborts
	// fail send when target method on multisig (recursive) aborts
}

// Wraps a multisig actor as a stage for nested sends.
type ms_stage struct {
	driver  *drivers.TestDriver
	creator address.Address // Address of the creator and sole signer of the multisig.
	msAddr  address.Address // Address of the multisig actor from which nested messages are sent.
}

// Creates a multisig actor with its creator as sole approver.
func prepareStage(td *drivers.TestDriver, creatorBalance, msBalance abi.TokenAmount) *ms_stage {
	_, creatorId := td.NewAccountActor(drivers.SECP, creatorBalance)

	msg := td.MessageProducer.CreateMultisigActor(creatorId, []address.Address{creatorId}, 0, 1, chain.Value(msBalance), chain.Nonce(0))
	result := td.ApplyMessage(msg)
	require.Equal(td.T, exitcode.Ok, result.Receipt.ExitCode)
	var ret init_.ExecReturn
	err := ret.UnmarshalCBOR(bytes.NewReader(result.Receipt.ReturnValue))
	require.NoError(td.T, err)

	return &ms_stage{
		driver:  td,
		creator: creatorId,
		msAddr:  ret.IDAddress,
	}
}

func (s *ms_stage) send(to address.Address, value abi.TokenAmount, method abi.MethodNum, params runtime.CBORMarshaler, approverNonce uint64) chain.ApplyResult {
	buf := bytes.Buffer{}
	if params != nil {
		err := params.MarshalCBOR(&buf)
		require.NoError(s.driver.T, err)
	}
	pparams := multisig.ProposeParams{
		To:     to,
		Value:  value,
		Method: method,
		Params: buf.Bytes(),
	}
	msg := s.driver.MessageProducer.MultisigPropose(s.msAddr, s.creator, pparams, chain.Nonce(approverNonce))
	return s.driver.ApplyMessage(msg)
}

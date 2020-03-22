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

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
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

	t.Run("ok basic send", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		_, aliceId := td.NewAccountActor(drivers.SECP, acctDefaultBalance)

		stage := prepareStage(td, aliceId, multisigBalance, 0)
		balanceBefore := td.GetBalance(aliceId)

		amtSent := abi.NewTokenAmount(1)
		recpt := stage.send(aliceId, amtSent, builtin.MethodSend, nil, 1)
		assert.Equal(t, exitcode.Ok, recpt.ExitCode)
		balanceAfter := td.GetBalance(aliceId)

		assert.Equal(t, big.Sub(big.Add(balanceBefore, amtSent), recpt.GasUsed), balanceAfter)
	})

	// TODO more tests:
	// ok send to new actor
	// ok send to self (the multisig)
	// fail send to invalid address
	// fail send negative value
	// fail send running out of gas on inner method
	// fail send to non-existent {ID|actor} address
	// fail send empty params to method with params
	// fail send non-CBOR params
	// fail send wrong shape params to method
	// fail send {existing|new} actor with invalid method id (replace TestInternalMessageApplicationFailure)
	// fail send when target method aborts
}

// Wraps a multisig actor as a stage for nested sends.
type ms_stage struct {
	driver   *drivers.TestDriver
	multisig address.Address
	signer   address.Address
}

// Creates a multisig actor with its creator as sole approver.
func prepareStage(td *drivers.TestDriver, creatorId address.Address, value abi.TokenAmount, creatorNonce int64) *ms_stage {
	msg := td.MessageProducer.CreateMultisigActor(creatorId, []address.Address{creatorId}, 0, 1, chain.Value(value), chain.Nonce(creatorNonce))
	receipt := td.ApplyMessage(msg)
	require.Equal(td.T, exitcode.Ok, receipt.ExitCode)
	var ret init_.ExecReturn
	err := ret.UnmarshalCBOR(bytes.NewReader(receipt.ReturnValue))
	require.NoError(td.T, err)

	return &ms_stage{
		driver:   td,
		multisig: ret.IDAddress,
		signer:   creatorId,
	}
}

func (s *ms_stage) send(to address.Address, value abi.TokenAmount, method abi.MethodNum, params runtime.CBORMarshaler, approverNonce int64) types.MessageReceipt {
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
	msg := s.driver.MessageProducer.MultisigPropose(s.multisig, s.signer, pparams, chain.Nonce(approverNonce))
	return s.driver.ApplyMessage(msg)
}

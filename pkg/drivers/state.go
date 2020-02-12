package drivers

import (
	"context"
	"fmt"
	"testing"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
)

var (
	SECP = address.SECP256K1
	BLS  = address.BLS
)

// Factories wraps up all the implementation-specific integration points.
type Factories interface {
	NewState() state.Wrapper

	chain.Applier
}

// StateDriver mutates and inspects a state.
type StateDriver struct {
	tb testing.TB
	st state.Wrapper
}

// NewStateDriver creates a new state driver for a state.
func NewStateDriver(tb testing.TB, w state.Wrapper) *StateDriver {
	return &StateDriver{tb, w}
}

// State returns the state.
func (d *StateDriver) State() state.Wrapper {
	return d.st
}

// NewAccountActor installs a new account actor, returning the address.
func (d *StateDriver) NewAccountActor(addrType address.Protocol, balanceAttoFil abi_spec.TokenAmount) address.Address {
	var addr address.Address
	switch addrType {
	case address.SECP256K1:
		addr = d.st.NewSecp256k1AccountAddress()
	case address.BLS:
		addr = d.st.NewBLSAccountAddress()
	default:
		require.FailNowf(d.tb, "unsupported address", "protocol for account actor: %v", addrType)
	}

	_, _, err := d.st.SetActor(addr, builtin_spec.AccountActorCodeID, balanceAttoFil)
	require.NoError(d.tb, err)
	return addr
}

// AssertBalance checks an actor has an expected balance.
func (d *StateDriver) AssertBalance(addr address.Address, expected big_spec.Int) {
	actr, err := d.st.Actor(addr)
	require.NoError(d.tb, err)
	assert.Equal(d.tb, expected, actr.Balance(), fmt.Sprintf("expected balance: %v, actual balance: %v", expected, actr.Balance().String()))
}

// AssertReceipt checks that a receipt is not nill and has values equal to `expected`.
func (d *StateDriver) AssertReceipt(receipt, expected chain.MessageReceipt) {
	assert.NotNil(d.tb, receipt)
	assert.Equal(d.tb, expected.GasUsed, receipt.GasUsed, fmt.Sprintf("expected gas: %v, actual gas: %v", expected.GasUsed, receipt.GasUsed))
	assert.Equal(d.tb, expected.ReturnValue, receipt.ReturnValue, fmt.Sprintf("expected return value: %v, actual return value: %v", expected.ReturnValue, receipt.ReturnValue))
	assert.Equal(d.tb, expected.ExitCode, receipt.ExitCode, fmt.Sprintf("expected exit code: %v, actual exit code: %v", expected.ExitCode, receipt.ExitCode))
}

func (d *StateDriver) AssertMultisigTransaction(multisigAddr address.Address, txnID multisig_spec.TxnID, txn multisig_spec.Transaction) {
	multisigActor, err := d.State().Actor(multisigAddr)
	require.NoError(d.tb, err)

	strg, err := d.State().Storage()
	require.NoError(d.tb, err)

	var multisig multisig_spec.State
	require.NoError(d.tb, strg.Get(context.Background(), multisigActor.Head(), &multisig))

	txnMap := adt_spec.AsMap(strg, multisig.PendingTxns)
	var actualTxn multisig_spec.Transaction
	found, err := txnMap.Get(txnID, &actualTxn)
	assert.NoError(d.tb, err)
	assert.True(d.tb, found)

	assert.Equal(d.tb, txn, actualTxn)
}

func (d *StateDriver) AssertMultisigContainsTransaction(multisigAddr address.Address, txnID multisig_spec.TxnID, contains bool) {
	multisigActor, err := d.State().Actor(multisigAddr)
	require.NoError(d.tb, err)

	strg, err := d.State().Storage()
	require.NoError(d.tb, err)

	var multisig multisig_spec.State
	require.NoError(d.tb, strg.Get(context.Background(), multisigActor.Head(), &multisig))

	txnMap := adt_spec.AsMap(strg, multisig.PendingTxns)
	var actualTxn multisig_spec.Transaction
	found, err := txnMap.Get(txnID, &actualTxn)
	require.NoError(d.tb, err)
	assert.Equal(d.tb, contains, found)

}

func (d *StateDriver) AssertMultisigState(multisigAddr address.Address, expected multisig_spec.State) {
	multisigActor, err := d.State().Actor(multisigAddr)
	require.NoError(d.tb, err)

	strg, err := d.State().Storage()
	require.NoError(d.tb, err)

	var multisig multisig_spec.State
	require.NoError(d.tb, strg.Get(context.Background(), multisigActor.Head(), &multisig))

	assert.NotNil(d.tb, multisig)
	assert.Equal(d.tb, expected.InitialBalance, multisig.InitialBalance, fmt.Sprintf("expected InitialBalance: %v, actual InitialBalance: %v", expected.InitialBalance, multisig.InitialBalance))
	assert.Equal(d.tb, expected.NextTxnID, multisig.NextTxnID, fmt.Sprintf("expected NextTxnID: %v, actual NextTxnID: %v", expected.NextTxnID, multisig.NextTxnID))
	assert.Equal(d.tb, expected.NumApprovalsThreshold, multisig.NumApprovalsThreshold, fmt.Sprintf("expected NumApprovalsThreshold: %v, actual NumApprovalsThreshold: %v", expected.NumApprovalsThreshold, multisig.NumApprovalsThreshold))
	assert.Equal(d.tb, expected.StartEpoch, multisig.StartEpoch, fmt.Sprintf("expected StartEpoch: %v, actual StartEpoch: %v", expected.StartEpoch, multisig.StartEpoch))
	assert.Equal(d.tb, expected.UnlockDuration, multisig.UnlockDuration, fmt.Sprintf("expected UnlockDuration: %v, actual UnlockDuration: %v", expected.UnlockDuration, multisig.UnlockDuration))

	for _, e := range expected.Signers {
		assert.Contains(d.tb, multisig.Signers, e, fmt.Sprintf("expected Signer: %v, actual Signer: %v", e, multisig.Signers))
	}
}

package suites

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
)

// Factories wraps up all the implementation-specific integration points.
type Factories interface {
	NewState() state.Wrapper
	NewMessageFactory(wrapper state.Wrapper) chain.MessageFactory

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
func (d *StateDriver) NewAccountActor(balanceAttoFil uint64) state.Address {
	addr, err := d.st.NewAccountAddress()
	require.NoError(d.tb, err)

	_, _, err = d.st.SetActor(addr, state.AccountActorCodeCid, af(balanceAttoFil))
	require.NoError(d.tb, err)
	return addr
}

// AssertBalance checks an actor has an expected balance.
func (d *StateDriver) AssertBalance(addr state.Address, expected uint64) {
	actr, err := d.st.Actor(addr)
	require.NoError(d.tb, err)
	assert.Equal(d.tb, af(expected), actr.Balance())
}

// AssertReceipt checks that a receipt is not nill and has values equal to `expected`.
func (d *StateDriver) AssertReceipt(receipt, expected chain.MessageReceipt) {
	// TODO uncomment when gas values stabalize
	//assert.Equal(t, expected.GasUsed, receipt.GasUsed)
	assert.NotNil(d.tb, receipt)
	assert.Equal(d.tb, expected.ReturnValue, receipt.ReturnValue)
	assert.Equal(d.tb, expected.ExitCode, receipt.ExitCode)
}

// Helpers

func af(v uint64) state.AttoFIL {
	return big.NewInt(0).SetUint64(v)
}

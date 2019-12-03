package suites

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/multsig"
	"github.com/filecoin-project/chain-validation/pkg/state/actors/strgminr"
	"github.com/filecoin-project/chain-validation/pkg/state/address"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
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
func (d *StateDriver) NewAccountActor(balanceAttoFil uint64) address.Address {
	addr, err := d.st.NewAccountAddress()
	require.NoError(d.tb, err)

	_, _, err = d.st.SetActor(addr, actors.AccountActorCodeCid, types.NewInt(balanceAttoFil))
	require.NoError(d.tb, err)
	return addr
}

func (d *StateDriver) NewAccountActorBigBalance(balanceAttoFil types.BigInt) address.Address {
	addr, err := d.st.NewAccountAddress()
	require.NoError(d.tb, err)

	_, _, err = d.st.SetActor(addr, actors.AccountActorCodeCid, balanceAttoFil)
	require.NoError(d.tb, err)
	return addr
}

// TODO add an assertBalancChanged method that can be used until gas is worked out. will need to track previous balances, could be tricky.
// AssertBalance checks an actor has an expected balance.
func (d *StateDriver) AssertBalance(addr address.Address, expected uint64) {
	actr, err := d.st.Actor(addr)
	require.NoError(d.tb, err)
	assert.Equal(d.tb, types.NewInt(expected), actr.Balance(), fmt.Sprintf("expected balance: %v, actual balance: %v", expected, actr.Balance().String()))
}

// AssertReceipt checks that a receipt is not nill and has values equal to `expected`.
func (d *StateDriver) AssertReceipt(receipt, expected chain.MessageReceipt) {
	assert.NotNil(d.tb, receipt)
	// leave gas uncheck for now as it is not speced
	//assert.Equal(d.tb, expected.GasUsed, receipt.GasUsed, fmt.Sprintf("expected gas: %v, actual gas: %v", expected.ExitCode, receipt.GasUsed))
	assert.Equal(d.tb, expected.ReturnValue, receipt.ReturnValue, fmt.Sprintf("expected return value: %v, actual return value: %v", expected.ReturnValue, receipt.ReturnValue))
	assert.Equal(d.tb, expected.ExitCode, receipt.ExitCode, fmt.Sprintf("expected exit code: %v, actual exit code: %v", expected.ExitCode, receipt.ExitCode))
}

func (d *StateDriver) AssertMinerInfo(miner, expected strgminr.MinerInfo) {
	assert.NotNil(d.tb, miner)
	assert.Equal(d.tb, expected.PeerID, miner.PeerID, fmt.Sprintf("expected peerID: %v, actual peerID: %v", expected.PeerID, miner.PeerID))
	assert.Equal(d.tb, expected.Owner, miner.Owner, fmt.Sprintf("expected owner: %v, actual owner: %v", expected.Owner, miner.Owner))
	assert.Equal(d.tb, expected.SectorSize, miner.SectorSize, fmt.Sprintf("expected sector size: %v, actual sector size: %v", expected.SectorSize, miner.SectorSize))
	assert.Equal(d.tb, expected.Worker, miner.Worker, fmt.Sprintf("expected worker: %v, actual worker: %v", expected.Worker, miner.Worker))
}

func (d *StateDriver) AssertMultisigState(multisigAddr address.Address, expected multsig.MultiSigActorState) {
	multisigActor, err := d.State().Actor(multisigAddr)
	require.NoError(d.tb, err)

	multisigStorage, err := d.State().Storage(multisigAddr)
	require.NoError(d.tb, err)

	var multisig multsig.MultiSigActorState
	require.NoError(d.tb, multisigStorage.Get(multisigActor.Head(), &multisig))

	assert.NotNil(d.tb, multisig)
	assert.Equal(d.tb, expected.InitialBalance, multisig.InitialBalance, fmt.Sprintf("expected InitialBalance: %v, actual InitialBalance: %v", expected.InitialBalance, multisig.InitialBalance))
	assert.Equal(d.tb, expected.NextTxID, multisig.NextTxID, fmt.Sprintf("expected NextTxID: %v, actual NextTxID: %v", expected.NextTxID, multisig.NextTxID))
	assert.Equal(d.tb, expected.Required, multisig.Required, fmt.Sprintf("expected Required: %v, actual Required: %v", expected.Required, multisig.Required))
	assert.Equal(d.tb, expected.StartingBlock, multisig.StartingBlock, fmt.Sprintf("expected StartingBlock: %v, actual StartingBlock: %v", expected.StartingBlock, multisig.StartingBlock))
	assert.Equal(d.tb, expected.Transactions, multisig.Transactions, fmt.Sprintf("expected Transactions: %v, actual Transactions: %v", expected.Transactions, multisig.Transactions))
	assert.Equal(d.tb, expected.UnlockDuration, multisig.UnlockDuration, fmt.Sprintf("expected UnlockDuration: %v, actual UnlockDuration: %v", expected.UnlockDuration, multisig.UnlockDuration))

	for _, e := range expected.Signers {
		assert.Contains(d.tb, multisig.Signers, e, fmt.Sprintf("expected Signer: %v, actual Signer: %v", e, multisig.Signers))
	}
}

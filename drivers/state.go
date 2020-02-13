package drivers

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"

	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"

	"github.com/filecoin-project/chain-validation/state"
)

var (
	SECP = address.SECP256K1
	BLS  = address.BLS
)

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

func (d *StateDriver) GetState(c cid.Cid, out cbg.CBORUnmarshaler) {
	strg, err := d.st.Storage()
	require.NoError(d.tb, err)

	err = strg.Get(context.Background(), c, out)
	require.NoError(d.tb, err)
}

func (d *StateDriver) GetActorState(actorAddr address.Address, out cbg.CBORUnmarshaler) {
	actor, err := d.State().Actor(actorAddr)
	require.NoError(d.tb, err)
	require.NotNil(d.tb, actor)

	d.GetState(actor.Head(), out)
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

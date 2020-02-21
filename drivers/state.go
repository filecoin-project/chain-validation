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
	account_spec "github.com/filecoin-project/specs-actors/actors/builtin/account"

	"github.com/filecoin-project/chain-validation/state"
)

var (
	SECP = address.SECP256K1
	BLS  = address.BLS
)

// StateDriver mutates and inspects a state.
type StateDriver struct {
	tb testing.TB
	st state.VMWrapper
	w  state.KeyManager
}

// NewStateDriver creates a new state driver for a state.
func NewStateDriver(tb testing.TB, st state.VMWrapper, w state.KeyManager) *StateDriver {
	return &StateDriver{tb, st, w}
}

// State returns the state.
func (d *StateDriver) State() state.VMWrapper {
	return d.st
}

func (d *StateDriver) Wallet() state.KeyManager {
	return d.w
}

func (d *StateDriver) GetState(c cid.Cid, out cbg.CBORUnmarshaler) {
	err := d.st.Store().Get(context.Background(), c, out)
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
		addr = d.w.NewSECP256k1AccountAddress()
	case address.BLS:
		addr = d.w.NewBLSAccountAddress()
	default:
		require.FailNowf(d.tb, "unsupported address", "protocol for account actor: %v", addrType)
	}

	_, err := d.st.CreateActor(builtin_spec.AccountActorCodeID, addr, balanceAttoFil, &account_spec.State{Address: addr})
	require.NoError(d.tb, err)
	return addr
}

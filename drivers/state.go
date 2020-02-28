package drivers

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	miner_spec "github.com/filecoin-project/specs-actors/actors/builtin/miner"
	power_spec "github.com/filecoin-project/specs-actors/actors/builtin/power"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"

	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	account_spec "github.com/filecoin-project/specs-actors/actors/builtin/account"

	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
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
	rs state.RandomnessSource
}

// NewStateDriver creates a new state driver for a state.
func NewStateDriver(tb testing.TB, st state.VMWrapper, w state.KeyManager, rs state.RandomnessSource) *StateDriver {
	return &StateDriver{tb, st, w, rs}
}

// State returns the state.
func (d *StateDriver) State() state.VMWrapper {
	return d.st
}

func (d *StateDriver) Wallet() state.KeyManager {
	return d.w
}

func (d *StateDriver) Randomness() state.RandomnessSource {
	return d.rs
}

func (d *StateDriver) GetState(c cid.Cid, out cbg.CBORUnmarshaler) {
	err := d.st.Store().Get(context.Background(), c, out)
	require.NoError(d.tb, err)
}

func (d *StateDriver) PutState(in cbg.CBORMarshaler) cid.Cid {
	c, err := d.st.Store().Put(context.Background(), in)
	require.NoError(d.tb, err)
	return c
}

func (d *StateDriver) GetActorState(actorAddr address.Address, out cbg.CBORUnmarshaler) {
	actor, err := d.State().Actor(actorAddr)
	require.NoError(d.tb, err)
	require.NotNil(d.tb, actor)

	d.GetState(actor.Head(), out)
}

// NewAccountActor installs a new account actor, returning the address.
func (d *StateDriver) NewAccountActor(addrType address.Protocol, balanceAttoFil abi_spec.TokenAmount) (pubkey address.Address, id address.Address) {
	var addr address.Address
	switch addrType {
	case address.SECP256K1:
		addr = d.w.NewSECP256k1AccountAddress()
	case address.BLS:
		addr = d.w.NewBLSAccountAddress()
	default:
		require.FailNowf(d.tb, "unsupported address", "protocol for account actor: %v", addrType)
	}

	_, idAddr, err := d.st.CreateActor(builtin_spec.AccountActorCodeID, addr, balanceAttoFil, &account_spec.State{Address: addr})
	require.NoError(d.tb, err)
	return addr, idAddr
}

// create miner without sending a message. modify the init and power actor manually
func (d *StateDriver) newMinerAccountActor() address.Address {

	// creat a miner, owner, and its worker
	_, minerOwnerID := d.NewAccountActor(address.SECP256K1, big_spec.NewInt(1_000_000_000))
	_, minerWorkerID := d.NewAccountActor(address.BLS, big_spec.Zero())
	expectedMinerActorIDAddress := utils.NewIDAddr(d.tb, utils.IdFromAddress(minerWorkerID)+1)
	minerActorAddrs := computeInitActorExecReturn(d.tb, builtin_spec.StoragePowerActorAddr, 0, expectedMinerActorIDAddress)

	// create the miner actor so it exists in the init actors map
	_, minerActorIDAddr, err := d.State().CreateActor(builtin_spec.StorageMinerActorCodeID, minerActorAddrs.RobustAddress, big_spec.Zero(), &miner_spec.State{
		PreCommittedSectors: EmptyMapCid,
		Sectors:             EmptyArrayCid,
		FaultSet:            abi_spec.NewBitField(),
		ProvingSet:          EmptyArrayCid,
		Info: miner_spec.MinerInfo{
			Owner:            minerOwnerID,
			Worker:           minerWorkerID,
			PendingWorkerKey: nil,
			PeerId:           "chain-validation",
			SectorSize:       0,
		},
		PoStState: miner_spec.PoStState{
			ProvingPeriodStart:     -1,
			NumConsecutiveFailures: 0,
		},
	})
	require.NoError(d.tb, err)
	// sanity check above code
	require.Equal(d.tb, expectedMinerActorIDAddress, minerActorIDAddr)
	// great the miner actor has been created, exists in the state tree, and has an entry in the init actor
	// now we need to update the storage power actor such that it is aware of the miner
	// get the spa state
	var spa power_spec.State
	d.GetActorState(builtin_spec.StoragePowerActorAddr, &spa)

	// set the miners balance in the storage power actors state
	table := adt_spec.AsBalanceTable(d.State().Store(), spa.EscrowTable)
	err = table.Set(minerActorIDAddr, big_spec.Zero())
	require.NoError(d.tb, err)
	spa.EscrowTable = table.Root()

	// set the miners claim in the storage power actors state
	hm := adt_spec.AsMap(d.State().Store(), spa.Claims)
	err = hm.Put(adt_spec.AddrKey(minerActorIDAddr), &power_spec.Claim{
		Power:  abi_spec.NewStoragePower(0),
		Pledge: abi_spec.NewTokenAmount(0),
	})
	require.NoError(d.tb, err)
	spa.Claims = hm.Root()

	// now update its state in the tree
	d.PutState(&spa)

	// tada a miner has been created without apply a message
	return minerActorIDAddr
}

package suites

import (
	"math/big"
	"testing"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
	"github.com/stretchr/testify/require"
)

// Driver wraps up all the implementation-specific integration points.
type Driver interface {
	state.Factory
	chain.MessageFactory
	chain.Applier
}

type SugarDrive struct {
	t *testing.T
	Driver
}

func NewSugarDrive(t *testing.T, driver Driver) *SugarDrive {
	return &SugarDrive{
		t:      t,
		Driver: driver,
	}
}

type SugarActor struct {
	actor   state.Actor
	address state.Address
}

func (s *SugarDrive) MakeAccountActor(balance int64) *SugarActor {
	addr, err := s.NewAddress()
	require.NoError(s.t, err)
	actor := s.NewActor(state.AccountActorCodeCid, big.NewInt(balance))
	return &SugarActor{
		actor:   actor,
		address: addr,
	}
}

func (s *SugarDrive) MakeStorageMinerActor(owner state.Address, wokerPk []byte, sectorSize uint64) *SugarActor {
	addr, err := s.NewAddress()
	require.NoError(s.t, err)
	actor := s.NewActor(state.StorageMinerActorCodeCid, big.NewInt(0))
	return &SugarActor{
		actor:   actor,
		address: addr,
	}
}

func (s *SugarDrive) ProducerWithActors(miner *SugarActor, actors ...*SugarActor) *SugarMessageProducer {
	var aa []state.ActorAndAddress
	for _, a := range actors {
		aa = append(aa, state.ActorAndAddress{
			Actor:   a.actor,
			Address: a.address,
		})
	}
	aa = append(aa, state.ActorAndAddress{
		Actor:   miner.actor,
		Address: miner.address,
	})
	tree, storage, err := s.NewState(aa)
	require.NoError(s.t, err)
	require.NotNil(s.t, tree)
	require.NotNil(s.t, storage)

	gasPrice := big.NewInt(1)
	gasLimit := state.GasUnit(1000)
	exeCtx := chain.NewExecutionContext(1, miner.address)
	validator := chain.NewValidator(s.Driver)

	return &SugarMessageProducer{
		t:         s.t,
		producer:  chain.NewMessageProducer(s),
		exeCtx:    exeCtx,
		validator: validator,

		tree:    tree,
		storage: storage,
		miner:   miner,

		gasPrice: gasPrice,
		gasLimit: gasLimit,
	}
}

type SugarMessageProducer struct {
	t         *testing.T
	producer  *chain.MessageProducer
	exeCtx    *chain.ExecutionContext
	validator *chain.Validator

	tree    state.Tree
	storage state.StorageMap
	miner   *SugarActor

	gasPrice state.AttoFIL
	gasLimit state.GasUnit
}

func (s *SugarMessageProducer) Transfer(from, to *SugarActor, amount int64) state.Tree {
	msg, err := s.producer.Transfer(from.address, to.address, big.NewInt(amount), s.gasPrice, s.gasLimit)
	require.NoError(s.t, err)

	endState, msgReceipt, err := s.validator.ApplyMessage(s.exeCtx, s.tree, s.storage, msg)
	require.NoError(s.t, err)
	require.NotNil(s.t, endState)
	require.NotNil(s.t, msgReceipt)

	require.Equal(s.t, uint8(0), msgReceipt.ExitCode)
	require.Equal(s.t, []byte{}, msgReceipt.ReturnValue)
	require.Equal(s.t, state.GasUnit(0), msgReceipt.GasUsed)

	newFromState, err := endState.Actor(from.address)
	require.NoError(s.t, err)
	from.actor = newFromState

	newToState, err := endState.Actor(to.address)
	require.NoError(s.t, err)
	to.actor = newToState

	newMinerState, err := endState.Actor(s.miner.address)
	require.NoError(s.t, err)
	s.miner.actor = newMinerState

	return endState
}

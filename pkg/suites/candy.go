package suites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

// I kinda hate this name
type Candy interface {
	TB() testing.TB
	Driver() *StateDriver
	Producer() *chain.MessageProducer
	Validator() *chain.Validator
	ExeCtx() *chain.ExecutionContext
}

func NewCandy(t testing.TB, factory Factories, singletons map[actors.SingletonActorID]types.BigInt) Candy {
	drv := NewStateDriver(t, factory.NewState())

	// TODO make these function opts
	gasPrice := types.NewInt(1)
	gasLimit := types.NewInt(1000000)

	for sa, balance := range singletons {
		_, _, err := drv.State().SetSingletonActor(sa, balance)
		require.NoError(t, err)
	}

	testMiner := drv.NewAccountActor(0)
	exeCtx := chain.NewExecutionContext(1, testMiner)
	producer := chain.NewMessageProducer(factory.NewMessageFactory(drv.State()), factory.NewActorInfoMapping(), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	return &candy{
		t:         t,
		driver:    drv,
		producer:  producer,
		validator: validator,
		exeCtx:    exeCtx,
	}
}

type candy struct {
	t         testing.TB
	driver    *StateDriver
	producer  *chain.MessageProducer
	validator *chain.Validator
	exeCtx    *chain.ExecutionContext
}

func (c candy) TB() testing.TB {
	return c.t
}

func (c candy) Driver() *StateDriver {
	return c.driver
}

func (c candy) Producer() *chain.MessageProducer {
	return c.producer
}

func (c candy) Validator() *chain.Validator {
	return c.validator
}

func (c candy) ExeCtx() *chain.ExecutionContext {
	return c.exeCtx
}

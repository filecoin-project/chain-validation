package suites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state/actors"
	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

type TestDriver interface {
	TB() testing.TB
	Driver() *StateDriver
	Producer() *chain.MessageProducer
	Validator() *chain.Validator
	ExeCtx() *chain.ExecutionContext
}

func NewTestDriver(t testing.TB, factory Factories, singletons map[actors.SingletonActorID]types.BigInt) TestDriver {
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
	producer := chain.NewMessageProducer(factory.NewMessageFactory(), factory.NewActorInfoMapping(), gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	return &testDriver{
		t:         t,
		driver:    drv,
		producer:  producer,
		validator: validator,
		exeCtx:    exeCtx,
	}
}

type testDriver struct {
	t         testing.TB
	driver    *StateDriver
	producer  *chain.MessageProducer
	validator *chain.Validator
	exeCtx    *chain.ExecutionContext
}

func (c testDriver) TB() testing.TB {
	return c.t
}

func (c testDriver) Driver() *StateDriver {
	return c.driver
}

func (c testDriver) Producer() *chain.MessageProducer {
	return c.producer
}

func (c testDriver) Validator() *chain.Validator {
	return c.validator
}

func (c testDriver) ExeCtx() *chain.ExecutionContext {
	return c.exeCtx
}

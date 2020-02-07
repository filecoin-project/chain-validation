package suites

import (
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/stretchr/testify/require"

	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"

	"github.com/filecoin-project/chain-validation/pkg/chain"
)

type TestDriver interface {
	TB() testing.TB
	Driver() *StateDriver
	Producer() *chain.MessageProducer
	Validator() *chain.Validator
	ExeCtx() *chain.ExecutionContext
}

func NewTestDriver(t testing.TB, factory Factories, singletons map[address.Address]big_spec.Int) TestDriver {
	drv := NewStateDriver(t, factory.NewState())

	// TODO make these function opts
	gasPrice := big_spec.NewInt(1)
	gasLimit := big_spec.NewInt(1000000)

	for sa, balance := range singletons {
		_, _, err := drv.State().SetSingletonActor(sa, balance)
		require.NoError(t, err)
	}

	testMiner := drv.NewAccountActor(0)
	exeCtx := chain.NewExecutionContext(1, testMiner)
	producer := chain.NewMessageProducer(gasLimit, gasPrice)
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

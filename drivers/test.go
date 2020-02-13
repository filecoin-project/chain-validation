package drivers

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/state"
)

var EmptyRetrunValueBytes []byte

func init() {
	buf := new(bytes.Buffer)
	ev := adt_spec.EmptyValue{}
	if err := ev.MarshalCBOR(buf); err != nil {
		panic(err)
	}
	EmptyRetrunValueBytes = buf.Bytes()
}

type mockStore struct {
	ctx context.Context
	cbor.IpldStore
}

func newMockStore() *mockStore {
	bs := blockstore.NewBlockstore(datastore.NewMapDatastore())
	cst := cbor.NewCborStore(bs)
	return &mockStore{
		ctx:       context.Background(),
		IpldStore: cst,
	}
}

func (m mockStore) Context() context.Context {
	return m.ctx
}

type TestDriverBuilder struct {
	ctx     context.Context
	factory state.Factories

	singletons map[address.Address]big_spec.Int
	actors     map[address.Address]big_spec.Int

	defaultMiner    address.Address
	defaultGasPrice big_spec.Int
	defaultGasLimit big_spec.Int
}

func NewBuilder(ctx context.Context, factory state.Factories) *TestDriverBuilder {
	return &TestDriverBuilder{
		factory: factory,
		ctx:     ctx,
	}
}

func (b *TestDriverBuilder) WithSingletonActors(singletons map[address.Address]big_spec.Int) *TestDriverBuilder {
	b.singletons = singletons
	return b
}

func (b *TestDriverBuilder) WithAccountActors(acts map[address.Address]big_spec.Int) *TestDriverBuilder {
	b.actors = acts
	return b
}

func (b *TestDriverBuilder) WithDefaultMiner(miner address.Address) *TestDriverBuilder {
	b.defaultMiner = miner
	return b
}

func (b *TestDriverBuilder) WithDefaultGasLimit(limit big_spec.Int) *TestDriverBuilder {
	b.defaultGasLimit = limit
	return b
}

func (b *TestDriverBuilder) WithDefaultGasPrice(price big_spec.Int) *TestDriverBuilder {
	b.defaultGasPrice = price
	return b
}

func (b *TestDriverBuilder) Build(t testing.TB) *TestDriver {
	sd := NewStateDriver(t, b.factory.NewState())
	for act, bal := range b.singletons {
		// TODO should not ignore the return value here as this should return the ID-address of the miner
		_, _, err := sd.State().SetSingletonActor(act, bal)
		require.NoError(t, err)
	}

	for act, bal := range b.actors {
		// TODO should not ignore the return value here as this should return the ID-address of the miner
		_, _, err := sd.State().SetActor(act, builtin_spec.AccountActorCodeID, bal)
		require.NoError(t, err)
	}

	// TODO should not ignore the return value here as this should return the ID-address of the miner
	_, _, err := sd.st.SetActor(b.defaultMiner, builtin_spec.AccountActorCodeID, big_spec.Zero())
	require.NoError(t, err)

	exeCtx := types.NewExecutionContext(1, b.defaultMiner)
	producer := chain.NewMessageProducer(b.defaultGasLimit, b.defaultGasPrice)
	validator := chain.NewValidator(b.factory)
	return &TestDriver{
		T:           t,
		StateDriver: sd,
		Producer:    producer,
		Validator:   validator,
		ExeCtx:      exeCtx,
	}
}

type TestDriver struct {
	*StateDriver

	T         testing.TB
	Producer  *chain.MessageProducer
	Validator *chain.Validator
	ExeCtx    *types.ExecutionContext
}

// TODO for failure cases we should consider catching panics here else they appear in the test output and obfuscate successful tests.
func (td *TestDriver) ApplyMessageExpectReceipt(msg *types.Message, receipt types.MessageReceipt) {
	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State(), msg)
	require.NoError(td.T, err)

	require.Equal(td.T, receipt.GasUsed, msgReceipt.GasUsed)
	require.Equal(td.T, receipt.ExitCode, msgReceipt.ExitCode)
	require.Equal(td.T, receipt.ReturnValue, msgReceipt.ReturnValue)
}

// AssertBalance checks an actor has an expected balance.
func (td *TestDriver) AssertBalance(addr address.Address, expected big_spec.Int) {
	actr, err := td.State().Actor(addr)
	require.NoError(td.T, err)
	assert.Equal(td.T, expected, actr.Balance(), fmt.Sprintf("expected balance: %v, actual balance: %v", expected, actr.Balance().String()))
}

func (td *TestDriver) AssertMultisigTransaction(multisigAddr address.Address, txnID multisig_spec.TxnID, txn multisig_spec.Transaction) {
	var msState multisig_spec.State
	td.GetActorState(multisigAddr, &msState)

	strg, err := td.State().Storage()
	require.NoError(td.T, err)

	txnMap := adt_spec.AsMap(strg, msState.PendingTxns)
	var actualTxn multisig_spec.Transaction
	found, err := txnMap.Get(txnID, &actualTxn)
	assert.NoError(td.T, err)
	assert.True(td.T, found)

	assert.Equal(td.T, txn, actualTxn)
}

func (td *TestDriver) AssertMultisigContainsTransaction(multisigAddr address.Address, txnID multisig_spec.TxnID, contains bool) {
	var msState multisig_spec.State
	td.GetActorState(multisigAddr, &msState)

	strg, err := td.State().Storage()
	require.NoError(td.T, err)

	txnMap := adt_spec.AsMap(strg, msState.PendingTxns)
	var actualTxn multisig_spec.Transaction
	found, err := txnMap.Get(txnID, &actualTxn)
	require.NoError(td.T, err)
	assert.Equal(td.T, contains, found)

}

func (td *TestDriver) AssertMultisigState(multisigAddr address.Address, expected multisig_spec.State) {
	var msState multisig_spec.State
	td.GetActorState(multisigAddr, &msState)
	assert.NotNil(td.T, msState)
	assert.Equal(td.T, expected.InitialBalance, msState.InitialBalance, fmt.Sprintf("expected InitialBalance: %v, actual InitialBalance: %v", expected.InitialBalance, msState.InitialBalance))
	assert.Equal(td.T, expected.NextTxnID, msState.NextTxnID, fmt.Sprintf("expected NextTxnID: %v, actual NextTxnID: %v", expected.NextTxnID, msState.NextTxnID))
	assert.Equal(td.T, expected.NumApprovalsThreshold, msState.NumApprovalsThreshold, fmt.Sprintf("expected NumApprovalsThreshold: %v, actual NumApprovalsThreshold: %v", expected.NumApprovalsThreshold, msState.NumApprovalsThreshold))
	assert.Equal(td.T, expected.StartEpoch, msState.StartEpoch, fmt.Sprintf("expected StartEpoch: %v, actual StartEpoch: %v", expected.StartEpoch, msState.StartEpoch))
	assert.Equal(td.T, expected.UnlockDuration, msState.UnlockDuration, fmt.Sprintf("expected UnlockDuration: %v, actual UnlockDuration: %v", expected.UnlockDuration, msState.UnlockDuration))

	for _, e := range expected.Signers {
		assert.Contains(td.T, msState.Signers, e, fmt.Sprintf("expected Signer: %v, actual Signer: %v", e, msState.Signers))
	}
}

func (td *TestDriver) MustCreateAndVerifyMultisigActor(nonce int64, value abi_spec.TokenAmount, multisigAddr address.Address, from address.Address, params *multisig_spec.ConstructorParams, receipt types.MessageReceipt) {
	/* Create the Multisig actor*/
	td.ApplyMessageExpectReceipt(
		td.Producer.CreateMultisigActor(from, params.Signers, params.UnlockDuration, params.NumApprovalsThreshold, chain.Nonce(nonce), chain.Value(value)),
		receipt,
	)
	/* Assert the actor state was setup as expected */
	pendingTxMap, err := adt_spec.MakeEmptyMap(newMockStore())
	require.NoError(td.T, err)
	td.AssertMultisigState(multisigAddr, multisig_spec.State{
		NextTxnID:      0,
		InitialBalance: value,
		StartEpoch:     td.ExeCtx.Epoch,

		Signers:               params.Signers,
		UnlockDuration:        params.UnlockDuration,
		NumApprovalsThreshold: params.NumApprovalsThreshold,

		PendingTxns: pendingTxMap.Root(),
	})
	td.AssertBalance(multisigAddr, value)
}

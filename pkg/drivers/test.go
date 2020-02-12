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
	init_spec "github.com/filecoin-project/specs-actors/actors/builtin/init"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/chain/types"
	"github.com/filecoin-project/chain-validation/pkg/state"
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
		T:         t,
		State:     sd,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}
}

type TestDriver struct {
	T         testing.TB
	State     *StateDriver
	Producer  *chain.MessageProducer
	Validator *chain.Validator
	ExeCtx    *types.ExecutionContext
}

// TODO for failure cases we should consider catching panics here else they appear in the test output and obfuscate successful tests.
func (td *TestDriver) ApplyMessageExpectReceipt(msg *types.Message, receipt types.MessageReceipt) {
	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State.st, msg)
	require.NoError(td.T, err)

	require.Equal(td.T, receipt.GasUsed, msgReceipt.GasUsed)
	require.Equal(td.T, receipt.ExitCode, msgReceipt.ExitCode)
	require.Equal(td.T, receipt.ReturnValue, msgReceipt.ReturnValue)
}

// AssertBalance checks an actor has an expected balance.
func (td *TestDriver) AssertBalance(addr address.Address, expected big_spec.Int) {
	actr, err := td.State.State().Actor(addr)
	require.NoError(td.T, err)
	assert.Equal(td.T, expected, actr.Balance(), fmt.Sprintf("expected balance: %v, actual balance: %v", expected, actr.Balance().String()))
}

// AssertReceipt checks that a receipt is not nill and has values equal to `expected`.
func (td *TestDriver) AssertReceipt(receipt, expected types.MessageReceipt) {
	assert.NotNil(td.T, receipt)
	assert.Equal(td.T, expected.GasUsed, receipt.GasUsed, fmt.Sprintf("expected gas: %v, actual gas: %v", expected.GasUsed, receipt.GasUsed))
	assert.Equal(td.T, expected.ReturnValue, receipt.ReturnValue, fmt.Sprintf("expected return value: %v, actual return value: %v", expected.ReturnValue, receipt.ReturnValue))
	assert.Equal(td.T, expected.ExitCode, receipt.ExitCode, fmt.Sprintf("expected exit code: %v, actual exit code: %v", expected.ExitCode, receipt.ExitCode))
}

func (td *TestDriver) AssertMultisigTransaction(multisigAddr address.Address, txnID multisig_spec.TxnID, txn multisig_spec.Transaction) {
	multisigActor, err := td.State.State().Actor(multisigAddr)
	require.NoError(td.T, err)

	strg, err := td.State.State().Storage()
	require.NoError(td.T, err)

	var multisig multisig_spec.State
	require.NoError(td.T, strg.Get(context.Background(), multisigActor.Head(), &multisig))

	txnMap := adt_spec.AsMap(strg, multisig.PendingTxns)
	var actualTxn multisig_spec.Transaction
	found, err := txnMap.Get(txnID, &actualTxn)
	assert.NoError(td.T, err)
	assert.True(td.T, found)

	assert.Equal(td.T, txn, actualTxn)
}

func (td *TestDriver) AssertMultisigContainsTransaction(multisigAddr address.Address, txnID multisig_spec.TxnID, contains bool) {
	multisigActor, err := td.State.State().Actor(multisigAddr)
	require.NoError(td.T, err)

	strg, err := td.State.State().Storage()
	require.NoError(td.T, err)

	var multisig multisig_spec.State
	require.NoError(td.T, strg.Get(context.Background(), multisigActor.Head(), &multisig))

	txnMap := adt_spec.AsMap(strg, multisig.PendingTxns)
	var actualTxn multisig_spec.Transaction
	found, err := txnMap.Get(txnID, &actualTxn)
	require.NoError(td.T, err)
	assert.Equal(td.T, contains, found)

}

func (td *TestDriver) AssertMultisigState(multisigAddr address.Address, expected multisig_spec.State) {
	multisigActor, err := td.State.State().Actor(multisigAddr)
	require.NoError(td.T, err)

	strg, err := td.State.State().Storage()
	require.NoError(td.T, err)

	var multisig multisig_spec.State
	require.NoError(td.T, strg.Get(context.Background(), multisigActor.Head(), &multisig))

	assert.NotNil(td.T, multisig)
	assert.Equal(td.T, expected.InitialBalance, multisig.InitialBalance, fmt.Sprintf("expected InitialBalance: %v, actual InitialBalance: %v", expected.InitialBalance, multisig.InitialBalance))
	assert.Equal(td.T, expected.NextTxnID, multisig.NextTxnID, fmt.Sprintf("expected NextTxnID: %v, actual NextTxnID: %v", expected.NextTxnID, multisig.NextTxnID))
	assert.Equal(td.T, expected.NumApprovalsThreshold, multisig.NumApprovalsThreshold, fmt.Sprintf("expected NumApprovalsThreshold: %v, actual NumApprovalsThreshold: %v", expected.NumApprovalsThreshold, multisig.NumApprovalsThreshold))
	assert.Equal(td.T, expected.StartEpoch, multisig.StartEpoch, fmt.Sprintf("expected StartEpoch: %v, actual StartEpoch: %v", expected.StartEpoch, multisig.StartEpoch))
	assert.Equal(td.T, expected.UnlockDuration, multisig.UnlockDuration, fmt.Sprintf("expected UnlockDuration: %v, actual UnlockDuration: %v", expected.UnlockDuration, multisig.UnlockDuration))

	for _, e := range expected.Signers {
		assert.Contains(td.T, multisig.Signers, e, fmt.Sprintf("expected Signer: %v, actual Signer: %v", e, multisig.Signers))
	}
}

func (td *TestDriver) MustCreateAndVerifyMultisigActor(nonce int64, value abi_spec.TokenAmount, multisigAddr address.Address, from address.Address, params *multisig_spec.ConstructorParams, receipt types.MessageReceipt) {
	/* Create the Multisig actor*/
	multiSigConstuctParams, err := chain.Serialize(params)
	require.NoError(td.T, err)

	msg := td.Producer.InitExec(builtin_spec.InitActorAddr, from, init_spec.ExecParams{
		CodeCID:           builtin_spec.MultisigActorCodeID,
		ConstructorParams: multiSigConstuctParams,
	}, chain.Value(value), chain.Nonce(nonce))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State.State(), msg)
	require.NoError(td.T, err)

	/* Assert the message was applied successfully  */
	td.AssertReceipt(msgReceipt, receipt)

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

func (td *TestDriver) MustProposeMultisigTransfer(nonce int64, value abi_spec.TokenAmount, txID multisig_spec.TxnID, multisigAddr, from address.Address, params multisig_spec.ProposeParams, receipt types.MessageReceipt) {
	/* Propose the transactions */
	msg := td.Producer.MultisigPropose(multisigAddr, from, params, chain.Value(value), chain.Nonce(nonce))
	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State.State(), msg)
	require.NoError(td.T, err)

	/* Assert it was applied successfully  */
	td.AssertReceipt(msgReceipt, receipt)
}

func (td *TestDriver) MustApproveMultisigActor(nonce int64, value abi_spec.TokenAmount, ms, from address.Address, txID multisig_spec.TxnID, receipt types.MessageReceipt) {
	msg := td.Producer.MultisigApprove(ms, from, multisig_spec.TxnIDParams{ID: txID}, chain.Value(value), chain.Nonce(nonce))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State.State(), msg)
	require.NoError(td.T, err)

	td.AssertReceipt(msgReceipt, receipt)
}

func (td *TestDriver) MustCancelMultisigActor(nonce int64, value abi_spec.TokenAmount, ms, from address.Address, txID multisig_spec.TxnID, receipt types.MessageReceipt) {
	msg := td.Producer.MultisigCancel(ms, from, multisig_spec.TxnIDParams{ID: txID}, chain.Value(value), chain.Nonce(nonce))

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State.State(), msg)
	require.NoError(td.T, err)

	td.AssertReceipt(msgReceipt, receipt)
}

package suites

import (
	"bytes"
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	init_spec "github.com/filecoin-project/specs-actors/actors/builtin/init"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/chain"
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
	factory Factories

	singletons map[address.Address]big_spec.Int
	actors     map[address.Address]big_spec.Int

	defaultMiner    address.Address
	defaultGasPrice big_spec.Int
	defaultGasLimit big_spec.Int
}

func NewBuilder(ctx context.Context, factory Factories) *TestDriverBuilder {
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

	exeCtx := chain.NewExecutionContext(1, b.defaultMiner)
	producer := chain.NewMessageProducer(b.defaultGasLimit, b.defaultGasPrice)
	validator := chain.NewValidator(b.factory)
	return &TestDriver{
		T:         t,
		Driver:    sd,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}
}

type TestDriver struct {
	T         testing.TB
	Driver    *StateDriver
	Producer  *chain.MessageProducer
	Validator *chain.Validator
	ExeCtx    *chain.ExecutionContext
}

// TODO for failure cases we should consider catching panics here else they appear in the test output and obfuscate successful tests.
func (td *TestDriver) ApplyMessageExpectReceipt(msgF func() (*chain.Message, error), receipt chain.MessageReceipt) {
	msg, err := msgF()
	require.NoError(td.T, err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.st, msg)
	require.NoError(td.T, err)

	require.Equal(td.T, receipt.GasUsed, msgReceipt.GasUsed)
	require.Equal(td.T, receipt.ExitCode, msgReceipt.ExitCode)
	require.Equal(td.T, receipt.ReturnValue, msgReceipt.ReturnValue)
}

func (td *TestDriver) MustCreateAndVerifyMultisigActor(nonce int64, value abi_spec.TokenAmount, multisigAddr address.Address, from address.Address, params *multisig_spec.ConstructorParams, receipt chain.MessageReceipt) {
	/* Create the Multisig actor*/
	multiSigConstuctParams, err := state.Serialize(params)
	require.NoError(td.T, err)

	msg, err := td.Producer.InitExec(builtin_spec.InitActorAddr, from, init_spec.ExecParams{
		CodeCID:           builtin_spec.MultisigActorCodeID,
		ConstructorParams: multiSigConstuctParams,
	}, chain.Value(value), chain.Nonce(nonce))
	require.NoError(td.T, err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(td.T, err)

	/* Assert the message was applied successfully  */
	td.Driver.AssertReceipt(msgReceipt, receipt)

	/* Assert the actor state was setup as expected */
	pendingTxMap, err := adt_spec.MakeEmptyMap(newMockStore())
	require.NoError(td.T, err)
	td.Driver.AssertMultisigState(multisigAddr, multisig_spec.State{
		NextTxnID:      0,
		InitialBalance: value,
		StartEpoch:     td.ExeCtx.Epoch,

		Signers:               params.Signers,
		UnlockDuration:        params.UnlockDuration,
		NumApprovalsThreshold: params.NumApprovalsThreshold,

		PendingTxns: pendingTxMap.Root(),
	})
	td.Driver.AssertBalance(multisigAddr, value)
}

func (td *TestDriver) MustProposeMultisigTransfer(nonce int64, value abi_spec.TokenAmount, txID multisig_spec.TxnID, multisigAddr, from address.Address, params multisig_spec.ProposeParams, receipt chain.MessageReceipt) {
	/* Propose the transactions */
	msg, err := td.Producer.MultisigPropose(multisigAddr, from, params, chain.Value(value), chain.Nonce(nonce))
	require.NoError(td.T, err)
	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(td.T, err)

	/* Assert it was applied successfully  */
	td.Driver.AssertReceipt(msgReceipt, receipt)
}

func (td *TestDriver) MustApproveMultisigActor(nonce int64, value abi_spec.TokenAmount, ms, from address.Address, txID multisig_spec.TxnID, receipt chain.MessageReceipt) {
	msg, err := td.Producer.MultisigApprove(ms, from, multisig_spec.TxnIDParams{ID: txID}, chain.Value(value), chain.Nonce(nonce))
	require.NoError(td.T, err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(td.T, err)

	td.Driver.AssertReceipt(msgReceipt, receipt)
}

func (td *TestDriver) MustCancelMultisigActor(nonce int64, value abi_spec.TokenAmount, ms, from address.Address, txID multisig_spec.TxnID, receipt chain.MessageReceipt) {
	msg, err := td.Producer.MultisigCancel(ms, from, multisig_spec.TxnIDParams{ID: txID}, chain.Value(value), chain.Nonce(nonce))
	require.NoError(td.T, err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(td.T, err)

	td.Driver.AssertReceipt(msgReceipt, receipt)
}

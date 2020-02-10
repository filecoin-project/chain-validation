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

type TestDriver struct {
	T         testing.TB
	Driver    *StateDriver
	Producer  *chain.MessageProducer
	Validator *chain.Validator
	ExeCtx    *chain.ExecutionContext
}

func NewTestDriver(t testing.TB, factory Factories, singletons map[address.Address]big_spec.Int) *TestDriver {
	drv := NewStateDriver(t, factory.NewState())

	// TODO make these function opts
	gasPrice := big_spec.NewInt(1)
	gasLimit := big_spec.NewInt(1000000)

	for sa, balance := range singletons {
		_, _, err := drv.State().SetSingletonActor(sa, balance)
		require.NoError(t, err)
	}

	testMiner := drv.NewAccountActor(BLS, big_spec.Zero())
	exeCtx := chain.NewExecutionContext(1, testMiner)
	producer := chain.NewMessageProducer(gasLimit, gasPrice)
	validator := chain.NewValidator(factory)

	return &TestDriver{
		T:         t,
		Driver:    drv,
		Producer:  producer,
		Validator: validator,
		ExeCtx:    exeCtx,
	}
}

func (td *TestDriver) ApplyMessageExpectReceipt(msgF func() (*chain.Message, error), receipt chain.MessageReceipt) {
	msg, err := msgF()
	require.NoError(td.T, err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.st, msg)
	require.NoError(td.T, err)

	require.Equal(td.T, receipt.GasUsed, msgReceipt.GasUsed)
	require.Equal(td.T, receipt.ExitCode, msgReceipt.ExitCode)
	require.Equal(td.T, receipt.ReturnValue, msgReceipt.ReturnValue)
}

// TODO all Must* methods need to be adapted to assert gas values correctly.

func (td *TestDriver) MustCreateAndVerifyMultisigActor(nonce int64, value abi_spec.TokenAmount, multisigAddr address.Address, from address.Address, params *multisig_spec.ConstructorParams) {
	/* Create the Multisig actor*/
	multiSigConstuctParams, err := state.Serialize(params)
	require.NoError(td.T, err)

	msg, err := td.Producer.InitExec(from, builtin_spec.MultisigActorCodeID, multiSigConstuctParams, chain.Value(value), chain.Nonce(nonce))
	require.NoError(td.T, err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(td.T, err)

	/* Assert the message was applied successfully  */
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: multisigAddr.Bytes(),
		GasUsed:     big_spec.Zero(),
	})

	/* Assert the actor state was setup as expected */
	pendingTxMap, err := adt_spec.MakeEmptyMap(newMockStore())
	require.NoError(td.T, err)
	td.Driver.AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
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

func (td *TestDriver) MustProposeMultisigTransfer(nonce int64, value abi_spec.TokenAmount, txID multisig_spec.TxnID, multisigAddr, from address.Address, params multisig_spec.ProposeParams) {
	/* Propose the transactions */
	msg, err := td.Producer.MultiSigPropose(multisigAddr, from, params, chain.Value(value), chain.Nonce(nonce))
	require.NoError(td.T, err)
	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(td.T, err)

	/* Assert it was applied successfully  */
	btxid, err := state.Serialize(&multisig_spec.TxnIDParams{ID: txID})
	require.NoError(td.T, err)
	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode: 0,
		// since the first byte is the cbor type indicator.
		ReturnValue: btxid[1:],
		GasUsed:     big_spec.NewInt(0),
	})
}

func (td *TestDriver) MustApproveMultisigActor(nonce int64, value abi_spec.TokenAmount, ms, from address.Address, txID multisig_spec.TxnID) {
	msg, err := td.Producer.MultiSigApprove(ms, from, txID, chain.Value(value), chain.Nonce(nonce))
	require.NoError(td.T, err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(td.T, err)

	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: EmptyRetrunValueBytes,
		GasUsed:     big_spec.NewInt(0),
	})
}

func (td *TestDriver) MustCancelMultisigActor(nonce int64, value abi_spec.TokenAmount, ms, from address.Address, txID multisig_spec.TxnID) {
	msg, err := td.Producer.MultiSigCancel(ms, from, txID, chain.Value(value), chain.Nonce(nonce))
	require.NoError(td.T, err)

	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.Driver.State(), msg)
	require.NoError(td.T, err)

	td.Driver.AssertReceipt(msgReceipt, chain.MessageReceipt{
		ExitCode:    0,
		ReturnValue: EmptyRetrunValueBytes,
		GasUsed:     big_spec.NewInt(0),
	})
}

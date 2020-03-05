package drivers

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	account_spec "github.com/filecoin-project/specs-actors/actors/builtin/account"
	cron_spec "github.com/filecoin-project/specs-actors/actors/builtin/cron"
	init_spec "github.com/filecoin-project/specs-actors/actors/builtin/init"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	power_spec "github.com/filecoin-project/specs-actors/actors/builtin/power"
	reward_spec "github.com/filecoin-project/specs-actors/actors/builtin/reward"
	runtime_spec "github.com/filecoin-project/specs-actors/actors/runtime"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/state"
)

var (

	// initialized by calling initializeStoreWithAdtRoots
	EmptyArrayCid    cid.Cid
	EmptyMapCid      cid.Cid
	EmptyMultiMapCid cid.Cid
	EmptySetCid      cid.Cid
)

var (
	DefaultInitActorState         ActorState
	DefaultRewardActorState       ActorState
	DefaultBurntFundsActorState   ActorState
	DefaultStoragePowerActorState ActorState
	DefaultSystemActorState       ActorState
	DefaultCronActorState         ActorState
	DefaultBuiltinActorsState     []ActorState
)

func init() {
	ms := newMockStore()
	if err := initializeStoreWithAdtRoots(ms); err != nil {
		panic(err)
	}

	DefaultInitActorState = ActorState{
		Addr:    builtin_spec.InitActorAddr,
		Balance: big_spec.Zero(),
		Code:    builtin_spec.InitActorCodeID,
		State:   init_spec.ConstructState(EmptyMapCid, "chain-validation"),
	}

	DefaultRewardActorState = ActorState{
		Addr:    builtin_spec.RewardActorAddr,
		Balance: TotalNetworkBalance,
		Code:    builtin_spec.RewardActorCodeID,
		State:   reward_spec.ConstructState(EmptyMultiMapCid),
	}

	DefaultBurntFundsActorState = ActorState{
		Addr:    builtin_spec.BurntFundsActorAddr,
		Balance: big_spec.Zero(),
		Code:    builtin_spec.AccountActorCodeID,
		State:   &account_spec.State{Address: builtin_spec.BurntFundsActorAddr},
	}

	DefaultStoragePowerActorState = ActorState{
		Addr:    builtin_spec.StoragePowerActorAddr,
		Balance: big_spec.Zero(),
		Code:    builtin_spec.StoragePowerActorCodeID,
		State: &power_spec.State{
			TotalNetworkPower:        abi_spec.NewStoragePower(0),
			EscrowTable:              EmptyMapCid,
			CronEventQueue:           EmptyMapCid,
			PoStDetectedFaultMiners:  EmptyMapCid,
			Claims:                   EmptyMapCid,
			NumMinersMeetingMinPower: 0,
		},
	}

	DefaultSystemActorState = ActorState{
		Addr:    builtin_spec.SystemActorAddr,
		Balance: big_spec.Zero(),
		Code:    builtin_spec.SystemActorCodeID,
		State:   &adt_spec.EmptyValue{},
	}

	DefaultCronActorState = ActorState{
		Addr:    builtin_spec.CronActorAddr,
		Balance: big_spec.Zero(),
		Code:    builtin_spec.CronActorCodeID,
		State:   &cron_spec.State{Entries: []cron_spec.Entry(nil)},
	}

	DefaultBuiltinActorsState = []ActorState{
		DefaultInitActorState,
		DefaultRewardActorState,
		DefaultBurntFundsActorState,
		DefaultStoragePowerActorState,
		DefaultSystemActorState,
		DefaultCronActorState,
	}
}

func initializeStoreWithAdtRoots(store adt_spec.Store) error {
	emptyArray, err := adt_spec.MakeEmptyArray(store)
	if err != nil {
		return err
	}
	EmptyArrayCid = emptyArray.Root()

	emptyMap, err := adt_spec.MakeEmptyMap(store)
	if err != nil {
		return err
	}
	EmptyMapCid = emptyMap.Root()

	emptyMultiMap, err := adt_spec.MakeEmptyMultimap(store)
	if err != nil {
		return err
	}
	EmptyMultiMapCid = emptyMultiMap.Root()

	emptySet, err := adt_spec.MakeEmptySet(store)
	if err != nil {
		return err
	}
	EmptySetCid = emptySet.Root()
	return nil
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

	actorStates []ActorState

	defaultGasPrice big_spec.Int
	defaultGasLimit int64
}

func NewBuilder(ctx context.Context, factory state.Factories) *TestDriverBuilder {
	return &TestDriverBuilder{
		factory: factory,
		ctx:     ctx,
	}
}

type ActorState struct {
	Addr    address.Address
	Balance big_spec.Int
	Code    cid.Cid
	State   runtime_spec.CBORMarshaler
}

func (b *TestDriverBuilder) WithActorState(acts []ActorState) *TestDriverBuilder {
	b.actorStates = acts
	return b
}

func (b *TestDriverBuilder) WithDefaultGasLimit(limit int64) *TestDriverBuilder {
	b.defaultGasLimit = limit
	return b
}

func (b *TestDriverBuilder) WithDefaultGasPrice(price big_spec.Int) *TestDriverBuilder {
	b.defaultGasPrice = price
	return b
}

func (b *TestDriverBuilder) Build(t testing.TB) *TestDriver {
	sd := NewStateDriver(t, b.factory.NewState(), b.factory.NewKeyManager(), b.factory.NewRandomnessSource())

	err := initializeStoreWithAdtRoots(sd.st.Store())
	require.NoError(t, err)

	for _, acts := range b.actorStates {
		_, _, err := sd.State().CreateActor(acts.Code, acts.Addr, acts.Balance, acts.State)
		require.NoError(t, err)
	}

	minerActorIDAddr := sd.newMinerAccountActor()

	exeCtx := types.NewExecutionContext(1, minerActorIDAddr)
	producer := chain.NewMessageProducer(b.defaultGasLimit, b.defaultGasPrice)
	validator := chain.NewValidator(b.factory)

	return &TestDriver{
		T:               t,
		StateDriver:     sd,
		MessageProducer: producer,
		Validator:       validator,
		ExeCtx:          exeCtx,

		Config: b.factory.NewValidationConfig(),
	}

}

type TestDriver struct {
	*StateDriver

	T                    testing.TB
	MessageProducer      *chain.MessageProducer
	TipSetMessageBuilder *chain.TipSetMessageBuilder
	Validator            *chain.Validator
	ExeCtx               *types.ExecutionContext

	Config state.ValidationConfig
}

// TODO for failure cases we should consider catching panics here else they appear in the test output and obfuscate successful tests.
func (td *TestDriver) ApplyMessageExpectReceipt(msg *types.Message, receipt types.MessageReceipt) {
	msgReceipt, err := td.Validator.ApplyMessage(td.ExeCtx, td.State(), msg)
	require.NoError(td.T, err)
	td.AssertReceipt(msgReceipt, receipt)

}

// AssertBalance checks an actor has an expected balance.
func (td *TestDriver) AssertBalance(addr address.Address, expected big_spec.Int) {
	actr, err := td.State().Actor(addr)
	require.NoError(td.T, err)
	assert.Equal(td.T, expected, actr.Balance(), fmt.Sprintf("expected balance: %v, actual balance: %v", expected, actr.Balance().String()))
}

func (td *TestDriver) AssertBalanceWithGas(addr address.Address, expected, gasUsed big_spec.Int) {
	if td.Config.ValidateGas() {
		actr, err := td.State().Actor(addr)
		require.NoError(td.T, err)
		assert.Equal(td.T, expected, big_spec.Sub(actr.Balance(), gasUsed), fmt.Sprintf("expected balance: %v, actual balance: %v", expected, actr.Balance().String()))
	}
}

func (td *TestDriver) AssertReceipt(actual, expected types.MessageReceipt) {
	if td.Config.ValidateGas() {
		assert.Equal(td.T, expected.GasUsed, actual.GasUsed, "Expected GasUsed: %s Actual GasUsed: %s", expected.GasUsed.String(), actual.GasUsed.String())
	}
	if td.Config.ValidateExitCode() {
		assert.Equal(td.T, expected.ExitCode, actual.ExitCode, "Expected ExitCode: %s Actual ExitCode: %s", expected.ExitCode.Error(), actual.ExitCode.Error())
	}
	if td.Config.ValidateReturnValue() {
		assert.Equal(td.T, expected.ReturnValue, actual.ReturnValue, "Expected ReturnValue: %v Actual ReturnValue: %v", expected.ReturnValue, actual.ReturnValue)
	}
}

func (td *TestDriver) AssertMultisigTransaction(multisigAddr address.Address, txnID multisig_spec.TxnID, txn multisig_spec.Transaction) {
	var msState multisig_spec.State
	td.GetActorState(multisigAddr, &msState)

	txnMap := adt_spec.AsMap(td.State().Store(), msState.PendingTxns)
	var actualTxn multisig_spec.Transaction
	found, err := txnMap.Get(txnID, &actualTxn)
	assert.NoError(td.T, err)
	assert.True(td.T, found)

	assert.Equal(td.T, txn, actualTxn)
}

func (td *TestDriver) AssertMultisigContainsTransaction(multisigAddr address.Address, txnID multisig_spec.TxnID, contains bool) {
	var msState multisig_spec.State
	td.GetActorState(multisigAddr, &msState)

	txnMap := adt_spec.AsMap(td.State().Store(), msState.PendingTxns)
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

func (td *TestDriver) ComputeInitActorExecReturn(from address.Address, callSeq uint64, expectedNewAddr address.Address) init_spec.ExecReturn {
	return computeInitActorExecReturn(td.T, from, callSeq, expectedNewAddr)
}

func computeInitActorExecReturn(t testing.TB, from address.Address, callSeq uint64, expectedNewAddr address.Address) init_spec.ExecReturn {
	buf := new(bytes.Buffer)

	require.NoError(t, from.MarshalCBOR(buf))
	require.NoError(t, binary.Write(buf, binary.BigEndian, callSeq))
	require.NoError(t, binary.Write(buf, binary.BigEndian, uint64(0))) // TODO this is wrong, but lotus does it this way

	out, err := address.NewActorAddress(buf.Bytes())
	require.NoError(t, err)

	return init_spec.ExecReturn{
		IDAddress:     expectedNewAddr,
		RobustAddress: out,
	}

}

func (td *TestDriver) MustCreateAndVerifyMultisigActor(nonce int64, value abi_spec.TokenAmount, multisigAddr address.Address, from address.Address, params *multisig_spec.ConstructorParams, receipt types.MessageReceipt) {
	/* Create the Multisig actor*/
	td.ApplyMessageExpectReceipt(
		td.MessageProducer.CreateMultisigActor(from, params.Signers, params.UnlockDuration, params.NumApprovalsThreshold, chain.Nonce(nonce), chain.Value(value)),
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

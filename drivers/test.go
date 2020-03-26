package drivers

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	account_spec "github.com/filecoin-project/specs-actors/actors/builtin/account"
	cron_spec "github.com/filecoin-project/specs-actors/actors/builtin/cron"
	init_spec "github.com/filecoin-project/specs-actors/actors/builtin/init"
	market_spec "github.com/filecoin-project/specs-actors/actors/builtin/market"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	power_spec "github.com/filecoin-project/specs-actors/actors/builtin/power"
	reward_spec "github.com/filecoin-project/specs-actors/actors/builtin/reward"
	"github.com/filecoin-project/specs-actors/actors/builtin/system"
	runtime_spec "github.com/filecoin-project/specs-actors/actors/runtime"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	adt_spec "github.com/filecoin-project/specs-actors/actors/util/adt"
	"github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/gasmeter"
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
	DefaultInitActorState          ActorState
	DefaultRewardActorState        ActorState
	DefaultBurntFundsActorState    ActorState
	DefaultStoragePowerActorState  ActorState
	DefaultStorageMarketActorState ActorState
	DefaultSystemActorState        ActorState
	DefaultCronActorState          ActorState
	DefaultBuiltinActorsState      []ActorState
)

// if this number is 0 we get a specs-actors panic since it divides by 0
const InitialTotalNetworkPower = 1

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
			TotalNetworkPower:        abi_spec.NewStoragePower(InitialTotalNetworkPower),
			EscrowTable:              EmptyMapCid,
			CronEventQueue:           EmptyMapCid,
			PoStDetectedFaultMiners:  EmptyMapCid,
			Claims:                   EmptyMapCid,
			NumMinersMeetingMinPower: 0,
		},
	}

	DefaultStorageMarketActorState = ActorState{
		Addr:    builtin_spec.StorageMarketActorAddr,
		Balance: big_spec.Zero(),
		Code:    builtin_spec.StorageMarketActorCodeID,
		State: &market_spec.State{
			Proposals:      EmptyArrayCid,
			States:         EmptyArrayCid,
			EscrowTable:    EmptyMapCid,
			LockedTable:    EmptyMapCid,
			NextID:         abi_spec.DealID(0),
			DealIDsByParty: EmptyMultiMapCid,
		},
	}

	DefaultSystemActorState = ActorState{
		Addr:    builtin_spec.SystemActorAddr,
		Balance: big_spec.Zero(),
		Code:    builtin_spec.SystemActorCodeID,
		State:   &system.State{},
	}

	DefaultCronActorState = ActorState{
		Addr:    builtin_spec.CronActorAddr,
		Balance: big_spec.Zero(),
		Code:    builtin_spec.CronActorCodeID,
		State: &cron_spec.State{Entries: []cron_spec.Entry{
			{
				Receiver:  builtin_spec.StoragePowerActorAddr,
				MethodNum: builtin_spec.MethodsPower.OnEpochTickEnd,
			},
		}},
	}

	DefaultBuiltinActorsState = []ActorState{
		DefaultInitActorState,
		DefaultRewardActorState,
		DefaultBurntFundsActorState,
		DefaultStoragePowerActorState,
		DefaultStorageMarketActorState,
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

	defaultGasPrice abi_spec.TokenAmount
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
	Balance abi_spec.TokenAmount
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

func (b *TestDriverBuilder) WithDefaultGasPrice(price abi_spec.TokenAmount) *TestDriverBuilder {
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
		validator:       validator,
		ExeCtx:          exeCtx,

		Config: b.factory.NewValidationConfig(),

		GasMeter: gasmeter.NewGasMeter(t),
	}
}

type TestDriver struct {
	*StateDriver

	T                    testing.TB
	MessageProducer      *chain.MessageProducer
	TipSetMessageBuilder *TipSetMessageBuilder
	validator            *chain.Validator
	ExeCtx               *types.ExecutionContext

	Config state.ValidationConfig

	GasMeter *gasmeter.GasMeter
}

func (td *TestDriver) Complete() {
	//
	// Gas expectation recording.
	// Uncomment the following line to persist the actual gas values used to file as the new set
	// of expectations.
	//
	//td.GasMeter.Record()
}

//
// Unsigned Message Appliers
//

func (td *TestDriver) ApplyMessage(msg *types.Message) (result chain.ApplyResult) {
	defer func() {
		if r := recover(); r != nil {
			result.Receipt.ExitCode = exitcode.SysErrInternal
			td.T.Fatalf("message application panicked: %v", r)
		}
	}()

	result, err := td.validator.ApplyMessage(td.ExeCtx, td.State(), msg)
	require.NoError(td.T, err)
	return result
}

func (td *TestDriver) ApplyOk(msg *types.Message) chain.ApplyResult {
	return td.ApplyExpect(msg, EmptyReturnValue)
}

func (td *TestDriver) ApplyExpect(msg *types.Message, retval []byte) chain.ApplyResult {
	return td.applyMessageExpectCodeAndReturn(msg, exitcode.Ok, retval)
}

func (td *TestDriver) ApplyFailure(msg *types.Message, code exitcode.ExitCode) chain.ApplyResult {
	return td.applyMessageExpectCodeAndReturn(msg, code, EmptyReturnValue)
}

func (td *TestDriver) applyMessageExpectCodeAndReturn(msg *types.Message, code exitcode.ExitCode, retval []byte) chain.ApplyResult {
	result := td.ApplyMessage(msg)
	if !td.validateAndTrackResult(result, code, retval) {
		td.T.Logf("WARNING (not a test failure): failed to find expected gas cost for message: %+v", msg)
	}
	return result
}

//
// Signed Message Appliers
//

func (td *TestDriver) ApplyMessageSigned(msg *types.Message) (result chain.ApplyResult) {
	defer func() {
		if r := recover(); r != nil {
			result.Receipt.ExitCode = exitcode.SysErrInternal
			td.T.Fatalf("message application panicked: %v", r)
		}
	}()
	serMsg, err := msg.Serialize()
	require.NoError(td.T, err)

	msgSig, err := td.Wallet().Sign(msg.From, serMsg)
	require.NoError(td.T, err)

	smgs := &types.SignedMessage{
		Message:   *msg,
		Signature: msgSig,
	}
	result, err = td.validator.ApplySignedMessage(td.ExeCtx, td.State(), smgs)
	require.NoError(td.T, err)
	return result
}

func (td *TestDriver) ApplySignedOk(msg *types.Message) chain.ApplyResult {
	return td.ApplySignedExpect(msg, EmptyReturnValue)
}

func (td *TestDriver) ApplySignedExpect(msg *types.Message, retval []byte) chain.ApplyResult {
	return td.applyMessageSignedExpectCodeAndReturn(msg, exitcode.Ok, retval)
}

func (td *TestDriver) ApplySignedFailure(msg *types.Message, code exitcode.ExitCode) chain.ApplyResult {
	return td.applyMessageExpectCodeAndReturn(msg, code, EmptyReturnValue)
}

func (td *TestDriver) applyMessageSignedExpectCodeAndReturn(msg *types.Message, code exitcode.ExitCode, retval []byte) chain.ApplyResult {
	result := td.ApplyMessageSigned(msg)
	if !td.validateAndTrackResult(result, code, retval) {
		td.T.Logf("WARNING (not a test failure): failed to find expected gas cost for message: %+v", msg)
	}
	return result
}

func (td *TestDriver) validateAndTrackResult(result chain.ApplyResult, code exitcode.ExitCode, retval []byte) (foundGas bool) {
	foundGas = true

	td.GasMeter.TrackStateRoot(result.Root)
	td.GasMeter.TrackReceipt(result.Receipt)
	if td.Config.ValidateExitCode() {
		assert.Equal(td.T, code, result.Receipt.ExitCode, "Expected ExitCode: %s Actual ExitCode: %s", code.Error(), result.Receipt.ExitCode.Error())
	}
	if td.Config.ValidateReturnValue() {
		assert.Equal(td.T, retval, result.Receipt.ReturnValue, "Expected ReturnValue: %v Actual ReturnValue: %v", retval, result.Receipt.ReturnValue)
	}
	if td.Config.ValidateGas() {
		expectedGasUsed, ok := td.GasMeter.NextExpectedGas()
		if ok {
			assert.Equal(td.T, expectedGasUsed, result.Receipt.GasUsed, "Expected GasUsed: %d Actual GasUsed: %d", expectedGasUsed, result.Receipt.GasUsed)
		} else {
			foundGas = false
		}
	}
	if td.Config.ValidateStateRoot() {
		expectedRoot, found := td.GasMeter.NextExpectedStateRoot()
		actualRoot := td.State().Root()
		if found {
			assert.Equal(td.T, expectedRoot, actualRoot, "Expected StateRoot: %s Actual StateRoot: %s", expectedRoot, actualRoot)
		} else {
			td.T.Log("WARNING: failed to find expected state  root for message number")
		}
	}
	return
}

func (td *TestDriver) AssertNoActor(addr address.Address) {
	_, err := td.State().Actor(addr)
	assert.Error(td.T, err, "expected no such actor %s", addr)
}

func (td *TestDriver) GetBalance(addr address.Address) abi_spec.TokenAmount {
	actr, err := td.State().Actor(addr)
	require.NoError(td.T, err)
	return actr.Balance()
}

// AssertBalance checks an actor has an expected balance.
func (td *TestDriver) AssertBalance(addr address.Address, expected abi_spec.TokenAmount) {
	actr, err := td.State().Actor(addr)
	require.NoError(td.T, err)
	assert.Equal(td.T, expected, actr.Balance(), fmt.Sprintf("expected actor %s balance: %s, actual balance: %s", addr, expected, actr.Balance()))
}

// Checks an actor's balance and callSeqNum.
func (td *TestDriver) AssertActor(addr address.Address, balance abi_spec.TokenAmount, callSeqNum uint64) {
	actr, err := td.State().Actor(addr)
	require.NoError(td.T, err)
	assert.Equal(td.T, balance, actr.Balance(), fmt.Sprintf("expected actor %s balance: %s, actual balance: %s", addr, balance, actr.Balance()))
	assert.Equal(td.T, callSeqNum, actr.CallSeqNum(), fmt.Sprintf("expected actor %s callSeqNum: %d, actual : %d", addr, callSeqNum, actr.CallSeqNum()))
}

func (td *TestDriver) AssertBalanceCallback(addr address.Address, thing func(actorBalance abi_spec.TokenAmount) bool) {
	actr, err := td.State().Actor(addr)
	require.NoError(td.T, err)
	assert.True(td.T, thing(actr.Balance()))
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

func (td *TestDriver) ComputeInitActorExecReturn(from address.Address, originatorCallSeq uint64, newActorAddressCount uint64, expectedNewAddr address.Address) init_spec.ExecReturn {
	td.T.Helper()
	return computeInitActorExecReturn(td.T, from, originatorCallSeq, newActorAddressCount, expectedNewAddr)
}

func computeInitActorExecReturn(t testing.TB, from address.Address, originatorCallSeq uint64, newActorAddressCount uint64, expectedNewAddr address.Address) init_spec.ExecReturn {
	t.Helper()
	buf := new(bytes.Buffer)
	if from.Protocol() == address.ID {
		t.Fatal("cannot compute init actor address return from ID address", from)
	}

	require.NoError(t, from.MarshalCBOR(buf))
	require.NoError(t, binary.Write(buf, binary.BigEndian, originatorCallSeq))
	require.NoError(t, binary.Write(buf, binary.BigEndian, newActorAddressCount))

	out, err := address.NewActorAddress(buf.Bytes())
	require.NoError(t, err)

	return init_spec.ExecReturn{
		IDAddress:     expectedNewAddr,
		RobustAddress: out,
	}
}

func (td *TestDriver) MustCreateAndVerifyMultisigActor(nonce uint64, value abi_spec.TokenAmount, multisigAddr address.Address, from address.Address, params *multisig_spec.ConstructorParams, code exitcode.ExitCode, retval []byte) {
	/* Create the Multisig actor*/
	td.applyMessageExpectCodeAndReturn(
		td.MessageProducer.CreateMultisigActor(from, params.Signers, params.UnlockDuration, params.NumApprovalsThreshold, chain.Nonce(nonce), chain.Value(value)),
		code, retval)
	/* Assert the actor state was setup as expected */
	pendingTxMap, err := adt_spec.MakeEmptyMap(newMockStore())
	require.NoError(td.T, err)
	initialBalance := big_spec.Zero()
	startEpoch := abi_spec.ChainEpoch(0)
	if params.UnlockDuration > 0 {
		initialBalance = value
		startEpoch = td.ExeCtx.Epoch
	}
	td.AssertMultisigState(multisigAddr, multisig_spec.State{
		NextTxnID:      0,
		InitialBalance: initialBalance,
		StartEpoch:     startEpoch,

		Signers:               params.Signers,
		UnlockDuration:        params.UnlockDuration,
		NumApprovalsThreshold: params.NumApprovalsThreshold,

		PendingTxns: pendingTxMap.Root(),
	})
	td.AssertBalance(multisigAddr, value)
}

type RewardSummary struct {
	Treasury    abi_spec.TokenAmount
	RewardTotal abi_spec.TokenAmount
	Rewards     map[address.Address]abi_spec.TokenAmount
}

func (r *RewardSummary) For(a address.Address) abi_spec.TokenAmount {
	v, ok := r.Rewards[a]
	if !ok {
		return big_spec.Zero()
	}
	return v
}

func (td *TestDriver) GetRewardSummary() *RewardSummary {
	var rst reward_spec.State
	td.GetActorState(builtin_spec.RewardActorAddr, &rst)
	rewards := make(map[address.Address]abi_spec.TokenAmount)
	// Traverse map keyed by miner address.
	var r cbg.CborCid
	err := adt_spec.AsMap(td.State().Store(), rst.RewardMap).ForEach(&r, func(key string) error {
		keyAddr, err := address.NewFromBytes([]byte(key))
		require.NoError(td.T, err)

		// Traverse array of reward entries.
		sum := big_spec.Zero()
		var rw reward_spec.Reward
		err = adt_spec.AsArray(td.State().Store(), cid.Cid(r)).ForEach(&rw, func(i int64) error {
			sum = big_spec.Sub(big_spec.Add(sum, rw.Value), rw.AmountWithdrawn)
			return nil
		})
		require.NoError(td.T, err)

		rewards[keyAddr] = sum
		return nil
	})
	require.NoError(td.T, err)
	return &RewardSummary{
		Treasury:    td.GetBalance(builtin_spec.RewardActorAddr),
		RewardTotal: rst.RewardTotal,
		Rewards:     rewards,
	}
}

package drivers

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/exitcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain/types"
)

type TipSetMessageBuilder struct {
	driver *TestDriver

	bbs []*BlockBuilder
}

func NewTipSetMessageBuilder(testDriver *TestDriver) *TipSetMessageBuilder {
	return &TipSetMessageBuilder{
		driver: testDriver,
		bbs:    nil,
	}
}

func (t *TipSetMessageBuilder) WithBlockBuilder(bb *BlockBuilder) *TipSetMessageBuilder {
	t.bbs = append(t.bbs, bb)
	return t
}

func (t *TipSetMessageBuilder) Apply() types.ApplyTipSetResult {
	result := t.apply()
	t.validateState(result)

	t.Clear()
	return result
}

func (t *TipSetMessageBuilder) ApplyAndValidate() types.ApplyTipSetResult {
	result := t.apply()

	t.validateResult(result)
	t.validateState(result)

	t.Clear()
	return result
}

func (t *TipSetMessageBuilder) apply() types.ApplyTipSetResult {
	var blks []types.BlockMessagesInfo
	for _, b := range t.bbs {
		blks = append(blks, b.build())
	}
	result, err := t.driver.validator.ApplyTipSetMessages(t.driver.ExeCtx.Epoch, blks, t.driver.Randomness())
	require.NoError(t.driver.T, err)

	t.driver.StateTracker.TrackResult(result)
	return result
}

func (t *TipSetMessageBuilder) validateResult(result types.ApplyTipSetResult) {
	expected := []ExpectedResult{}
	for _, b := range t.bbs {
		expected = append(expected, b.expectedResults...)
	}

	if len(result.Receipts) > len(expected) {
		t.driver.T.Fatalf("ApplyTipSetMessages returned more result than expected. Expected: %d, Actual: %d", len(expected), len(result.Receipts))
		return
	}

	for i := range result.Receipts {
		if t.driver.Config.ValidateExitCode() {
			assert.Equal(t.driver.T, expected[i].ExitCode, result.Receipts[i].ExitCode, "Message Number: %d Expected ExitCode: %s Actual ExitCode: %s", i, expected[i].ExitCode.Error(), result.Receipts[i].ExitCode.Error())
		}
		if t.driver.Config.ValidateReturnValue() {
			assert.Equal(t.driver.T, expected[i].ReturnVal, result.Receipts[i].ReturnValue, "Message Number: %d Expected ReturnValue: %v Actual ReturnValue: %v", i, expected[i].ReturnVal, result.Receipts[i].ReturnValue)
		}
	}
}

func (t *TipSetMessageBuilder) validateState(result types.ApplyTipSetResult) {
	if t.driver.Config.ValidateGas() {
		for i := range result.Receipts {
			expectedGas, found := t.driver.StateTracker.NextExpectedGas()
			if found {
				assert.Equal(t.driver.T, expectedGas, result.Receipts[i].GasUsed, "Message Number: %d Expected GasUsed: %d Actual GasUsed: %d", i, expectedGas, result.Receipts[i].GasUsed)
			} else {
				t.driver.T.Logf("WARNING: failed to find expected gas cost for message number: %d", i)
			}
		}
	}
	if t.driver.Config.ValidateStateRoot() {
		expectedRoot, found := t.driver.StateTracker.NextExpectedStateRoot()
		actualRoot := t.driver.State().Root()
		if found {
			assert.Equal(t.driver.T, expectedRoot, actualRoot, "Expected StateRoot: %s Actual StateRoot: %s", expectedRoot, actualRoot)
		} else {
			t.driver.T.Log("WARNING: failed to find expected state  root for message number")
		}
	}
}

func (t *TipSetMessageBuilder) Clear() {
	t.bbs = nil
}

type BlockBuilder struct {
	TD *TestDriver

	miner       address.Address
	ticketCount int64

	secpMsgs []*types.SignedMessage
	blsMsgs  []*types.Message

	expectedResults []ExpectedResult
}

type ExpectedResult struct {
	ExitCode  exitcode.ExitCode
	ReturnVal []byte
}

func NewBlockBuilder(td *TestDriver, miner address.Address) *BlockBuilder {
	return &BlockBuilder{
		TD:              td,
		miner:           miner,
		ticketCount:     1,
		secpMsgs:        nil,
		blsMsgs:         nil,
		expectedResults: nil,
	}
}

func (bb *BlockBuilder) addResult(code exitcode.ExitCode, retval []byte) {
	bb.expectedResults = append(bb.expectedResults, ExpectedResult{
		ExitCode:  code,
		ReturnVal: retval,
	})
}

func (bb *BlockBuilder) WithBLSMessageOk(blsMsg *types.Message) *BlockBuilder {
	bb.blsMsgs = append(bb.blsMsgs, blsMsg)
	bb.addResult(exitcode.Ok, EmptyReturnValue)
	return bb
}

func (bb *BlockBuilder) WithBLSMessageDropped(blsMsg *types.Message) *BlockBuilder {
	bb.blsMsgs = append(bb.blsMsgs, blsMsg)
	return bb
}

func (bb *BlockBuilder) WithBLSMessageAndCode(bm *types.Message, code exitcode.ExitCode) *BlockBuilder {
	bb.blsMsgs = append(bb.blsMsgs, bm)
	bb.addResult(code, EmptyReturnValue)
	return bb
}

func (bb *BlockBuilder) WithBLSMessageAndRet(bm *types.Message, retval []byte) *BlockBuilder {
	bb.blsMsgs = append(bb.blsMsgs, bm)
	bb.addResult(exitcode.Ok, retval)
	return bb
}

func (bb *BlockBuilder) WithSECPMessageAndCode(bm *types.Message, code exitcode.ExitCode) *BlockBuilder {
	secpMsg := bb.toSignedMessage(bm)
	bb.secpMsgs = append(bb.secpMsgs, secpMsg)
	bb.addResult(code, EmptyReturnValue)
	return bb
}

func (bb *BlockBuilder) WithSECPMessageAndRet(bm *types.Message, retval []byte) *BlockBuilder {
	secpMsg := bb.toSignedMessage(bm)
	bb.secpMsgs = append(bb.secpMsgs, secpMsg)
	bb.addResult(exitcode.Ok, retval)
	return bb
}

func (bb *BlockBuilder) WithSECPMessageOk(bm *types.Message) *BlockBuilder {
	secpMsg := bb.toSignedMessage(bm)
	bb.secpMsgs = append(bb.secpMsgs, secpMsg)
	bb.addResult(exitcode.Ok, EmptyReturnValue)
	return bb
}

func (bb *BlockBuilder) WithSECPMessageDropped(bm *types.Message) *BlockBuilder {
	secpMsg := bb.toSignedMessage(bm)
	bb.secpMsgs = append(bb.secpMsgs, secpMsg)
	return bb
}

func (bb *BlockBuilder) WithTicketCount(count int64) *BlockBuilder {
	bb.ticketCount = count
	return bb
}

func (bb *BlockBuilder) toSignedMessage(m *types.Message) *types.SignedMessage {
	from := m.From
	if from.Protocol() == address.ID {
		from = bb.TD.ActorPubKey(from)
	}
	if from.Protocol() != address.SECP256K1 {
		bb.TD.T.Fatalf("Invalid address for SECP signature, address protocol: %v", from.Protocol())
	}
	raw, err := m.Serialize()
	require.NoError(bb.TD.T, err)

	sig, err := bb.TD.Wallet().Sign(from, raw)
	require.NoError(bb.TD.T, err)

	return &types.SignedMessage{
		Message:   *m,
		Signature: sig,
	}
}

func (bb *BlockBuilder) build() types.BlockMessagesInfo {
	return types.BlockMessagesInfo{
		BLSMessages:  bb.blsMsgs,
		SECPMessages: bb.secpMsgs,
		Miner:        bb.miner,
		TicketCount:  bb.ticketCount,
	}
}

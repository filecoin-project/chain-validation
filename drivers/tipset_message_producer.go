package drivers

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
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

func (t *TipSetMessageBuilder) Apply() []types.MessageReceipt {
	var blks []types.BlockMessagesInfo
	for _, b := range t.bbs {
		blks = append(blks, b.build())
	}
	receipts, err := t.driver.validator.ApplyTipSetMessages(t.driver.ExeCtx.Epoch, t.driver.State(), blks, t.driver.Randomness())
	require.NoError(t.driver.T, err)

	return receipts
}

func (t *TipSetMessageBuilder) ApplyAndValidate() {
	receipts := t.Apply()

	var results []Result
	for _, b := range t.bbs {
		results = append(results, b.expectedResults...)
	}

	if len(receipts) > len(results) {
		t.driver.T.Fatalf("ApplyTipSetMessages returned more receipts than expected. Expected: %d, Actual: %d", len(results), len(receipts))
	}

	for i := range receipts {
		t.driver.GasMeter.Track(receipts[i])
		if t.driver.Config.ValidateExitCode() {
			assert.Equal(t.driver.T, results[i].ExitCode, receipts[i].ExitCode, "Message Number: %d Expected ExitCode: %s Actual ExitCode: %s", i, results[i].ExitCode.Error(), receipts[i].ExitCode.Error())
		}
		if t.driver.Config.ValidateReturnValue() {
			assert.Equal(t.driver.T, results[i].ReturnVal, receipts[i].ReturnValue, "Message Number: %d Expected ReturnValue: %v Actual ReturnValue: %v", i, results[i].ReturnVal, receipts[i].ReturnValue)
		}
		if t.driver.Config.ValidateGas() {
			expectedGas, found := t.driver.GasMeter.NextExpectedGas()
			if found {
				assert.Equal(t.driver.T, expectedGas, receipts[i].GasUsed.Int64(), "Message Number: %d Expected GasUsed: %s Actual GasUsed: %s", i, expectedGas, receipts[i].GasUsed.String())
			} else {
				t.driver.T.Logf("WARNING: failed to find expected gas cost for message number: %d", i)
			}
		}
	}
	t.Clear()
}

func (t *TipSetMessageBuilder) Clear() {
	t.bbs = nil
}

type BlockBuilder struct {
	miner       address.Address
	ticketCount int64

	secpMsgs []*types.SignedMessage
	blsMsgs  []*types.Message

	expectedResults []Result
}

type Result struct {
	ExitCode  exitcode.ExitCode
	ReturnVal []byte
}

func NewBlockBuilder(miner address.Address) *BlockBuilder {
	return &BlockBuilder{
		miner:           miner,
		ticketCount:     0,
		secpMsgs:        nil,
		blsMsgs:         nil,
		expectedResults: nil,
	}
}

func (bb *BlockBuilder) addResult(code exitcode.ExitCode, retval []byte) {
	bb.expectedResults = append(bb.expectedResults, Result{
		ExitCode:  code,
		ReturnVal: retval,
	})
}

func (bb *BlockBuilder) WithSECPMessageOk(secpMsg *types.SignedMessage) *BlockBuilder {
	bb.secpMsgs = append(bb.secpMsgs, secpMsg)
	bb.addResult(exitcode.Ok, EmptyReturnValue)
	return bb
}

func (bb *BlockBuilder) WithBLSMessageOk(blsMsg *types.Message) *BlockBuilder {
	bb.blsMsgs = append(bb.blsMsgs, blsMsg)
	bb.addResult(exitcode.Ok, EmptyReturnValue)
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

func (bb *BlockBuilder) WithSECPMessageAndCode(sm *types.SignedMessage, code exitcode.ExitCode) *BlockBuilder {
	bb.secpMsgs = append(bb.secpMsgs, sm)
	bb.addResult(code, EmptyReturnValue)
	return bb
}

func (bb *BlockBuilder) WithSECPMessageAndRet(sm *types.SignedMessage, retval []byte) *BlockBuilder {
	bb.secpMsgs = append(bb.secpMsgs, sm)
	bb.addResult(exitcode.Ok, retval)
	return bb
}

func (bb *BlockBuilder) WithResult(code exitcode.ExitCode, retval []byte) *BlockBuilder {
	bb.addResult(code, retval)
	return bb
}

func (bb *BlockBuilder) WithTicketCount(count int64) *BlockBuilder {
	bb.ticketCount = count
	return bb
}

func (bb *BlockBuilder) build() types.BlockMessagesInfo {
	return types.BlockMessagesInfo{
		BLSMessages:  bb.blsMsgs,
		SECPMessages: bb.secpMsgs,
		Miner:        bb.miner,
		TicketCount:  bb.ticketCount,
	}
}

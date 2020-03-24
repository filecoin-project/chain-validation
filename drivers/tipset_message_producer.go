package drivers

import (
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain/types"
)

type TipSetMessageBuilder struct {
	driver *TestDriver

	secpMsgs []*types.SignedMessage
	blsMsgs  []*types.Message

	expectedResults []Result
	ticketCount     int64
}

type Result struct {
	ExitCode  exitcode.ExitCode
	ReturnVal []byte
}

func NewTipSetMessageBuilder(testDriver *TestDriver) *TipSetMessageBuilder {
	return &TipSetMessageBuilder{
		driver:          testDriver,
		ticketCount:     0,
		secpMsgs:        nil,
		blsMsgs:         nil,
		expectedResults: nil,
	}
}

func (t *TipSetMessageBuilder) addResult(code exitcode.ExitCode, retval []byte) {
	t.expectedResults = append(t.expectedResults, Result{
		ExitCode:  code,
		ReturnVal: retval,
	})
}

func (t *TipSetMessageBuilder) WithSECPMessageOk(secpMsg *types.SignedMessage) *TipSetMessageBuilder {
	t.secpMsgs = append(t.secpMsgs, secpMsg)
	t.addResult(exitcode.Ok, EmptyReturnValue)
	return t
}

func (t *TipSetMessageBuilder) WithBLSMessageOk(blsMsg *types.Message) *TipSetMessageBuilder {
	t.blsMsgs = append(t.blsMsgs, blsMsg)
	t.addResult(exitcode.Ok, EmptyReturnValue)
	return t
}

func (t *TipSetMessageBuilder) WithBLSMessageAndCode(bm *types.Message, code exitcode.ExitCode) *TipSetMessageBuilder {
	t.blsMsgs = append(t.blsMsgs, bm)
	t.addResult(code, EmptyReturnValue)
	return t
}

func (t *TipSetMessageBuilder) WithBLSMessageAndRet(bm *types.Message, retval []byte) *TipSetMessageBuilder {
	t.blsMsgs = append(t.blsMsgs, bm)
	t.addResult(exitcode.Ok, retval)
	return t
}

func (t *TipSetMessageBuilder) WithSECPMessageAndCode(sm *types.SignedMessage, code exitcode.ExitCode) *TipSetMessageBuilder {
	t.secpMsgs = append(t.secpMsgs, sm)
	t.addResult(code, EmptyReturnValue)
	return t
}

func (t *TipSetMessageBuilder) WithSECPMessageAndRet(sm *types.SignedMessage, retval []byte) *TipSetMessageBuilder {
	t.secpMsgs = append(t.secpMsgs, sm)
	t.addResult(exitcode.Ok, retval)
	return t
}

func (t *TipSetMessageBuilder) WithResult(code exitcode.ExitCode, retval []byte) *TipSetMessageBuilder {
	t.addResult(code, retval)
	return t
}

func (t *TipSetMessageBuilder) WithTicketCount(count int64) *TipSetMessageBuilder {
	t.ticketCount = count
	return t
}

func (t *TipSetMessageBuilder) Build() types.BlockMessagesInfo {
	return types.BlockMessagesInfo{
		BLSMessages:  t.blsMsgs,
		SECPMessages: t.secpMsgs,
		TicketCount:  t.ticketCount,
		Miner:        t.driver.ExeCtx.Miner,
	}
}

func (t *TipSetMessageBuilder) Apply() []types.MessageReceipt {
	receipts, err := t.driver.validator.ApplyTipSetMessages(t.driver.ExeCtx, t.driver.State(), []types.BlockMessagesInfo{t.Build()}, t.driver.Randomness())
	require.NoError(t.driver.T, err)

	return receipts
}

func (t *TipSetMessageBuilder) ApplyAndValidate() {
	receipts := t.Apply()
	if len(receipts) > len(t.expectedResults) {
		t.driver.T.Fatalf("ApplyTipSetMessages returned more receipts than expected. Expected: %d, Actual: %d", len(t.expectedResults), len(receipts))
	}
	for i := range receipts {
		t.driver.GasMeter.Track(receipts[i])
		if t.driver.Config.ValidateExitCode() {
			assert.Equal(t.driver.T, t.expectedResults[i].ExitCode, receipts[i].ExitCode, "Message Number: %d Expected ExitCode: %s Actual ExitCode: %s", i, t.expectedResults[i].ExitCode.Error(), receipts[i].ExitCode.Error())
		}
		if t.driver.Config.ValidateReturnValue() {
			assert.Equal(t.driver.T, t.expectedResults[i].ReturnVal, receipts[i].ReturnValue, "Message Number: %d Expected ReturnValue: %v Actual ReturnValue: %v", i, t.expectedResults[i].ReturnVal, receipts[i].ReturnValue)
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
	t.expectedResults = nil
	t.secpMsgs = nil
	t.blsMsgs = nil
	t.ticketCount = 0
}

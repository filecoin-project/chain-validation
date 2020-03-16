package drivers

import (
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/chain/types"
)

type TipSetMessageBuilder struct {
	driver *TestDriver

	secpMsgs []*types.SignedMessage
	blsMsgs  []*types.Message

	msgReceipts []types.MessageReceipt
	ticketCount int64
}

func NewTipSetMessageBuilder(testDriver *TestDriver) *TipSetMessageBuilder {
	return &TipSetMessageBuilder{
		driver:      testDriver,
		ticketCount: 0,
		secpMsgs:    nil,
		blsMsgs:     nil,
		msgReceipts: nil,
	}
}

func (t *TipSetMessageBuilder) WithSECPMessage(secpMsg *types.SignedMessage) *TipSetMessageBuilder {
	t.secpMsgs = append(t.secpMsgs, secpMsg)
	return t
}

func (t *TipSetMessageBuilder) WithBLSMessage(blsMsg *types.Message) *TipSetMessageBuilder {
	t.blsMsgs = append(t.blsMsgs, blsMsg)
	return t
}

func (t *TipSetMessageBuilder) WithBLSMessageAndReceipt(bm *types.Message, rc types.MessageReceipt) *TipSetMessageBuilder {
	t.blsMsgs = append(t.blsMsgs, bm)
	t.msgReceipts = append(t.msgReceipts, rc)
	return t
}

func (t *TipSetMessageBuilder) WithSECPMessageAndReceipt(sm *types.SignedMessage, rc types.MessageReceipt) *TipSetMessageBuilder {
	t.secpMsgs = append(t.secpMsgs, sm)
	t.msgReceipts = append(t.msgReceipts, rc)
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
	receipts, err := t.driver.Validator.ApplyTipSetMessages(t.driver.ExeCtx, t.driver.State(), []types.BlockMessagesInfo{t.Build()}, t.driver.Randomness())
	require.NoError(t.driver.T, err)
	return receipts
}

func (t *TipSetMessageBuilder) ApplyAndValidate() {
	receipts := t.Apply()
	for i := range receipts {
		t.driver.AssertReceipt(receipts[i], t.msgReceipts[i])
	}
	t.Clear()
}

func (t *TipSetMessageBuilder) Clear() {
	t.msgReceipts = nil
	t.secpMsgs = nil
	t.blsMsgs = nil
	t.ticketCount = 0
}

package chain

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/chain-validation/chain/types"
)

type TipSetMessageBuilder struct {
	miner address.Address

	secpMsgs []*types.SignedMessage
	blsMsgs  []*types.Message
}

func NewTipSetMessageBuilder() *TipSetMessageBuilder {
	return &TipSetMessageBuilder{}
}

func (t *TipSetMessageBuilder) WithMiner(miner address.Address) *TipSetMessageBuilder {
	t.miner = miner
	return t
}

func (t *TipSetMessageBuilder) WithSECPMessage(secpMsg *types.SignedMessage) *TipSetMessageBuilder {
	t.secpMsgs = append(t.secpMsgs, secpMsg)
	return t
}

func (t *TipSetMessageBuilder) WithBLSMessage(blsMsg *types.Message) *TipSetMessageBuilder {
	t.blsMsgs = append(t.blsMsgs, blsMsg)
	return t
}

func (t *TipSetMessageBuilder) Build() types.BlockMessagesInfo {
	return types.BlockMessagesInfo{
		BLSMessages:  t.blsMsgs,
		SECPMessages: t.secpMsgs,
		Miner:        t.miner,
	}
}
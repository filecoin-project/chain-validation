package types

import "github.com/filecoin-project/go-address"

// BlockMessagesInfo contains messages for one block in a tipset.
type BlockMessagesInfo struct {
	BLSMessages  []*Message
	SECPMessages []*SignedMessage
	Miner        address.Address
	TicketCount  int64
}

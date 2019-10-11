package chainvalidation

import (
	"github.com/ipfs/go-cid"
)

var AccountActorCodeCid cid.Cid
var StorageMarketActorCodeCid cid.Cid
var StorageMinerCodeCid cid.Cid
var MultisigActorCodeCid cid.Cid
var InitActorCodeCid cid.Cid
var PaymentChannelActorCodeCid cid.Cid

type Actor interface {
	Code() cid.Cid
	Head() cid.Cid
	Nonce() uint64
	Balance() uint64
}

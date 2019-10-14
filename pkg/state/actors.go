package state

import (
	"github.com/ipfs/go-cid"
	dag "github.com/ipfs/go-merkledag"
)

// Builtin actor code CIDs.
var (
	InitActorCodeCid          cid.Cid
	AccountActorCodeCid       cid.Cid
	StorageMarketActorCodeCid cid.Cid
)

// Builtin actor addresses.
var (
	InitAddress          Address
	NetworkAddress       Address
	StorageMarketAddress Address
	BurntFundsAddress    Address
)

func init() {
	var err error

	InitActorCodeObj := dag.NewRawNode([]byte("initactor"))
	InitActorCodeCid = InitActorCodeObj.Cid()

	AccountActorCodeObj := dag.NewRawNode([]byte("accountactor"))
	AccountActorCodeCid = AccountActorCodeObj.Cid()

	StorageMarketActorCodeObj := dag.NewRawNode([]byte("storagemarket"))
	StorageMarketActorCodeCid = StorageMarketActorCodeObj.Cid()

	InitAddress, err = NewIDAddress(0)
	if err != nil {
		panic(err)
	}
	NetworkAddress, err = NewIDAddress(1)
	if err != nil {
		panic(err)
	}
	StorageMarketAddress, err = NewIDAddress(2)
	if err != nil {
		panic(err)
	}
	BurntFundsAddress, err = NewIDAddress(99)
	if err != nil {
		panic(err)
	}
}

package state

import (
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
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
	pref := cid.NewPrefixV1(cid.Raw, mh.ID)
	mustSum := func(s string) cid.Cid {
		c, err := pref.Sum([]byte(s))
		if err != nil {
			panic(err)
		}
		return c
	}

	AccountActorCodeCid = mustSum("account")
	StorageMarketActorCodeCid = mustSum("smarket")
	InitActorCodeCid = mustSum("init")

	var err error
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


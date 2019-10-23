package state

type ActorCodeCid int
const (
	AccountActorCodeCid  = ActorCodeCid(iota)
	StorageMinerCodeCid
	MultisigActorCodeCid
	PaymentChannelActorCodeCid
)

type SingletonActorAddress int
const (
	InitAddress = SingletonActorAddress(iota)
	NetworkAddress
	StorageMarketAddress
	BurntFundsAddress
)

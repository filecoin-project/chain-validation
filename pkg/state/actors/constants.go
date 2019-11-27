package actors

type ActorCodeID int

const (
	AccountActorCodeCid = ActorCodeID(iota)
	StorageMinerCodeCid
	MultisigActorCodeCid
	PaymentChannelActorCodeCid
)

type SingletonActorID int

const (
	InitAddress = SingletonActorID(iota)
	NetworkAddress
	StorageMarketAddress
	BurntFundsAddress
	StoragePowerAddress
)

package chainvalidation

import (
	"github.com/filecoin-project/chain-validation/address"
)

type StateTree interface {
	GetActor(address.Address) (Actor, error)
}

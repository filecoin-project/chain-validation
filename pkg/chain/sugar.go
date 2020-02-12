package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"

	"github.com/filecoin-project/chain-validation/pkg/chain/types"
)

var noParams []byte

// Transfer builds a simple value transfer message and returns it.
func (mp *MessageProducer) Transfer(to, from address.Address, opts ...MsgOpt) *types.Message {
	return mp.Build(to, from, builtin_spec.MethodSend, noParams, opts...)
}

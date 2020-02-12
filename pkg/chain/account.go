package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

func (mp *MessageProducer) AccountConstructor(to, from address.Address, params address.Address, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsAccount.Constructor, ser, opts...)
}
func (mp *MessageProducer) AccountPubkeyAddress(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsAccount.PubkeyAddress, ser, opts...)
}

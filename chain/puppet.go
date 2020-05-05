package chain

import (
	"github.com/filecoin-project/go-address"
	puppet "github.com/filecoin-project/specs-actors/actors/puppet"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func (mp *MessageProducer) PuppetConstructor(to, from address.Address, params *adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, puppet.MethodsPuppet.Constructor, ser, opts...)
}
func (mp *MessageProducer) PuppetSend(to, from address.Address, params *puppet.SendParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, puppet.MethodsPuppet.Send, ser, opts...)
}

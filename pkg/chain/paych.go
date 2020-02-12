package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"
	"github.com/filecoin-project/specs-actors/actors/util/adt"
)

func (mp *MessageProducer) PaychConstructor(to, from address.Address, params paych.ConstructorParams, opts ...MsgOpt) *Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPaych.Constructor, ser, opts...)
}
func (mp *MessageProducer) PaychUpdateChannelState(to, from address.Address, params paych.UpdateChannelStateParams, opts ...MsgOpt) *Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPaych.UpdateChannelState, ser, opts...)
}
func (mp *MessageProducer) PaychSettle(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPaych.Settle, ser, opts...)
}
func (mp *MessageProducer) PaychCollect(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPaych.Collect, ser, opts...)
}

package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	init_ "github.com/filecoin-project/specs-actors/actors/builtin/init"
	"github.com/filecoin-project/specs-actors/actors/util/adt"
)

func (mp *MessageProducer) InitConstructor(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsInit.Constructor, ser, opts...)
}
func (mp *MessageProducer) InitExec(to, from address.Address, params init_.ExecParams, opts ...MsgOpt) *Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsInit.Exec, ser, opts...)
}

package chain

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func (mp *MessageProducer) PowerConstructor(from, to address.Address, params *adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.Constructor, ser, opts...)
}
func (mp *MessageProducer) PowerCreateMiner(from, to address.Address, params *power.CreateMinerParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.CreateMiner, ser, opts...)
}
func (mp *MessageProducer) PowerDeleteMiner(from, to address.Address, params *power.DeleteMinerParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.DeleteMiner, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorProveCommit(from, to address.Address, params *power.OnSectorProveCommitParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.OnSectorProveCommit, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorTerminate(from, to address.Address, params *power.OnSectorTerminateParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.OnSectorTerminate, ser, opts...)
}
func (mp *MessageProducer) PowerOnFaultBegin(from, to address.Address, params *power.OnFaultBeginParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.OnFaultBegin, ser, opts...)
}
func (mp *MessageProducer) PowerOnFaultEnd(from, to address.Address, params *power.OnFaultEndParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.OnFaultEnd, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorModifyWeightDesc(from, to address.Address, params *power.OnSectorModifyWeightDescParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.OnSectorModifyWeightDesc, ser, opts...)
}
func (mp *MessageProducer) PowerEnrollCronEvent(from, to address.Address, params *power.EnrollCronEventParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.EnrollCronEvent, ser, opts...)
}
func (mp *MessageProducer) PowerOnEpochTickEnd(from, to address.Address, params *adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.OnEpochTickEnd, ser, opts...)
}
func (mp *MessageProducer) PowerUpdatePledgeTotal(from, to address.Address, params *big.Int, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.UpdatePledgeTotal, ser, opts...)
}
func (mp *MessageProducer) PowerOnConsensusFault(from, to address.Address, params *big.Int, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsPower.OnConsensusFault, ser, opts...)
}

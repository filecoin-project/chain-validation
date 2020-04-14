package chain

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func (mp *MessageProducer) PowerConstructor(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.Constructor, ser, opts...)
}
func (mp *MessageProducer) PowerCreateMiner(to, from address.Address, params power.CreateMinerParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.CreateMiner, ser, opts...)
}
func (mp *MessageProducer) PowerDeleteMiner(to, from address.Address, params power.DeleteMinerParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.DeleteMiner, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorProveCommit(to, from address.Address, params power.OnSectorProveCommitParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorProveCommit, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorTerminate(to, from address.Address, params power.OnSectorTerminateParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorTerminate, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorTemporaryFaultEffectiveBegin(to, from address.Address, params power.OnSectorTemporaryFaultEffectiveBeginParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorTemporaryFaultEffectiveBegin, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorTemporaryFaultEffectiveEnd(to, from address.Address, params power.OnSectorTemporaryFaultEffectiveEndParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorTemporaryFaultEffectiveEnd, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorModifyWeightDesc(to, from address.Address, params power.OnSectorModifyWeightDescParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorModifyWeightDesc, ser, opts...)
}
func (mp *MessageProducer) PowerOnMinerWindowedPoStSuccess(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnMinerWindowedPoStSuccess, ser, opts...)
}
func (mp *MessageProducer) PowerOnMinerWindowedPoStFailure(to, from address.Address, params power.OnMinerWindowedPoStFailureParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnMinerWindowedPoStFailure, ser, opts...)
}
func (mp *MessageProducer) PowerEnrollCronEvent(to, from address.Address, params power.EnrollCronEventParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.EnrollCronEvent, ser, opts...)
}
func (mp *MessageProducer) PowerOnEpochTickEnd(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnEpochTickEnd, ser, opts...)
}
func (mp *MessageProducer) PowerUpdatePledgeTotal(to, from address.Address, params big.Int, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.UpdatePledgeTotal, ser, opts...)
}
func (mp *MessageProducer) PowerOnConsensusFault(to, from address.Address, params big.Int, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnConsensusFault, ser, opts...)
}

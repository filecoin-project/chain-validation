package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func (mp *MessageProducer) PowerConstructor(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.Constructor, ser, opts...)
}
func (mp *MessageProducer) PowerAddBalance(to, from address.Address, params power.AddBalanceParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.AddBalance, ser, opts...)
}
func (mp *MessageProducer) PowerWithdrawBalance(to, from address.Address, params power.WithdrawBalanceParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.WithdrawBalance, ser, opts...)
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
func (mp *MessageProducer) PowerReportConsensusFault(to, from address.Address, params power.ReportConsensusFaultParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.ReportConsensusFault, ser, opts...)
}
func (mp *MessageProducer) PowerOnEpochTickEnd(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnEpochTickEnd, ser, opts...)
}

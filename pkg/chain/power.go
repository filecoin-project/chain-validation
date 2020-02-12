package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

func (mp *MessageProducer) PowerConstructor(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.Constructor, ser, opts...)
}
func (mp *MessageProducer) PowerAddBalance(to, from address.Address, params power.AddBalanceParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.AddBalance, ser, opts...)
}
func (mp *MessageProducer) PowerWithdrawBalance(to, from address.Address, params power.WithdrawBalanceParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.WithdrawBalance, ser, opts...)
}
func (mp *MessageProducer) PowerCreateMiner(to, from address.Address, params power.CreateMinerParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.CreateMiner, ser, opts...)
}
func (mp *MessageProducer) PowerDeleteMiner(to, from address.Address, params power.DeleteMinerParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.DeleteMiner, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorProveCommit(to, from address.Address, params power.OnSectorProveCommitParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorProveCommit, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorTerminate(to, from address.Address, params power.OnSectorTerminateParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorTerminate, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorTemporaryFaultEffectiveBegin(to, from address.Address, params power.OnSectorTemporaryFaultEffectiveBeginParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorTemporaryFaultEffectiveBegin, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorTemporaryFaultEffectiveEnd(to, from address.Address, params power.OnSectorTemporaryFaultEffectiveEndParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorTemporaryFaultEffectiveEnd, ser, opts...)
}
func (mp *MessageProducer) PowerOnSectorModifyWeightDesc(to, from address.Address, params power.OnSectorModifyWeightDescParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorModifyWeightDesc, ser, opts...)
}
func (mp *MessageProducer) PowerOnMinerSurprisePoStSuccess(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnMinerSurprisePoStSuccess, ser, opts...)
}
func (mp *MessageProducer) PowerOnMinerSurprisePoStFailure(to, from address.Address, params power.OnMinerSurprisePoStFailureParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnMinerSurprisePoStFailure, ser, opts...)
}
func (mp *MessageProducer) PowerEnrollCronEvent(to, from address.Address, params power.EnrollCronEventParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.EnrollCronEvent, ser, opts...)
}
func (mp *MessageProducer) PowerReportConsensusFault(to, from address.Address, params power.ReportConsensusFaultParams, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.ReportConsensusFault, ser, opts...)
}
func (mp *MessageProducer) PowerOnEpochTickEnd(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *Message {
	ser := state.MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsPower.OnEpochTickEnd, ser, opts...)
}

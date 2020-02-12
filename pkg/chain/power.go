package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

func (mp *MessageProducer) PowerConstructor(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.Constructor, ser, opts...), nil
}

func (mp *MessageProducer) PowerAddBalance(to, from address.Address, params power.AddBalanceParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.AddBalance, ser, opts...), nil
}

func (mp *MessageProducer) PowerWithdrawBalance(to, from address.Address, params power.WithdrawBalanceParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.WithdrawBalance, ser, opts...), nil
}

func (mp *MessageProducer) PowerCreateMiner(to, from address.Address, params power.CreateMinerParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.CreateMiner, ser, opts...), nil
}

func (mp *MessageProducer) PowerDeleteMiner(to, from address.Address, params power.DeleteMinerParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.DeleteMiner, ser, opts...), nil
}

func (mp *MessageProducer) PowerOnSectorProveCommit(to, from address.Address, params power.OnSectorProveCommitParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorProveCommit, ser, opts...), nil
}

func (mp *MessageProducer) PowerOnSectorTerminate(to, from address.Address, params power.OnSectorTerminateParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorTerminate, ser, opts...), nil
}

func (mp *MessageProducer) PowerOnSectorTemporaryFaultEffectiveBegin(to, from address.Address, params power.OnSectorTemporaryFaultEffectiveBeginParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorTemporaryFaultEffectiveBegin, ser, opts...), nil
}

func (mp *MessageProducer) PowerOnSectorTemporaryFaultEffectiveEnd(to, from address.Address, params power.OnSectorTemporaryFaultEffectiveEndParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorTemporaryFaultEffectiveEnd, ser, opts...), nil
}

func (mp *MessageProducer) PowerOnSectorModifyWeightDesc(to, from address.Address, params power.OnSectorModifyWeightDescParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.OnSectorModifyWeightDesc, ser, opts...), nil
}

func (mp *MessageProducer) PowerOnMinerSurprisePoStSuccess(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.OnMinerSurprisePoStSuccess, ser, opts...), nil
}

func (mp *MessageProducer) PowerOnMinerSurprisePoStFailure(to, from address.Address, params power.OnMinerSurprisePoStFailureParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.OnMinerSurprisePoStFailure, ser, opts...), nil
}

func (mp *MessageProducer) PowerEnrollCronEvent(to, from address.Address, params power.EnrollCronEventParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.EnrollCronEvent, ser, opts...), nil
}

func (mp *MessageProducer) PowerReportConsensusFault(to, from address.Address, params power.ReportConsensusFaultParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.ReportConsensusFault, ser, opts...), nil
}

func (mp *MessageProducer) PowerOnEpochTickEnd(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsPower.OnEpochTickEnd, ser, opts...), nil
}

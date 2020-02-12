package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

func (mp *MessageProducer) MinerConstructor(to, from address.Address, params power.MinerConstructorParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.Constructor, ser, opts...), nil
}

func (mp *MessageProducer) MinerControlAddresses(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.ControlAddresses, ser, opts...), nil
}

func (mp *MessageProducer) MinerChangeWorkerAddress(to, from address.Address, params miner.ChangeWorkerAddressParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.ChangeWorkerAddress, ser, opts...), nil
}

func (mp *MessageProducer) MinerOnSurprisePoStChallenge(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.OnSurprisePoStChallenge, ser, opts...), nil
}

func (mp *MessageProducer) MinerSubmitSurprisePoStResponse(to, from address.Address, params miner.SubmitSurprisePoStResponseParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.SubmitSurprisePoStResponse, ser, opts...), nil
}

func (mp *MessageProducer) MinerOnDeleteMiner(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.OnDeleteMiner, ser, opts...), nil
}

func (mp *MessageProducer) MinerOnVerifiedElectionPoSt(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.OnVerifiedElectionPoSt, ser, opts...), nil
}

func (mp *MessageProducer) MinerPreCommitSector(to, from address.Address, params miner.PreCommitSectorParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.PreCommitSector, ser, opts...), nil
}

func (mp *MessageProducer) MinerProveCommitSector(to, from address.Address, params miner.ProveCommitSectorParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.ProveCommitSector, ser, opts...), nil
}

func (mp *MessageProducer) MinerExtendSectorExpiration(to, from address.Address, params miner.ExtendSectorExpirationParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.ExtendSectorExpiration, ser, opts...), nil
}

func (mp *MessageProducer) MinerTerminateSectors(to, from address.Address, params miner.TerminateSectorsParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.TerminateSectors, ser, opts...), nil
}

func (mp *MessageProducer) MinerDeclareTemporaryFaults(to, from address.Address, params miner.DeclareTemporaryFaultsParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.DeclareTemporaryFaults, ser, opts...), nil
}

func (mp *MessageProducer) MinerOnDeferredCronEvent(to, from address.Address, params miner.OnDeferredCronEventParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMiner.OnDeferredCronEvent, ser, opts...), nil
}

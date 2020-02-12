package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/reward"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

func (mp *MessageProducer) RewardConstructor(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsReward.Constructor, ser, opts...), nil
}

func (mp *MessageProducer) RewardAwardBlockReward(to, from address.Address, params reward.AwardBlockRewardParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsReward.AwardBlockReward, ser, opts...), nil
}

func (mp *MessageProducer) RewardWithdrawReward(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsReward.WithdrawReward, ser, opts...), nil
}

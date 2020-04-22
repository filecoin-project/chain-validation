package chain

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/reward"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func (mp *MessageProducer) RewardConstructor(to, from address.Address, params *adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsReward.Constructor, ser, opts...)
}
func (mp *MessageProducer) RewardAwardBlockReward(to, from address.Address, params reward.AwardBlockRewardParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsReward.AwardBlockReward, ser, opts...)
}
func (mp *MessageProducer) RewardLastPerEpochReward(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsReward.LastPerEpochReward, ser, opts...)
}
func (mp *MessageProducer) RewardUpdateNetworkKPI(to, from address.Address, params big.Int, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsReward.UpdateNetworkKPI, ser, opts...)
}

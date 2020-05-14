package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func (mp *MessageProducer) MarketConstructor(from, to address.Address, params *adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsMarket.Constructor, ser, opts...)
}
func (mp *MessageProducer) MarketAddBalance(from, to address.Address, params *address.Address, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsMarket.AddBalance, ser, opts...)
}
func (mp *MessageProducer) MarketWithdrawBalance(from, to address.Address, params *market.WithdrawBalanceParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsMarket.WithdrawBalance, ser, opts...)
}
func (mp *MessageProducer) MarketPublishStorageDeals(from, to address.Address, params *market.PublishStorageDealsParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsMarket.PublishStorageDeals, ser, opts...)
}
func (mp *MessageProducer) MarketVerifyDealsOnSectorProveCommit(from, to address.Address, params *market.VerifyDealsOnSectorProveCommitParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsMarket.VerifyDealsOnSectorProveCommit, ser, opts...)
}
func (mp *MessageProducer) MarketOnMinerSectorsTerminate(from, to address.Address, params *market.OnMinerSectorsTerminateParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsMarket.OnMinerSectorsTerminate, ser, opts...)
}
func (mp *MessageProducer) MarketComputeDataCommitment(from, to address.Address, params *market.ComputeDataCommitmentParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsMarket.ComputeDataCommitment, ser, opts...)
}
func (mp *MessageProducer) MarketCronTick(from, to address.Address, params *adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(from, to, builtin_spec.MethodsMarket.CronTick, ser, opts...)
}

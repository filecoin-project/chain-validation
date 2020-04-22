package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func (mp *MessageProducer) MarketConstructor(to, from address.Address, params *adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMarket.Constructor, ser, opts...)
}
func (mp *MessageProducer) MarketAddBalance(to, from address.Address, params address.Address, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMarket.AddBalance, ser, opts...)
}
func (mp *MessageProducer) MarketWithdrawBalance(to, from address.Address, params market.WithdrawBalanceParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMarket.WithdrawBalance, ser, opts...)
}
func (mp *MessageProducer) MarketHandleExpiredDeals(to, from address.Address, params market.HandleExpiredDealsParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMarket.HandleExpiredDeals, ser, opts...)
}
func (mp *MessageProducer) MarketPublishStorageDeals(to, from address.Address, params market.PublishStorageDealsParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMarket.PublishStorageDeals, ser, opts...)
}
func (mp *MessageProducer) MarketVerifyDealsOnSectorProveCommit(to, from address.Address, params market.VerifyDealsOnSectorProveCommitParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMarket.VerifyDealsOnSectorProveCommit, ser, opts...)
}
func (mp *MessageProducer) MarketOnMinerSectorsTerminate(to, from address.Address, params market.OnMinerSectorsTerminateParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMarket.OnMinerSectorsTerminate, ser, opts...)
}
func (mp *MessageProducer) MarketComputeDataCommitment(to, from address.Address, params market.ComputeDataCommitmentParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMarket.ComputeDataCommitment, ser, opts...)
}
func (mp *MessageProducer) MarketHandleInitTimeoutDeals(to, from address.Address, params market.HandleInitTimeoutDealsParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMarket.HandleInitTimeoutDeals, ser, opts...)
}

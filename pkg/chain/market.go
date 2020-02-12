package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

func (mp *MessageProducer) MarketConstructor(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMarket.Constructor, ser, opts...), nil
}
func (mp *MessageProducer) MarketAddBalance(to, from address.Address, params address.Address, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMarket.AddBalance, ser, opts...), nil
}
func (mp *MessageProducer) MarketWithdrawBalance(to, from address.Address, params market.WithdrawBalanceParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMarket.WithdrawBalance, ser, opts...), nil
}
func (mp *MessageProducer) MarketHandleExpiredDeals(to, from address.Address, params market.HandleExpiredDealsParams, opts ...MsgOpt) (*Message, error) {
	// FIXME params does not fulfill cbor marshal interface
	panic("TODO HandleExpiredDealsParams does not implement a CBOR marshaller")
	/*
		ser, err := state.Serialize(&params)
		if err != nil {
			return nil, err
		}
		return mp.Build(to, from, builtin_spec.MethodsMarket.HandleExpiredDeals, ser, opts...), nil
	*/
}
func (mp *MessageProducer) MarketPublishStorageDeals(to, from address.Address, params market.PublishStorageDealsParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMarket.PublishStorageDeals, ser, opts...), nil
}
func (mp *MessageProducer) MarketVerifyDealsOnSectorProveCommit(to, from address.Address, params market.VerifyDealsOnSectorProveCommitParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMarket.VerifyDealsOnSectorProveCommit, ser, opts...), nil
}
func (mp *MessageProducer) MarketOnMinerSectorsTerminate(to, from address.Address, params market.OnMinerSectorsTerminateParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMarket.OnMinerSectorsTerminate, ser, opts...), nil
}
func (mp *MessageProducer) MarketComputeDataCommitment(to, from address.Address, params market.ComputeDataCommitmentParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMarket.ComputeDataCommitment, ser, opts...), nil
}

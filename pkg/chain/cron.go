package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/cron"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

func (mp *MessageProducer) CronConstructor(to, from address.Address, params cron.ConstructorParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsCron.Constructor, ser, opts...), nil
}
func (mp *MessageProducer) CronEpochTick(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsCron.EpochTick, ser, opts...), nil
}

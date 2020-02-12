package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	multisig "github.com/filecoin-project/specs-actors/actors/builtin/multisig"

	"github.com/filecoin-project/chain-validation/pkg/state"
)

func (mp *MessageProducer) MultisigConstructor(to, from address.Address, params multisig.ConstructorParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Constructor, ser, opts...), nil
}
func (mp *MessageProducer) MultisigPropose(to, from address.Address, params multisig.ProposeParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Propose, ser, opts...), nil
}
func (mp *MessageProducer) MultisigApprove(to, from address.Address, params multisig.TxnIDParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Approve, ser, opts...), nil
}
func (mp *MessageProducer) MultisigCancel(to, from address.Address, params multisig.TxnIDParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Cancel, ser, opts...), nil
}
func (mp *MessageProducer) MultisigAddSigner(to, from address.Address, params multisig.AddSignerParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.AddSigner, ser, opts...), nil
}
func (mp *MessageProducer) MultisigRemoveSigner(to, from address.Address, params multisig.RemoveSignerParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.RemoveSigner, ser, opts...), nil
}
func (mp *MessageProducer) MultisigSwapSigner(to, from address.Address, params multisig.SwapSignerParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.SwapSigner, ser, opts...), nil
}
func (mp *MessageProducer) MultisigChangeNumApprovalsThreshold(to, from address.Address, params multisig.ChangeNumApprovalsThresholdParams, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	return mp.Build(to, from, builtin_spec.MethodsMultisig.ChangeNumApprovalsThreshold, ser, opts...), nil
}

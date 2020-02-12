package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	multisig "github.com/filecoin-project/specs-actors/actors/builtin/multisig"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func (mp *MessageProducer) MultisigConstructor(to, from address.Address, params multisig.ConstructorParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Constructor, ser, opts...)
}
func (mp *MessageProducer) MultisigPropose(to, from address.Address, params multisig.ProposeParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Propose, ser, opts...)
}
func (mp *MessageProducer) MultisigApprove(to, from address.Address, params multisig.TxnIDParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Approve, ser, opts...)
}
func (mp *MessageProducer) MultisigCancel(to, from address.Address, params multisig.TxnIDParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMultisig.Cancel, ser, opts...)
}
func (mp *MessageProducer) MultisigAddSigner(to, from address.Address, params multisig.AddSignerParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMultisig.AddSigner, ser, opts...)
}
func (mp *MessageProducer) MultisigRemoveSigner(to, from address.Address, params multisig.RemoveSignerParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMultisig.RemoveSigner, ser, opts...)
}
func (mp *MessageProducer) MultisigSwapSigner(to, from address.Address, params multisig.SwapSignerParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMultisig.SwapSigner, ser, opts...)
}
func (mp *MessageProducer) MultisigChangeNumApprovalsThreshold(to, from address.Address, params multisig.ChangeNumApprovalsThresholdParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsMultisig.ChangeNumApprovalsThreshold, ser, opts...)
}

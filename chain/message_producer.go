package chain

import (
	address "github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/go-state-types/abi"
	big_spec "github.com/filecoin-project/go-state-types/big"

	"github.com/filecoin-project/chain-validation/chain/types"
)

// The created messages are retained for subsequent export or evaluation in a VM.
type MessageProducer struct {
	defaults msgOpts // Note non-pointer reference.

	messages []*types.Message
}

// NewMessageProducer creates a new message producer, delegating message creation to `factory`.
func NewMessageProducer(defaultGasFeeCap abi_spec.TokenAmount, defaultGasPremium abi_spec.TokenAmount, defaultGasLimit int64) *MessageProducer {
	return &MessageProducer{
		defaults: msgOpts{
			value:      big_spec.Zero(),
			gasLimit:   defaultGasLimit,
			gasFeeCap:  defaultGasFeeCap,
			gasPremium: defaultGasPremium,
		},
	}
}

// Messages returns a slice containing all messages created by the producer.
func (mp *MessageProducer) Messages() []*types.Message {
	return mp.messages
}

// BuildFull creates and returns a single message.
func (mp *MessageProducer) BuildFull(from, to address.Address, method abi_spec.MethodNum, callSeq uint64, value, gasFeeCap abi_spec.TokenAmount, gasPremium abi_spec.TokenAmount, gasLimit int64, params []byte) *types.Message {
	fm := &types.Message{
		To:         to,
		From:       from,
		CallSeqNum: callSeq,
		Value:      value,
		Method:     method,
		Params:     params,
		GasLimit:   gasLimit,
		GasFeeCap:  gasFeeCap,
		GasPremium: gasPremium,
	}
	mp.messages = append(mp.messages, fm)
	return fm
}

// Build creates and returns a single message, using default gas parameters unless modified by `opts`.
func (mp *MessageProducer) Build(from, to address.Address, method abi_spec.MethodNum, params []byte, opts ...MsgOpt) *types.Message {
	values := mp.defaults
	for _, opt := range opts {
		opt(&values)
	}

	return mp.BuildFull(from, to, method, values.nonce, values.value, values.gasFeeCap, values.gasPremium, values.gasLimit, params)
}

// msgOpts specifies value and gas parameters for a message, supporting a functional options pattern
// for concise but customizable message construction.
type msgOpts struct {
	nonce      uint64
	value      big_spec.Int
	gasLimit   int64
	gasFeeCap  abi_spec.TokenAmount
	gasPremium abi_spec.TokenAmount
}

// MsgOpt is an option configuring message value or gas parameters.
type MsgOpt func(*msgOpts)

func Value(value big_spec.Int) MsgOpt {
	return func(opts *msgOpts) {
		opts.value = value
	}
}

func Nonce(n uint64) MsgOpt {
	return func(opts *msgOpts) {
		opts.nonce = n
	}
}

func GasLimit(limit int64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasLimit = limit
	}
}

func GasFeeCap(feeCap int64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasFeeCap = abi_spec.NewTokenAmount(feeCap)
	}
}

func GasPremium(premium int64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasPremium = abi_spec.NewTokenAmount(premium)
	}
}

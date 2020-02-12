package chain

import (
	address "github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
)

type Message struct {
	// Address of the receiving actor.
	To address.Address
	// Address of the sending actor.
	From address.Address
	// Expected CallSeqNum of the sending actor (only for top-level messages).
	CallSeqNum int64

	// Amount of value to transfer from sender's to receiver's balance.
	Value big_spec.Int

	// Optional method to invoke on receiver, zero for a plain value send.
	Method abi_spec.MethodNum
	/// Serialized parameters to the method (if method is non-zero).
	Params []byte

	GasPrice big_spec.Int
	GasLimit big_spec.Int
}

type MessageFactory struct{}

func (b *MessageFactory) MakeMessage(to, from address.Address, method abi_spec.MethodNum, callSeq int64, value, gasPrice, gasLimit big_spec.Int, params []byte) *Message {
	return &Message{
		To:         to,
		From:       from,
		CallSeqNum: callSeq,
		Value:      value,
		Method:     method,
		Params:     params,
		GasPrice:   gasPrice,
		GasLimit:   gasLimit,
	}
}

// MessageProducer presents a convenient API for scripting the creation of long and complex message sequences.
// The created messages are retained for subsequent export or evaluation in a VM.
// Actual message construction is delegated to a `MessageFactory`.
type MessageProducer struct {
	factory  *MessageFactory
	defaults msgOpts // Note non-pointer reference.

	messages []*Message
}

// NewMessageProducer creates a new message producer, delegating message creation to `factory`.
func NewMessageProducer(defaultGasLimit, defaultGasPrice big_spec.Int) *MessageProducer {
	return &MessageProducer{
		factory: &MessageFactory{},
		defaults: msgOpts{
			gasLimit: defaultGasLimit,
			gasPrice: defaultGasPrice,
		},
	}
}

// Messages returns a slice containing all messages created by the producer.
func (mp *MessageProducer) Messages() []*Message {
	return mp.messages
}

// msgOpts specifies value and gas parameters for a message, supporting a functional options pattern
// for concise but customizable message construction.
type msgOpts struct {
	nonce    int64
	value    big_spec.Int
	gasLimit big_spec.Int
	gasPrice big_spec.Int
}

// MsgOpt is an option configuring message value or gas parameters.
type MsgOpt func(*msgOpts)

func Value(value big_spec.Int) MsgOpt {
	return func(opts *msgOpts) {
		opts.value = value
	}
}

func Nonce(n int64) MsgOpt {
	return func(opts *msgOpts) {
		opts.nonce = n
	}
}

func GasLimit(limit int64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasLimit = big_spec.NewInt(limit)
	}
}

func GasPrice(price int64) MsgOpt {
	return func(opts *msgOpts) {
		opts.gasPrice = big_spec.NewInt(price)
	}
}

// BuildFull creates and returns a single message.
func (mp *MessageProducer) BuildFull(to, from address.Address, method abi_spec.MethodNum, callSeq int64, value, gasLimit, gasPrice big_spec.Int, params []byte) *Message {
	fm := &Message{
		To:         to,
		From:       from,
		CallSeqNum: callSeq,
		Value:      value,
		Method:     method,
		Params:     params,
		GasPrice:   gasPrice,
		GasLimit:   gasLimit,
	}
	mp.messages = append(mp.messages, fm)
	return fm
}

// Build creates and returns a single message, using default gas parameters unless modified by `opts`.
func (mp *MessageProducer) Build(to, from address.Address, method abi_spec.MethodNum, params []byte, opts ...MsgOpt) *Message {
	values := mp.defaults
	for _, opt := range opts {
		opt(&values)
	}

	return mp.BuildFull(to, from, method, values.nonce, values.value, values.gasLimit, values.gasPrice, params)
}

//
// Sugar methods for type-checked construction of specific messages.
//

// Transfer builds a simple value transfer message and returns it.
func (mp *MessageProducer) Transfer(to, from address.Address, opts ...MsgOpt) *Message {
	return mp.Build(to, from, builtin_spec.MethodSend, noParams, opts...)
}

var noParams []byte

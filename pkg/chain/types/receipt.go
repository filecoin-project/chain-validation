package types

import (
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
)

// MessageReceipt is the return value of message application.
type MessageReceipt struct {
	ExitCode    exitcode.ExitCode
	ReturnValue []byte
	GasUsed     big.Int
}

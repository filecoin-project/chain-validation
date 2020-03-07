package drivers

import (
	"fmt"
	"os"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/fixtures"
)

type GasRecorder struct {
	file *os.File

	recordF func(oldState, newState cid.Cid, msg *types.Message, gas int64)
}

func NewGasRecorder(file string, record bool) (*GasRecorder, error) {
	gv := &GasRecorder{}
	gv.recordF = func(oldState, newState cid.Cid, msg *types.Message, gas int64) {}
	if record {
		f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		gv.recordF = func(oldState, newState cid.Cid, msg *types.Message, gas int64) {
			if _, err := fmt.Fprintf(f, "\"%s-%s-%s\":%d,\n", oldState, msg.Cid(), newState, gas); err != nil {
				panic(err)
			}
		}
		gv.file = f
	}

	return gv, nil
}

func (g *GasRecorder) Record(oldState, newState cid.Cid, msg *types.Message, gas int64) {
	g.recordF(oldState, newState, msg, gas)
}

func (g *GasRecorder) GasFor(oldState, newState cid.Cid, msg *types.Message) int64 {
	stateMsgStr := fmt.Sprintf("%s-%s-%s", oldState, msg.Cid(), newState)
	gas, ok := fixtures.StateToGasMap[stateMsgStr]
	if !ok {
		panic(fmt.Sprintf("Unknown message and states: %v", msg))
	}
	return gas
}

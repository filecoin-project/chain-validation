package drivers

import (
	"container/list"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/gas_gen/gas"
)

type trackerElement struct {
	oldState cid.Cid
	newState cid.Cid
	msg      *types.Message
	receipt  types.MessageReceipt
}

func (te *trackerElement) fileKey() string {
	return fmt.Sprintf("%s,%s,%s,%d", te.oldState, te.msg.Cid(), te.newState, te.receipt.GasUsed.Int64())
}

type GasMeter struct {
	tracker *list.List
	T       testing.TB
}

func NewGasMeter(t testing.TB) *GasMeter {
	return &GasMeter{
		tracker: list.New(),
		T:       t,
	}
}

func (gm *GasMeter) Track(oldState, newState cid.Cid, msg *types.Message, receipt types.MessageReceipt) {
	gm.tracker.PushBack(&trackerElement{
		oldState: oldState,
		newState: newState,
		msg:      msg,
		receipt:  receipt,
	})
}

func (gm *GasMeter) Expected(oldState cid.Cid, msg *types.Message) (int64, bool) {
	key := makeKey(oldState, msg)
	gasUsed, ok := gas.GasConstants[key]
	return gasUsed, ok
}

func (gm *GasMeter) Record() {
	file := fileNameFromTest(gm.T)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		gm.T.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	for e := gm.tracker.Front(); e != nil; e = e.Next() {
		_, err := fmt.Fprintf(f, "%s\n", e.Value.(*trackerElement).fileKey())
		if err != nil {
			gm.T.Fatal(err)
		}
	}
}

func fileNameFromTest(t testing.TB) string {
	// need to remove all '/' from the file name
	fileName := strings.ReplaceAll(t.Name(), "/", "_")
	prefix := "GASFILE_"
	return fmt.Sprintf("%s%s", prefix, fileName)
}

func makeKey(oldState cid.Cid, msg *types.Message) string {
	return fmt.Sprintf("%s-%s", oldState, msg.Cid())
}
package gasmeter

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/chain-validation/box"
	"github.com/filecoin-project/chain-validation/chain/types"
)

const ValidationDataEnvVar = "CHAIN_VALIDATION_DATA"

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
	// index in gasUnits of expected gas
	gasIdx int
	// slice of gas units used by the test
	expectedGasUnits []int64
}

func NewGasMeter(t testing.TB) *GasMeter {
	var gasUnits []int64
	return &GasMeter{
		tracker:          list.New(),
		T:                t,
		gasIdx:           0,
		expectedGasUnits: gasUnits,
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

func (gm *GasMeter) ExpectedGasUnit() (int64, bool) {
	defer func() { gm.gasIdx += 1 }()
	if gm.gasIdx > len(gm.expectedGasUnits)-1 {
		// didn't find any gas
		return 0, false
	}
	return gm.expectedGasUnits[gm.gasIdx], true
}

// write the contents of gm.tracker to a file using the format:
// oldStateCID,msgCID,newStateCID,GasUnits
func (gm *GasMeter) Record() {
	file := getTestDataFilePath(gm.T)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		gm.T.Fatal(err)
	}
	defer f.Close()

	for e := gm.tracker.Front(); e != nil; e = e.Next() {
		_, err := fmt.Fprintf(f, "%s\n", e.Value.(*trackerElement).fileKey())
		if err != nil {
			gm.T.Fatal(err)
		}
	}
}

// Given a testing T, load the gas file associated with it and return a slice of the gas used by the test
// an index in the slice represents the order of apply message calls.
func LoadGasForTest(t testing.TB) []int64 {
	f, found := box.Get(filenameFromTest(t))
	if !found {
		t.Fatal("can't find file")
	}

	var gasUnits []int64
	scanner := bufio.NewScanner(bytes.NewReader(f))
	for scanner.Scan() {
		gas, err := gasFromTestFileLine(scanner.Text())
		if err != nil {
			t.Fatal(err)
		}
		gasUnits = append(gasUnits, gas)
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	return gasUnits
}

func getTestDataFilePath(t testing.TB) string {
	dataPath := os.Getenv(ValidationDataEnvVar)
	if dataPath == "" {
		// FIXME before mergin, make dev easier
		//t.Fatalf("failed to find validation data path, make sure %s is set", ValidationDataEnvVar)
		dataPath = "/home/frrist/src/github.com/filecoin-project/chain-validation/gasmeter/gas_files"
	}
	return filepath.Join(dataPath, filenameFromTest(t))
}

// given a line of the form:
// oldState,messageCid,newState,gasUnits
// return gasUnits as an int64
func gasFromTestFileLine(l string) (int64, error) {
	tokens := strings.Split(l, ",")
	// expect tokens to always be length 4 (oldState,msgCid,newState,gasCost)
	if len(tokens) != 4 {
		return -1, fmt.Errorf("invalid gas line, expected 4 tokens, got %d: %s", len(tokens), tokens)
	}

	//oldState := tokens[0]
	//msgCid := tokens[1]
	//newState := tokens[2]
	gasUnits := tokens[3]

	g, err := strconv.ParseInt(gasUnits, 10, 64)
	if err != nil {
		return -1, err
	}
	return g, nil
}

// return a string containing only letters and number.
func filenameFromTest(t testing.TB) string {
	// only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("/%s", reg.ReplaceAllString(t.Name(), ""))
}

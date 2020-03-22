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

	"github.com/filecoin-project/chain-validation/box"
	"github.com/filecoin-project/chain-validation/chain/types"
)

const ValidationDataEnvVar = "CHAIN_VALIDATION_DATA"

type trackerElement struct {
	receipt types.MessageReceipt
}

func (te *trackerElement) fileKey() string {
	return fmt.Sprintf("%d", te.receipt.GasUsed.Int64())
}

type GasMeter struct {
	tracker *list.List
	T       testing.TB

	record bool
	// index in gasUnits of expected gas
	gasIdx int
	// slice of gas units used by the test
	expectedGasUnits []int64
}

func NewGasMeter(t testing.TB, record bool) *GasMeter {
	return &GasMeter{
		tracker:          list.New(),
		T:                t,
		record:           record,
		gasIdx:           0,
		expectedGasUnits: LoadGasForTest(t),
	}
}

func (gm *GasMeter) Track(receipt types.MessageReceipt) {
	gm.tracker.PushBack(&trackerElement{
		receipt: receipt,
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
	if !gm.record {
		return
	}
	file := getTestDataFilePath(gm.T)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		gm.T.Log(err)
		return
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
	fileName := filenameFromTest(t)
	f, found := box.Get(fileName)
	if !found {
		t.Logf("can't find file: %s", fileName)
		// return an empty slice here since `ExpectedGasUnit` performs bounds checking
		return []int64{}
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
		t.Fatalf("failed to find validation data path, make sure %s is set", ValidationDataEnvVar)
	}
	return filepath.Join(dataPath, filenameFromTest(t))
}

// given a line of the form:
// oldState,messageCid,newState,gasUnits
// return gasUnits as an int64
func gasFromTestFileLine(l string) (int64, error) {
	tokens := strings.Split(l, ",")
	// expect tokens to always be length 3 (oldState,newState,gasCost)
	if len(tokens) != 1 {
		return -1, fmt.Errorf("invalid gas line, expected 1 tokens, got %d: %s", len(tokens), tokens)
	}

	gasUnits := tokens[0]

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

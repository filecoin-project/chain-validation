package gasmeter

import (
	"container/list"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/chain-validation/boxs/gas"
	"github.com/filecoin-project/chain-validation/boxs/stateroot"
	"github.com/filecoin-project/chain-validation/chain/types"
)

const ValidationDataEnvVar = "CHAIN_VALIDATION_DATA"

type receiptElement struct {
	receipt types.MessageReceipt
}

func (re *receiptElement) fileKey() string {
	return fmt.Sprintf("%d", re.receipt.GasUsed)
}

type stateRootElement struct {
	state cid.Cid
}

func (se *stateRootElement) fileKey() string {
	return se.state.String()
}

type GasMeter struct {
	T testing.TB

	receipts   *list.List
	stateroots *list.List

	// index in gasUnits of expected gas
	gasIdx int
	// slice of gas units used by the test
	expectedGasUnits []int64

	// index in stateRoots of expected state root
	staterootIdx int
	// slice of state stateroots used by test
	expectedStateRoots []cid.Cid
}

func NewGasMeter(t testing.TB) *GasMeter {
	return &GasMeter{
		T: t,

		receipts:   list.New(),
		stateroots: list.New(),

		gasIdx:           0,
		expectedGasUnits: LoadGasForTest(t),

		staterootIdx:       0,
		expectedStateRoots: LoadStateRootsForTest(t),
	}
}

func (gm *GasMeter) TrackReceipt(receipt types.MessageReceipt) {
	gm.receipts.PushBack(&receiptElement{receipt: receipt})
}

func (gm *GasMeter) TrackStateRoot(root cid.Cid) {
	gm.stateroots.PushBack(&stateRootElement{state: root})
}

func (gm *GasMeter) NextExpectedGas() (types.GasUnits, bool) {
	defer func() { gm.gasIdx += 1 }()
	if gm.gasIdx > len(gm.expectedGasUnits)-1 {
		// didn't find any gas
		return 0, false
	}
	return types.GasUnits(gm.expectedGasUnits[gm.gasIdx]), true
}

func (gm *GasMeter) NextExpectedStateRoot() (cid.Cid, bool) {
	defer func() { gm.staterootIdx += 1 }()
	if gm.staterootIdx > len(gm.expectedStateRoots)-1 {
		return cid.Undef, false
	}
	return gm.expectedStateRoots[gm.staterootIdx], true
}

func (gm *GasMeter) Record() {
	gm.recordGas()
	gm.recordStateRoots()
}

// write the contents of gm.receipts to a file using the format:
// GasUnit
// GasUnit
// ...
func (gm *GasMeter) recordGas() {
	file := getTestDataFilePath(gm.T)
	f, err := os.OpenFile(fmt.Sprintf("%s_GAS", file), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		gm.T.Log(err)
		return
	}
	defer f.Close()

	for e := gm.receipts.Front(); e != nil; e = e.Next() {
		_, err := fmt.Fprintf(f, "%s\n", e.Value.(*receiptElement).fileKey())
		if err != nil {
			gm.T.Fatal(err)
		}
	}
}

// write the contents of gm.stateroots to a file using the format:
// StateRoot
// StateRoot
// ...
func (gm *GasMeter) recordStateRoots() {
	file := getTestDataFilePath(gm.T)
	f, err := os.OpenFile(fmt.Sprintf("%s_STATEROOT", file), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		gm.T.Log(err)
		return
	}
	defer f.Close()

	for e := gm.stateroots.Front(); e != nil; e = e.Next() {
		_, err := fmt.Fprintf(f, "%s\n", e.Value.(*stateRootElement).fileKey())
		if err != nil {
			gm.T.Fatal(err)
		}
	}
}

// Given a testing T, load the gas file associated with it and return a slice of the gas used by the test
// an index in the slice represents the order of apply message calls.
func LoadGasForTest(t testing.TB) []int64 {
	fileName := filenameFromTest(t)
	f, found := gas.Get(fmt.Sprintf("%s_GAS", fileName))
	if !found {
		t.Logf("WARNING (does NOT indicate test failure): can't find gas file: %s", fileName)
		// return an empty slice here since `NextExpectedGas` performs bounds checking
		return []int64{}
	}
	return f
}

func LoadStateRootsForTest(t testing.TB) []cid.Cid {
	fileName := filenameFromTest(t)
	f, found := stateroot.Get(fmt.Sprintf("%s_STATEROOT", fileName))
	if !found {
		t.Logf("WARNING (does NOT indicate test failure): can't find stateroot file: %s", fileName)
		// return an empty slice here since `NextExpectedStateRoot` performs bounds checking
		return []cid.Cid{}
	}
	var out []cid.Cid
	for _, c := range f {
		root, err := cid.Decode(c)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, root)
	}
	return out
}

func getTestDataFilePath(t testing.TB) string {
	dataPath := os.Getenv(ValidationDataEnvVar)
	if dataPath == "" {
		t.Fatalf("failed to find validation data path, make sure %s is set", ValidationDataEnvVar)
	}
	return filepath.Join(dataPath, filenameFromTest(t))
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

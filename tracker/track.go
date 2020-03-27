package tracker

import (
	"container/list"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/chain-validation/box"
	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
)

const ValidationDataEnvVar = "CHAIN_VALIDATION_DATA"

type StateTracker struct {
	tracker *list.List
	T       testing.TB

	// index in gasUnits of expected gas
	gasIdx int
	// slice of gas units used by the test
	expectedGasUnits []types.GasUnits

	// index in stateroots of expected cid
	rootIdx int
	// slice of state roots used by the test
	expectedStateRoots []cid.Cid
}

func NewStateTracker(t testing.TB) *StateTracker {
	gasUsed, stateRoots := LoadDataForTest(t)
	return &StateTracker{
		tracker:            list.New(),
		T:                  t,
		gasIdx:             0,
		expectedGasUnits:   gasUsed,
		rootIdx:            0,
		expectedStateRoots: stateRoots,
	}
}

func (st *StateTracker) TrackMessageResult(result chain.ApplyMessageResult) {
	st.tracker.PushBack(result)
}

func (st *StateTracker) TrackTipSetMessagesResult(result chain.ApplyTipSetMessagesResult) {
	st.tracker.PushBack(result)
}

func (st *StateTracker) NextExpectedGas() (types.GasUnits, bool) {
	defer func() { st.gasIdx += 1 }()
	if st.gasIdx > len(st.expectedGasUnits)-1 {
		// didn't find any gas
		return 0, false
	}
	return st.expectedGasUnits[st.gasIdx], true
}

func (st *StateTracker) NextExpectedStateRoot() (cid.Cid, bool) {
	defer func() { st.rootIdx += 1 }()
	if st.rootIdx > len(st.expectedStateRoots)-1 {
		// didn't find any gas
		return cid.Undef, false
	}
	return st.expectedStateRoots[st.rootIdx], true
}

// write the contents of gm.tracker to a file using the format:
// GasUnit
// GasUnit
// ...
func (st *StateTracker) Record() {
	file := getTestDataFilePath(st.T)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		st.T.Log(err)
		return
	}
	defer func() { _ = f.Close() }()
	enc := json.NewEncoder(f)

	for e := st.tracker.Front(); e != nil; e = e.Next() {
		switch ele := e.Value.(type) {
		case chain.ApplyMessageResult:
			if err := enc.Encode(ele); err != nil {
				st.T.Fatal(err)
			}
		case chain.ApplyTipSetMessagesResult:
			if err := enc.Encode(ele); err != nil {
				st.T.Fatal(err)
			}
		default:
			st.T.Fatalf("Unknown type: %T", ele)
		}
	}
}

func LoadDataForTest(t testing.TB) (gasUsed []types.GasUnits, stateRoots []cid.Cid) {
	fileName := filenameFromTest(t)
	data, found := box.Get(fileName)
	if !found {
		t.Logf("WARNING (does NOT indicate test failure): can't find file: %s", fileName)
		// return an empty slice here since `NextExpectedGas` performs bounds checking
		return []types.GasUnits{}, []cid.Cid{}
	}
	switch v := data.(type) {
	case chain.ApplyMessageResult:
		root, err := cid.Decode(v.Root)
		if err != nil {
			t.Fatal(err)
		}
		gasUsed = append(gasUsed, v.Receipt.GasUsed)
		stateRoots = append(stateRoots, root)
		return
	case []chain.ApplyMessageResult:
		for _, res := range v {
			gasUsed = append(gasUsed, res.Receipt.GasUsed)
			root, err := cid.Decode(res.Root)
			if err != nil {
				t.Fatal(err)
			}
			stateRoots = append(stateRoots, root)
		}
		return
	case chain.ApplyTipSetMessagesResult:
		root, err := cid.Decode(v.Root)
		if err != nil {
			t.Fatal(err)
		}
		for _, rect := range v.Receipts {
			gasUsed = append(gasUsed, rect.GasUsed)
		}
		stateRoots = append(stateRoots, root)
		return
	case []chain.ApplyTipSetMessagesResult:
		for _, res := range v {
			for _, rect := range res.Receipts {
				gasUsed = append(gasUsed, rect.GasUsed)
			}
			root, err := cid.Decode(res.Root)
			if err != nil {
				t.Fatal(err)
			}
			stateRoots = append(stateRoots, root)
		}
		return
	default:
		t.Fatalf("Unknown Test Data Type: %T", v)
	}
	panic("unreachable")
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

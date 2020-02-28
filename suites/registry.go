package suites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/message"
)

type ChainValidationTestCase struct {
	CreateSecpAccountActor ChainValidationTest
	CreateBlsAccountActor  ChainValidationTest

	FailCreateSecpAccountActor ChainValidationTest
	FailCreateBlsAccountActor  ChainValidationTest

	InitActorSequentialID ChainValidationTest
}

type ChainValidationTest func(v *drivers.ValidationHarness, factory state.Factories)

func ValueTransferTestCases() ChainValidationTestCase {
	return ChainValidationTestCase{
		CreateSecpAccountActor: message.SuccessfullyCreateSECP256K1AccountActor,
		CreateBlsAccountActor:  message.SuccessfullyCreateBLSAccountActor,

		FailCreateSecpAccountActor: message.FailToCreateSECP256K1AccountActorWithInsufficientBalance,
		FailCreateBlsAccountActor:  message.FailToCreateBLSAccountActorWithInsufficientBalance,

		InitActorSequentialID: message.InitActorCreatesActorsWithSequentialIDAddresses,
	}
}

func TestStuff() {

}

type TestResult bool

var Pass = TestResult(true)
var Fail = TestResult(false)

func AssertTestResult(t *testing.T, testCase ChainValidationTest, result TestResult, f state.Factories) {
	harness := drivers.NewValidationHarness(t)

	testCase(harness, f)

	if result == Pass {
		require.False(t, harness.Failed(), "expected %v pass but failed")
	}

	if result == Fail {
		require.True(t, harness.Failed(), "expected %v fail but passes")
	}
}

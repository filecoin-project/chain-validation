package suites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/chain-validation/state"

	"github.com/filecoin-project/chain-validation/suites/message/actors/account"
)

type ExpectFailure struct {
	hasFailed bool
}

func (ef *ExpectFailure) Private() {
	panic("implement me")
}

func (ef *ExpectFailure) Fatal(args ...interface{}) {
	ef.hasFailed = true
}

func (ef *ExpectFailure) Fatalf(format string, args ...interface{}) {
	ef.hasFailed = true
}

func (ef *ExpectFailure) Log(args ...interface{}) {
	panic("implement me")
}

func (ef *ExpectFailure) Logf(format string, args ...interface{}) {
	panic("implement me")
}

func (ef *ExpectFailure) Name() string {
	panic("implement me")
}

func (ef *ExpectFailure) Skip(args ...interface{}) {
	panic("implement me")
}

func (ef *ExpectFailure) SkipNow() {
	panic("implement me")
}

func (ef *ExpectFailure) Skipf(format string, args ...interface{}) {
	panic("implement me")
}

func (ef *ExpectFailure) Skipped() bool {
	panic("implement me")
}

func (ef *ExpectFailure) Helper() {
	panic("implement me")
}

func (ef *ExpectFailure) Fail() {
	ef.hasFailed = true
}

func (ef *ExpectFailure) FailNow() {
	ef.hasFailed = true
}
func (ef *ExpectFailure) Error(args ...interface{}) {

	ef.hasFailed = true
}

func (ef *ExpectFailure) Errorf(format string, args ...interface{}) {
	ef.hasFailed = true
}

func (ef *ExpectFailure) Failed() bool {
	return ef.hasFailed
}

type ChainValidationTestCase struct {
	CreateSecpAccountActor     ChainValidationTest
	FailCreateSecpAccountActor ChainValidationTest

	CreateBlsAccountActor     ChainValidationTest
	FailCreateBlsAccountActor ChainValidationTest
}

type ChainValidationTest func(t testing.TB, factory state.Factories)

func ValueTransferTestCases() ChainValidationTestCase {
	return ChainValidationTestCase{
		account.SuccessfullyCreateSECP256K1AccountActor,
		account.FailToCreateSECP256K1AccountActorWithInsufficientBalance,
		account.SuccessfullyCreateBLSAccountActor,
		account.FailToCreateBLSAccountActorWithInsufficientBalance,
	}
}

type TestResult bool

var Pass = TestResult(true)
var Fail = TestResult(false)

func AssertTestResult(t *testing.T, test ChainValidationTest, result TestResult, factory state.Factories, config state.ValidationConfig) {
	ef := &ExpectFailure{
		hasFailed: false,
	}
	test(ef, factory)

	if result == Pass {
		require.Equal(t, false, ef.hasFailed)
	}
	if result == Fail {
		require.Equal(t, true, ef.hasFailed)
	}
}

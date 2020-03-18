package suites

import (
	"testing"

	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/message"
	"github.com/filecoin-project/chain-validation/suites/tipset"
)

type TestCase func(t *testing.T, factory state.Factories)

func MessageTestCases() []TestCase {
	return []TestCase{
		message.TestAccountActorCreation,
		message.TestInitActorSequentialIDAddressCreate,
		message.TestMessageApplicationEdgecases,
		message.TestMultiSigActor,
		message.TestPaych,
		message.TestValueTransferAdvance,
		message.TestValueTransferSimple,
	}
}

func TipSetTestCases() []TestCase {
	return []TestCase{
		tipset.TestBlockMessageDeduplication,
		tipset.TestInternalMessageApplicationFailure,
		tipset.TestInvalidSenderAddress,
		tipset.TestMinerMissPoStChallengeWindow,
		tipset.TestMinerSubmitFallbackPoSt,
	}
}

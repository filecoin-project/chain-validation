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
		message.MessageTest_AccountActorCreation,
		message.MessageTest_InitActorSequentialIDAddressCreate,
		message.MessageTest_MessageApplicationEdgecases,
		message.MessageTest_MultiSigActor,
		message.MessageTest_NestedSends,
		message.MessageTest_Paych,
		message.MessageTest_ValueTransferAdvance,
		message.MessageTest_ValueTransferSimple,
	}
}

func TipSetTestCases() []TestCase {
	return []TestCase{
		tipset.TipSetTest_BlockMessageApplication,
		tipset.TipSetTest_BlockMessageDeduplication,
		tipset.TipSetTest_MinerRewardsAndPenalties,
	}
}

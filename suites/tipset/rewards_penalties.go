package tipset

import (
	"context"
	"testing"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/reward"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

// Test for semantically in/valid messages, including miner penalties.
func TestMinerRewardsAndPenalties(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	acctDefaultBalance := abi.NewTokenAmount(1_000_000_000)
	sendValue := abi.NewTokenAmount(1)

	t.Run("ok simple send", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		tipB := drivers.NewTipSetMessageBuilder(td)
		miner := td.ExeCtx.Miner

		alicePk, aliceId := td.NewAccountActor(drivers.SECP, acctDefaultBalance)
		bobPk, bobId := td.NewAccountActor(drivers.SECP, acctDefaultBalance)

		// Exercise all combinations of ID and PK address for the sender.
		callSeq := uint64(0)
		for _, alice := range []addr.Address{alicePk, aliceId} {
			for _, bob := range []addr.Address{bobPk, bobId} {
				aBal := td.GetBalance(aliceId)
				bBal := td.GetBalance(bobId)
				prevRewards := td.GetRewardSummary()

				// Process a block with two messages, a simple send back and forth between accounts.
				rcpts := tipB.WithBlockBuilder(
					drivers.NewBlockBuilder(td.ExeCtx.Miner).
						WithTicketCount(1).
						WithBLSMessageOk(
							td.MessageProducer.Transfer(bob, alice, chain.Value(sendValue), chain.Nonce(callSeq)),
						).
						WithBLSMessageOk(
							td.MessageProducer.Transfer(alice, bob, chain.Value(sendValue), chain.Nonce(callSeq)),
						),
				).Apply()
				assert.Equal(t, exitcode.Ok, rcpts[0].ExitCode)
				assert.Equal(t, exitcode.Ok, rcpts[1].ExitCode)
				tipB.Clear()

				// Each account has paid gas fees.
				td.AssertBalance(aliceId, big.Sub(aBal, rcpts[0].GasUsed.AsBigInt()))
				td.AssertBalance(bobId, big.Sub(bBal, rcpts[1].GasUsed.AsBigInt()))
				gasSum := big.Add(rcpts[0].GasUsed.AsBigInt(), rcpts[1].GasUsed.AsBigInt()) // Exploit gas price = 1

				// Validate rewards.
				// No reward is paid to the miner directly. The funds for block reward were already held by the
				// reward actor. The gas reward should be added to the treasury. The sum of block and gas reward
				// should be chalked up to the miner address.
				thisReward := big.Add(reward.BlockRewardTarget, gasSum)
				newRewards := td.GetRewardSummary()
				assert.Equal(t, big.Add(prevRewards.Treasury, gasSum), newRewards.Treasury)
				assert.Equal(t, big.Add(prevRewards.For(miner), thisReward), newRewards.For(miner))
				assert.Equal(t, big.Add(prevRewards.RewardTotal, thisReward), newRewards.RewardTotal)
				assert.Equal(t, big.Zero(), td.GetBalance(builtin.BurntFundsActorAddr))

				callSeq++
			}
		}
	})

	t.Run("penalize sender does't exist", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		blkBuilder := drivers.NewBlockBuilder(td.ExeCtx.Miner)
		miner := td.ExeCtx.Miner

		_, receiver := td.NewAccountActor(drivers.SECP, acctDefaultBalance)
		badSenders := []addr.Address{
			utils.NewIDAddr(t, 1234),
			utils.NewSECP256K1Addr(t, "1234"),
			utils.NewBLSAddr(t, 1234),
			utils.NewActorAddr(t, "1234"),
		}

		bb := blkBuilder.WithTicketCount(1)
		for _, s := range badSenders {
			bb.WithBLSMessageAndCode(td.MessageProducer.Transfer(receiver, s, chain.Value(sendValue)),
				exitcode.SysErrActorNotFound,
			)
		}

		prevRewards := td.GetRewardSummary()
		drivers.NewTipSetMessageBuilder(td).WithBlockBuilder(bb).ApplyAndValidate()

		// Nothing received, no actors created.
		td.AssertBalance(receiver, acctDefaultBalance)
		for _, s := range badSenders {
			td.AssertNoActor(s)
		}

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(342)

		// The penalty amount has been burnt by the reward actor, and subtracted from the miner's block reward
		rwd := big.Sub(reward.BlockRewardTarget, gasPenalty)
		assert.Equal(t, big.Sub(prevRewards.Treasury, gasPenalty), newRewards.Treasury)
		assert.Equal(t, gasPenalty, td.GetBalance(builtin.BurntFundsActorAddr))
		assert.Equal(t, rwd, newRewards.For(miner))
		assert.Equal(t, rwd, newRewards.RewardTotal)
	})

	// TODO more tests:
	// - sender exists but isn't an account (miner penalty)
	// - mismatched callseqnum (miner penalty)
	// - sender cannot cover value + gas cost (miner penalty)
	// - miner penalty followed by non-miner penalty with same nonce
}

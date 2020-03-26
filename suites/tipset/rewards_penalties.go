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
				result := tipB.WithBlockBuilder(
					drivers.NewBlockBuilder(td, td.ExeCtx.Miner).
						WithBLSMessageOk(
							td.MessageProducer.Transfer(bob, alice, chain.Value(sendValue), chain.Nonce(callSeq)),
						).
						WithBLSMessageOk(
							td.MessageProducer.Transfer(alice, bob, chain.Value(sendValue), chain.Nonce(callSeq)),
						),
				).ApplyAndValidate()
				tipB.Clear()

				// Each account has paid gas fees.
				td.AssertBalance(aliceId, big.Sub(aBal, result.Receipts[0].GasUsed.Big()))
				td.AssertBalance(bobId, big.Sub(bBal, result.Receipts[1].GasUsed.Big()))
				gasSum := big.Add(result.Receipts[0].GasUsed.Big(), result.Receipts[1].GasUsed.Big()) // Exploit gas price = 1

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
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)
		miner := td.ExeCtx.Miner

		_, receiver := td.NewAccountActor(drivers.SECP, acctDefaultBalance)
		badSenders := []addr.Address{
			utils.NewIDAddr(t, 1234),
			utils.NewSECP256K1Addr(t, "1234"),
			utils.NewBLSAddr(t, 1234),
			utils.NewActorAddr(t, "1234"),
		}

		for _, s := range badSenders {
			bb.WithBLSMessageAndCode(td.MessageProducer.Transfer(receiver, s, chain.Value(sendValue)),
				exitcode.SysErrSenderInvalid,
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
		validateRewards(t, prevRewards, newRewards, miner, big.Zero(), gasPenalty)
		td.AssertBalance(builtin.BurntFundsActorAddr, gasPenalty)
	})

	t.Run("penalize sender non account", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, miner)

		_, receiver := td.NewAccountActor(drivers.SECP, acctDefaultBalance)
		// Various non-account actors that can't be top-level senders.
		senders := []addr.Address{
			builtin.SystemActorAddr,
			builtin.InitActorAddr,
			builtin.CronActorAddr,
			miner,
		}

		for _, sender := range senders {
			bb.WithBLSMessageAndCode(td.MessageProducer.Transfer(receiver, sender, chain.Value(sendValue)),
				exitcode.SysErrSenderInvalid)
		}
		prevRewards := td.GetRewardSummary()
		tb.WithBlockBuilder(bb).ApplyAndValidate()
		td.AssertBalance(receiver, acctDefaultBalance)

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(168)

		// The penalty amount has been burnt by the reward actor, and subtracted from the miner's block reward.
		validateRewards(t, prevRewards, newRewards, miner, big.Zero(), gasPenalty)
		td.AssertBalance(builtin.BurntFundsActorAddr, gasPenalty)
	})

	t.Run("penalize wrong callseqnum", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		_, aliceId := td.NewAccountActor(drivers.BLS, acctDefaultBalance)

		bb.WithBLSMessageAndCode(
			td.MessageProducer.Transfer(builtin.BurntFundsActorAddr, aliceId, chain.Nonce(1)),
			exitcode.SysErrSenderStateInvalid,
		)

		prevRewards := td.GetRewardSummary()
		tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(38)
		validateRewards(t, prevRewards, newRewards, miner, big.Zero(), gasPenalty)
		td.AssertBalance(builtin.BurntFundsActorAddr, gasPenalty)
	})

	t.Run("penalize sender insufficient balance", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		halfBalance := abi.NewTokenAmount(10_000_000)
		_, aliceId := td.NewAccountActor(drivers.BLS, big.Add(halfBalance, halfBalance))

		// Attempt to whole balance, in two parts.
		// The second message should fail (insufficient balance to pay fees).
		bb.WithBLSMessageOk(
			td.MessageProducer.Transfer(builtin.BurntFundsActorAddr, aliceId, chain.Value(halfBalance)),
		).WithBLSMessageAndCode(
			td.MessageProducer.Transfer(builtin.BurntFundsActorAddr, aliceId, chain.Value(halfBalance), chain.Nonce(1)),
			exitcode.SysErrSenderStateInvalid,
		)

		prevRewards := td.GetRewardSummary()
		result := tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(46)
		validateRewards(t, prevRewards, newRewards, miner, result.Receipts[0].GasUsed.Big(), gasPenalty)
		td.AssertBalance(builtin.BurntFundsActorAddr, big.Add(halfBalance, gasPenalty))
	})

	// TODO more tests:
	// - miner penalty causes subsequent otherwise-valid message to have wrong nonce (another miner penalty)
	// - miner penalty followed by non-miner penalty with same nonce (in different block)
}

func validateRewards(t testing.TB, prevRewards *drivers.RewardSummary, newRewards *drivers.RewardSummary, miner addr.Address, gasReward big.Int, gasPenalty big.Int) {
	rwd := big.Add(big.Sub(reward.BlockRewardTarget, gasPenalty), gasReward)
	assert.Equal(t, big.Add(big.Sub(prevRewards.Treasury, gasPenalty), gasReward), newRewards.Treasury)
	assert.Equal(t, big.Add(prevRewards.For(miner), rwd), newRewards.For(miner))
	assert.Equal(t, big.Add(prevRewards.RewardTotal, rwd), newRewards.RewardTotal)
}

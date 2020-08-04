package tipset

import (
	"context"
	"testing"

	miner_spec "github.com/filecoin-project/specs-actors/actors/builtin/miner"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

// Test for semantically in/valid messages, including miner penalties.
func TipSetTest_MinerRewardsAndPenalties(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000_000).
		WithDefaultGasPrice(big.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState...)

	acctDefaultBalance := abi.NewTokenAmount(10_000_000_000)
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
				burnBal := td.GetBalance(builtin.BurntFundsActorAddr)
				prevRewards := td.GetRewardSummary()
				prevMinerBal := td.GetBalance(miner)

				// Process a block with two messages, a simple send back and forth between accounts.
				msg1 := td.MessageProducer.Transfer(alice, bob, chain.Value(sendValue), chain.Nonce(callSeq))
				msg2 := td.MessageProducer.Transfer(bob, alice, chain.Value(sendValue), chain.Nonce(callSeq))
				result := tipB.WithBlockBuilder(
					drivers.NewBlockBuilder(td, td.ExeCtx.Miner).
						WithBLSMessageOk(msg1).
						WithBLSMessageOk(msg2),
				).ApplyAndValidate()
				tipB.Clear()

				td.ExeCtx.Epoch++

				// Each account has paid gas fees.
				td.AssertBalance(aliceId, big.Sub(aBal, td.CalcMessageCost(msg1.GasLimit, msg1.GasPrice, big.Zero(), result.Receipts[0])))
				td.AssertBalance(bobId, big.Sub(bBal, td.CalcMessageCost(msg2.GasLimit, msg2.GasPrice, big.Zero(), result.Receipts[1])))

				gasSum := big.Add(result.Receipts[0].GasUsed.Big(), result.Receipts[1].GasUsed.Big()) // Exploit gas price = 1

				// Validate rewards are paid directly to miner
				newRewards := td.GetRewardSummary()

				// total supply should decrease by the last reward amount
				assert.Equal(t, big.Sub(prevRewards.Treasury, prevRewards.NextPerBlockReward), newRewards.Treasury)

				// the miners balance should have increased by the reward amount
				thisReward := big.Add(prevRewards.NextPerBlockReward, gasSum)
				assert.Equal(t, big.Add(prevMinerBal, thisReward), td.GetBalance(miner))

				newBurn := big.Add(drivers.GetBurn(msg1.GasLimit, result.Receipts[0].GasUsed), drivers.GetBurn(msg2.GasLimit, result.Receipts[1].GasUsed))
				td.AssertBalance(builtin.BurntFundsActorAddr, big.Add(burnBal, newBurn))

				callSeq++
			}
		}
	})

	t.Run("penalize sender doesn't exist", func(t *testing.T) {
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
			bb.WithBLSMessageAndCode(td.MessageProducer.Transfer(s, receiver, chain.Value(sendValue)),
				exitcode.SysErrSenderInvalid,
			)
		}

		prevRewards := td.GetRewardSummary()
		prevMinerBalance := td.GetBalance(miner)
		drivers.NewTipSetMessageBuilder(td).WithBlockBuilder(bb).ApplyAndValidate()

		// Nothing received, no actors created.
		td.AssertBalance(receiver, acctDefaultBalance)
		for _, s := range badSenders {
			td.AssertNoActor(s)
		}

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		newMinerBalance := td.GetBalance(miner)
		gasPenalty := big.NewInt(867548)

		// The penalty amount has been burnt by the reward actor, and subtracted from the miner's block reward
		validateRewards(td, prevRewards, newRewards, prevMinerBalance, newMinerBalance, big.Zero(), gasPenalty)
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
			bb.WithBLSMessageAndCode(td.MessageProducer.Transfer(sender, receiver, chain.Value(sendValue)),
				exitcode.SysErrSenderInvalid)
		}
		prevRewards := td.GetRewardSummary()
		prevMinerBalance := td.GetBalance(miner)
		tb.WithBlockBuilder(bb).ApplyAndValidate()
		td.AssertBalance(receiver, acctDefaultBalance)

		newRewards := td.GetRewardSummary()
		newMinerBalance := td.GetBalance(miner)
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(780548)

		// The penalty amount has been burnt by the reward actor, and subtracted from the miner's block reward.
		validateRewards(td, prevRewards, newRewards, prevMinerBalance, newMinerBalance, big.Zero(), gasPenalty)
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
			td.MessageProducer.Transfer(aliceId, builtin.BurntFundsActorAddr, chain.Nonce(1)),
			exitcode.SysErrSenderStateInvalid,
		)

		prevRewards := td.GetRewardSummary()
		prevMinerBalance := td.GetBalance(miner)
		tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		newMinerBalance := td.GetBalance(miner)
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(193137)
		validateRewards(td, prevRewards, newRewards, prevMinerBalance, newMinerBalance, big.Zero(), gasPenalty)
		td.AssertBalance(builtin.BurntFundsActorAddr, gasPenalty)
	})

	t.Run("miner penalty exceeds declared gas limit for BLS message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		alice, _ := td.NewAccountActor(drivers.BLS, acctDefaultBalance)

		gasPrice := int64(2)
		gasPenalty := int64(482274)
		gasLimit := gasPenalty - gasPenalty/gasPrice

		// nonce == 1 causes the message application to fail resulting in a miner penalty.
		bb.WithBLSMessageAndCode(
			td.MessageProducer.Transfer(alice, builtin.BurntFundsActorAddr, chain.Nonce(1), chain.GasPrice(gasPrice), chain.GasLimit(gasLimit)),
			exitcode.SysErrSenderStateInvalid,
		)

		prevRewards := td.GetRewardSummary()
		prevMinerBalance := td.GetBalance(miner)
		tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		newMinerBalance := td.GetBalance(miner)
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		validateRewards(td, prevRewards, newRewards, prevMinerBalance, newMinerBalance, big.Zero(), big.NewInt(gasPenalty))
		td.AssertBalance(builtin.BurntFundsActorAddr, big.NewInt(gasPenalty))
		td.AssertBalance(alice, acctDefaultBalance)
	})

	t.Run("miner penalty exceeds declared gas limit for SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		alice, _ := td.NewAccountActor(drivers.SECP, acctDefaultBalance)

		gasPrice := int64(2)
		gasPenalty := int64(562274)
		gasLimit := gasPenalty - gasPenalty/gasPrice
		// nonce == 1 causes the message application to fail resulting in a miner penalty.
		bb.WithSECPMessageAndCode(
			td.MessageProducer.Transfer(alice, builtin.BurntFundsActorAddr, chain.Nonce(1), chain.GasPrice(gasPrice), chain.GasLimit(gasLimit)),
			exitcode.SysErrSenderStateInvalid,
		)

		prevRewards := td.GetRewardSummary()
		prevMinerBalance := td.GetBalance(miner)
		tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		newMinerBalance := td.GetBalance(miner)
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		validateRewards(td, prevRewards, newRewards, prevMinerBalance, newMinerBalance, big.Zero(), big.NewInt(gasPenalty))
		td.AssertBalance(builtin.BurntFundsActorAddr, big.NewInt(gasPenalty))
		td.AssertBalance(alice, acctDefaultBalance)
	})

	t.Run("no penalty if the balance is not sufficient to cover transfer", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		halfBalance := abi.NewTokenAmount(5_000_000_000)
		_, aliceId := td.NewAccountActor(drivers.BLS, big.Add(halfBalance, halfBalance))

		// Attempt to whole balance, in two parts.
		// The second message should fail (insufficient balance to pay fees).
		msgOk := td.MessageProducer.Transfer(aliceId, builtin.BurntFundsActorAddr, chain.Value(halfBalance))
		msgFail := td.MessageProducer.Transfer(aliceId, builtin.BurntFundsActorAddr, chain.Value(halfBalance), chain.Nonce(1))
		bb.WithBLSMessageOk(msgOk).WithBLSMessageAndCode(msgFail, exitcode.SysErrInsufficientFunds)

		prevRewards := td.GetRewardSummary()
		prevMinerBalance := td.GetBalance(miner)
		result := tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		newMinerBalance := td.GetBalance(miner)
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(0)
		validateRewards(td, prevRewards, newRewards, prevMinerBalance, newMinerBalance, big.Add(result.Receipts[0].GasUsed.Big(), result.Receipts[1].GasUsed.Big()), gasPenalty)

		burn := big.Add(drivers.GetBurn(msgOk.GasLimit, result.Receipts[0].GasUsed), drivers.GetBurn(msgFail.GasLimit, result.Receipts[1].GasUsed))
		td.AssertBalance(builtin.BurntFundsActorAddr, big.Add(burn, big.Add(halfBalance, gasPenalty)))
	})

	t.Run("insufficient gas to cover return value", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)

		alice, _ := td.NewAccountActor(drivers.BLS, acctDefaultBalance)

		// get a successful result so we can determine how much gas it costs. We'll reduce this by 1 in a subsequent call
		// to test insufficient gas to cover return value.
		tracerResult := tb.WithBlockBuilder(
			drivers.NewBlockBuilder(td, td.ExeCtx.Miner).
				WithBLSMessageAndRet(td.MessageProducer.MinerControlAddresses(alice, miner, nil, chain.Nonce(0)),
					// required to satisfy testing methods, unrelated to current test.
					chain.MustSerialize(&miner_spec.GetControlAddressesReturn{
						Owner:  td.StateDriver.BuiltinMinerInfo().OwnerID,
						Worker: td.StateDriver.BuiltinMinerInfo().WorkerID,
					}),
				),
		).ApplyAndValidate()
		requiredGasLimit := tracerResult.Receipts[0].GasUsed

		/* now the test */
		tb.Clear()
		rewardsBefore := td.GetRewardSummary()
		minerBalanceBefore := td.GetBalance(miner)
		senderBalanceBefore := td.GetBalance(alice)
		td.ExeCtx.Epoch++

		// Apply the message again with a reduced gas limit
		// A value just one less than the required limit for success ensures that the gas limit will be reached
		// at the last possible gas charge, which is that for the return value size.
		gasLimit := requiredGasLimit - 1
		result := tb.WithBlockBuilder(
			drivers.NewBlockBuilder(td, td.ExeCtx.Miner).
				WithBLSMessageAndCode(td.MessageProducer.MinerControlAddresses(alice, miner, nil, chain.Nonce(1), chain.GasLimit(int64(gasLimit))),
					exitcode.SysErrOutOfGas,
				),
		).ApplyAndValidate()
		gasUsed := result.Receipts[0].GasUsed
		gasCost := gasLimit.Big() // Gas price is 1
		newRewards := td.GetRewardSummary()

		// Check the actual gas charged is equal to the gas limit rather than the amount consumed up to but excluding
		// the return value which is smaller than the gas limit.
		assert.Equal(t, gasLimit, gasUsed)

		// Check sender charged exactly the max cost.
		assert.Equal(td.T, big.Sub(senderBalanceBefore, gasCost), td.GetBalance(alice))

		// Check the miner earned exactly the max cost (plus block reward).
		thisRwd := big.Add(rewardsBefore.NextPerBlockReward, gasCost)
		assert.Equal(td.T, big.Add(minerBalanceBefore, thisRwd), td.GetBalance(miner))
		assert.Equal(td.T, big.Sub(rewardsBefore.Treasury, rewardsBefore.NextPerBlockReward), newRewards.Treasury)
	})

	// TODO more tests:
	// - miner penalty causes subsequent otherwise-valid message to have wrong nonce (another miner penalty)
	// - miner penalty followed by non-miner penalty with same nonce (in different block)
}

func validateRewards(td *drivers.TestDriver, prevRewards *drivers.RewardSummary, newRewards *drivers.RewardSummary, oldMinerBalance abi.TokenAmount, newMinerBalance abi.TokenAmount, gasReward big.Int, gasPenalty big.Int) {
	rwd := big.Add(big.Sub(prevRewards.NextPerBlockReward, gasPenalty), gasReward)
	assert.Equal(td.T, big.Add(oldMinerBalance, rwd), newMinerBalance)
	assert.Equal(td.T, big.Sub(prevRewards.Treasury, prevRewards.NextPerBlockReward), newRewards.Treasury)
}

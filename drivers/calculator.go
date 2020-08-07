package drivers

import (
	"github.com/filecoin-project/chain-validation/chain/types"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
)

const (
	overuseNum = 11
	overuseDen = 10
)

const baseFee = 100

func GetMinerPenalty(gasLimit int64) big_spec.Int {
	return big_spec.NewInt(baseFee * gasLimit)
}

func GetBurn(gasLimit types.GasUnits, gasUsed types.GasUnits) big_spec.Int {
	over := gasLimit - (overuseNum*gasUsed)/overuseDen
	if over < 0 {
		over = 0
	}
	if over > gasUsed {
		over = gasUsed
	}

	overestimateGas := big_spec.NewInt(int64(gasLimit - gasUsed))
	overestimateGas = big_spec.Mul(overestimateGas, big_spec.NewInt(int64(over)))
	overestimateGas = big_spec.Div(overestimateGas, big_spec.NewInt(int64(gasUsed)))

	totalBurnGas := big_spec.Add(overestimateGas, gasUsed.Big())
	return big_spec.Mul(big_spec.NewInt(baseFee), totalBurnGas)
}

func (d *StateDriver) CalcMessageCost(gasLimit int64, gasPremium big_spec.Int, transferred big_spec.Int, rct types.MessageReceipt) big_spec.Int {
	minerReward := big_spec.Mul(big_spec.NewInt(gasLimit), gasPremium)
	burn := GetBurn(types.GasUnits(gasLimit), rct.GasUsed)
	cost := big_spec.Add(minerReward, burn)

	if rct.ExitCode.IsSuccess() {
		cost = big_spec.Add(cost, transferred)
	}

	return cost
}

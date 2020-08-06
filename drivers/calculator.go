package drivers

import (
	"github.com/filecoin-project/chain-validation/chain/types"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
)

const (
	overuseNum = 11
	overuseDen = 10
)

func GetBurn(gasLimit types.GasUnits, gasUsed types.GasUnits) big_spec.Int {
	over := gasLimit - (overuseNum*gasUsed)/overuseDen
	if over < 0 {
		return big_spec.Zero()
	}
	if over > gasUsed {
		over = gasUsed
	}

	gasToBurn := big_spec.NewInt(int64(gasLimit - gasUsed))
	gasToBurn = big_spec.Mul(gasToBurn, big_spec.NewInt(int64(over)))
	gasToBurn = big_spec.Div(gasToBurn, big_spec.NewInt(int64(gasUsed)))

	return gasToBurn
}

func (d *StateDriver) CalcMessageCost(gasLimit int64, gasPremium big_spec.Int, transferred big_spec.Int, rct types.MessageReceipt) big_spec.Int {
	change := big_spec.Add(rct.GasUsed.Big(), GetBurn(types.GasUnits(gasLimit), rct.GasUsed))
	change = big_spec.Mul(change, gasPremium)
	if rct.ExitCode.IsSuccess() {
		change = big_spec.Add(change, transferred)
	}
	return change
}

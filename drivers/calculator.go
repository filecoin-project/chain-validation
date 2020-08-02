package drivers

import (
	"github.com/filecoin-project/chain-validation/chain/types"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
)

const (
	overuseNum = 3
	overuseDen = 10
)

func GetBurn(gasLimit int64, gasUsed types.GasUnits) big_spec.Int {
	overestimate := int64(gasUsed + (gasUsed * overuseNum / overuseDen))
	if overestimate < gasLimit {
		return big_spec.NewInt(gasLimit - overestimate)
	} else {
		return big_spec.Zero()
	}
}

func (d *StateDriver) CalcMessageCost(gasLimit int64, gasPrice big_spec.Int, transferred big_spec.Int, rct types.MessageReceipt) big_spec.Int {
	change := big_spec.Add(rct.GasUsed.Big(), GetBurn(gasLimit, rct.GasUsed))
	change = big_spec.Mul(change, gasPrice)
	if rct.ExitCode.IsSuccess() {
		change = big_spec.Add(change, transferred)
	}
	return change
}

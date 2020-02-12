package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/filecoin-project/specs-actors/actors/builtin/account"
	"github.com/filecoin-project/specs-actors/actors/builtin/cron"
	init_ "github.com/filecoin-project/specs-actors/actors/builtin/init"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/builtin/reward"
)

// This is a tool used to save some time writing all the message builder methods for actors.
// It produces code that does NOT compile but saves a lot of typing.
func main() {

	accountExports := account.Actor{}.Exports()
	accountDetails := ParseGenerationFields("account_messages", "Account", accountExports)
	f := jen.NewFile(accountDetails.file)
	MakeMethods(f, accountDetails)
	fmt.Printf("%#v", f)

	cronExports := cron.Actor{}.Exports()
	details := ParseGenerationFields("cron_messages", "Cron", cronExports)
	f = jen.NewFile(details.file)
	MakeMethods(f, details)
	fmt.Printf("%#v", f)

	initExports := init_.Actor{}.Exports()
	details = ParseGenerationFields("init_messages", "Init", initExports)
	f = jen.NewFile(details.file)
	MakeMethods(f, details)
	fmt.Printf("%#v", f)

	marketExports := market.Actor{}.Exports()
	details = ParseGenerationFields("market_messages", "Market", marketExports)
	f = jen.NewFile(details.file)
	MakeMethods(f, details)
	fmt.Printf("%#v", f)

	minerExports := miner.Actor{}.Exports()
	details = ParseGenerationFields("miner_messages", "Miner", minerExports)
	f = jen.NewFile(details.file)
	MakeMethods(f, details)
	fmt.Printf("%#v", f)

	multisigExports := multisig.Actor{}.Exports()
	details = ParseGenerationFields("multisig_messages", "Multisig", multisigExports)
	f = jen.NewFile(details.file)
	MakeMethods(f, details)
	fmt.Printf("%#v", f)

	paychExports := paych.Actor{}.Exports()
	details = ParseGenerationFields("paych_messages", "Paych", paychExports)
	f = jen.NewFile(details.file)
	MakeMethods(f, details)
	fmt.Printf("%#v", f)

	powerExports := power.Actor{}.Exports()
	details = ParseGenerationFields("power_messages", "Power", powerExports)
	f = jen.NewFile(details.file)
	MakeMethods(f, details)
	fmt.Printf("%#v", f)

	rewardExports := reward.Actor{}.Exports()
	details = ParseGenerationFields("reward_messages", "Reward", rewardExports)
	f = jen.NewFile(details.file)
	MakeMethods(f, details)
	fmt.Printf("%#v", f)
}

type GenDetails struct {
	file  string
	pairs []MethodParam
}

type MethodParam struct {
	method       string
	methodPrefix string
	params       string
}

func ParseGenerationFields(file string, methodPrefix string, exports []interface{}) GenDetails {
	details := GenDetails{
		file: file,
	}
	for _, m := range exports {
		if m == nil {
			continue
		}

		// ignore things that are not functions
		meth := reflect.ValueOf(m)
		t := meth.Type()
		if t.Kind() != reflect.Func {
			continue
		}

		methodPkgPath := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
		tokens := strings.Split(methodPkgPath, "/")

		pkgStructMethod := strings.Split(tokens[len(tokens)-1], ".")

		method := strings.TrimRight(pkgStructMethod[2], "-fm")

		methodParam := strings.TrimLeft(t.In(1).String(), "*")

		details.pairs = append(details.pairs, MethodParam{
			method:       method,
			methodPrefix: methodPrefix,
			params:       methodParam,
		})
	}
	return details
}

func MakeMethods(jenFile *jen.File, details GenDetails) {
	for _, d := range details.pairs {
		jenFile.Func().Params(
			jen.Id("mp").Id("*MessageProducer"),
		).Id(fmt.Sprintf("%s%s", d.methodPrefix, d.method)).Params(
			jen.Id("to"),
			jen.Id("from").Id("address.Address"),
			jen.Id("params").Id(d.params),
			jen.Id("opts").Id("...MsgOpt"),
		).Params(
			jen.Id("*Message"),
			jen.Id("error"),
		).Block(
			jen.List(jen.Id("ser"), jen.Err()).Op(":=").Id("state.Serialize").Call(jen.Id("&params")),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return(jen.Id("nil"), jen.Id("err")),
			),
			jen.Return(jen.Id("mp.Build").Params(jen.Id("to"), jen.Id("from"), jen.Id(fmt.Sprintf("builtin_spec.Methods%s", d.methodPrefix)), jen.Id("ser"), jen.Id("opts...")), jen.Nil()),
		)
	}
}

package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/dave/jennifer/jen"
	//account_actor "github.com/filecoin-project/specs-actors/actors/builtin/account"
	multisig_actor "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
)

type Field struct {
	Name    string
	Pointer bool
	Type    reflect.Type
	Pkg     string

	IterLabel string
}

type GenTypeInfo struct {
	Name   string
	Fields []Field
}

func nameIsExported(name string) bool {
	return strings.ToUpper(name[0:1]) == name[0:1]
}

func ParseTypeInfo(pkg string, i interface{}) (*GenTypeInfo, error) {
	t := reflect.TypeOf(i)

	out := GenTypeInfo{
		Name: t.Name(),
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !nameIsExported(f.Name) {
			continue
		}

		ft := f.Type
		var pointer bool
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
			pointer = true
		}

		out.Fields = append(out.Fields, Field{
			Name:    f.Name,
			Pointer: pointer,
			Type:    ft,
			Pkg:     pkg,
		})
	}

	return &out, nil
}

var messageProducerTemplate = `
func (mp *MessageProducer) {{.Name}}(to, from address.Address, params {{.Type}}, opts ...MsgOpt) (*Message, error) {
	ser, err := state.Serialize(&params)
	if err != nil {
		return nil, err
	}
	//return mp.Build(to, from, , ser, opts...), nil
}
`

func main() {
	exports := multisig_actor.MultiSigActor{}.Exports()

	for _, m := range exports {
		if m == nil {
			continue
		}
		meth := reflect.ValueOf(m)
		t := meth.Type()
		if t.Kind() != reflect.Func {
			fmt.Printf("this isn't a function!\n")
			continue
		}

		methodPkgPath := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
		tokens := strings.Split(methodPkgPath, "/")

		pkgStructMethod := strings.Split(tokens[len(tokens)-1], ".")
		//fmt.Printf("Package path: %s\n", pkgStructMethod)

		pkg := pkgStructMethod[0]
		//actorType := pkgStructMethod[1]
		method := strings.TrimRight(pkgStructMethod[2], "-fm")

		importPath := strings.TrimRight(methodPkgPath, pkgStructMethod[2]+".")
		importPath = strings.TrimRight(importPath, pkgStructMethod[1])
		importPath = strings.TrimRight(importPath, ".")
		//fmt.Printf("Import Path: %s\n", importPath)

		GenerateMessageMethods3(pkg, importPath, fmt.Sprintf("%s%s", pkgStructMethod[1], method), t.In(1).String())
		//fmt.Printf("Package: %s\nActor: %s\nMethod: %s\n\n", pkg, actorType, method)
		//fmt.Printf("Actor parameters: %s\n", t.In(1))
		break
	}
}

func GenerateMessageMethods(file, importPkg, method, actParam string) {
	f := jen.NewFile(file)
	f.Func().Id(method).Params().Block(
		jen.Qual(importPkg, method).Call(jen.Lit(actParam)),
	)
	fmt.Printf("%#v", f)
}

func GenerateMessageMethods2(file, importPkg, method, actParam string) {
	f := jen.NewFile(file)
	f.Func().Params(
		jen.Id("mp *MessageProducer"),
	).Id(method).Params(
		jen.Id("to"),
		jen.Id("from address.Address"),
		jen.Id(fmt.Sprintf("params %s", actParam)),
		jen.Id("opts ...MsgOpt"),
	).Block(
		jen.Id("ser").Id("err").Op(":=").Id("state.Serialize").Call(
			jen.Id("&params"),
		),
	)
	fmt.Printf("%#v", f)
}

func GenerateMessageMethods3(file, importPkg, method, actParam string) {
	f := jen.NewFile(file)
	f.Func().Params(
		jen.Id("mp").Id("*MessageProducer"),
	).Id(method).Params(
		jen.Id("to"),
		jen.Id("from").Id("address.Address"),
		jen.Id("params").Id(actParam),
		jen.Id("opts").Id("...MsgOpt"),
	).Params(
		jen.Id("*Message"),
		jen.Id("error"),
	).Block(
		jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
			jen.Return(jen.Id("nil"), jen.Id("err")),
		),
	)
	fmt.Printf("%#v", f)
}

//go:generate go run gen.go

// Based on code from https://github.com/wso2/product-apim-tooling/tree/6292e88350868ddaa17ec58501672f008639dd5a/import-export-cli on 2020-03-21 under the Apache License 2.0.

package stateroot

import "github.com/ipfs/go-cid"

type resourceBox struct {
	storage map[string][]cid.Cid
}

func newResourceBox() *resourceBox {
	return &resourceBox{storage: make(map[string][]cid.Cid)}
}

// Find a file
func (r *resourceBox) Has(file string) bool {
	if _, ok := r.storage[file]; ok {
		return true
	}
	return false
}

// Get file's content
func (r *resourceBox) Get(file string) ([]cid.Cid, bool) {
	if f, ok := r.storage[file]; ok {
		return f, ok
	}
	return nil, false
}

// Add a file to boxs
func (r *resourceBox) Add(file string, content []cid.Cid) {
	r.storage[file] = content
}

// Resource expose
var resources = newResourceBox()

// Get a file from boxs
func Get(file string) ([]cid.Cid, bool) {
	return resources.Get(file)
}

// Add a file content to boxs
func Add(file string, content []cid.Cid) {
	resources.Add(file, content)
}

// Has a file in boxs
func Has(file string) bool {
	return resources.Has(file)
}

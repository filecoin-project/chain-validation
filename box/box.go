//go:generate go run gen.go

// Source: https://levelup.gitconnected.com/how-i-embedded-resources-in-go-514b72f6ef0a

package box

type resourceBox struct {
	storage map[string][]byte
}

func newResourceBox() *resourceBox {
	return &resourceBox{storage: make(map[string][]byte)}
}

// Find a file
func (r *resourceBox) Has(file string) bool {
	if _, ok := r.storage[file]; ok {
		return true
	}
	return false
}

// Get file's content
func (r *resourceBox) Get(file string) ([]byte, bool) {
	if f, ok := r.storage[file]; ok {
		return f, ok
	}
	return nil, false
}

// Add a file to box
func (r *resourceBox) Add(file string, content []byte) {
	r.storage[file] = content
}

// Resource expose
var resources = newResourceBox()

// Get a file from box
func Get(file string) ([]byte, bool) {
	return resources.Get(file)
}

// Add a file content to box
func Add(file string, content []byte) {
	resources.Add(file, content)
}

// Has a file in box
func Has(file string) bool {
	return resources.Has(file)
}

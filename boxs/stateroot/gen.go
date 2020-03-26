//+build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/ipfs/go-cid"
)

const blob = "blob.go"

var packageTemplate = template.Must(template.New("").Funcs(map[string]interface{}{"conv": FormatCIDSlice}).Parse(`// Code generated by go generate; DO NOT EDIT.
// generated using files from resources directory
package boxs
func init(){
	{{- range $name, $file := . }}
    	resources.Add("{{ $name }}", []int64{ {{ conv $file }} })
	{{- end }}
}
`))

// don't think this will work
func FormatCIDSlice(sl []cid.Cid) string {
	builder := strings.Builder{}
	for _, v := range sl {
		builder.WriteString(fmt.Sprintf("%d,", v.String()))
	}
	return builder.String()
}

func main() {
	log.Println("Baking resources... \U0001F4E6")

	if _, err := os.Stat("resources"); os.IsNotExist(err) {
		log.Fatal("Resources directory does not exists")
	}

	resources := make(map[string][]int64)
	err := filepath.Walk("../resources", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error :", err)
			return err
		}
		relativePath := filepath.ToSlash(strings.TrimPrefix(path, "resources"))
		if info.IsDir() {
			log.Println(path, "is a directory, skipping... \U0001F47B")
			return nil
		} else {
			log.Println(path, "is a file, baking in... \U0001F31F")
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				value, err := strconv.ParseInt(scanner.Text(), 10, 64)
				if err != nil {
					return err
				}
				resources[relativePath] = append(resources[relativePath], value)
			}

		}
		return nil
	})

	if err != nil {
		log.Fatal("Error walking through resources directory:", err)
	}

	f, err := os.Create(blob)
	if err != nil {
		log.Fatal("Error creating blob file:", err)
	}
	defer f.Close()

	builder := &bytes.Buffer{}

	err = packageTemplate.Execute(builder, resources)
	if err != nil {
		log.Fatal("Error executing template", err)
	}

	data, err := format.Source(builder.Bytes())
	if err != nil {
		log.Fatal("Error formatting generated code", err)
	}
	err = ioutil.WriteFile(blob, data, os.ModePerm)
	if err != nil {
		log.Fatal("Error writing blob file", err)
	}

	log.Println("Baking resources done... \U0001F680")
}

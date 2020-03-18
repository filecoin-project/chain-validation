package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/* This program reads all files in ./gas_files which it expects to be comma separated values of the form:
*	oldStateCID,messageCID,newStateCID,GasUsed
*
*  It then generates a file ./gas/constants.go that contains a go map of the form:
*  var GasConstants = map[string]int64{
*      "oldStateCID-messageCID-newStateCID: GasUsed",
*      ...
*      ...
*  }
*
*  This go file may then be used by validation to ensure the correct amount of gas is being charged.
*
 */
func main() {
	var gasPath string
	var codePath string
	flag.StringVar(&gasPath, "gaspath", "./gas_files", "sets location where to read gas files from, expect the directory to only contain gasfiles")
	flag.StringVar(&codePath, "codepath", "./gas/constants.go", "sets location where to write go code to, generates a map")

	// exit on error is set
	flag.Parse()

	if gasPath == "" {
		panic("must provide a gaspath")
	}

	// get a slice of all the gas files
	gasFiles, err := getGasFilesFromPath(gasPath)
	if err != nil {
		panic(err)
	}

	code, err := generateGoCodeFromGasFiles(gasFiles)
	if err != nil {
		panic(err)
	}

	var codeFile *os.File
	if codePath == "" {
		codeFile = os.Stdout
	} else {
		codeFile, err = os.OpenFile(codePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer codeFile.Close()
	}

	_, err = codeFile.WriteString(code)
	if err != nil {
		panic(err)
	}
}

func generateGoCodeFromGasFiles(gasFiles []string) (string, error) {
	var sb strings.Builder
	if _, err := sb.WriteString("package gas\n\n"); err != nil {
		return "", err
	}
	if _, err := sb.WriteString("var GasConstants = map[string]int64{\n"); err != nil {
		return "", err
	}

	seen := make(map[string]bool)
	for _, gf := range gasFiles {
		gasFile, err := os.OpenFile(gf, os.O_RDONLY, 0644)
		if err != nil {
			return "", err
		}
		// read the gas file line by line
		scanner := bufio.NewScanner(gasFile)
		for scanner.Scan() {
			l, err := gasFileLineToMapEntry(scanner.Text())
			if err != nil {
				return "", err
			}

			// no duplicates
			if !seen[l] {
				if _, err := sb.WriteString(l); err != nil {
					return "", err
				}
			}
			seen[l] = true
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	}
	sb.WriteString("}")
	return sb.String(), nil

}

func gasFileLineToMapEntry(line string) (string, error) {
	tokens := strings.Split(line, ",")
	// expect tokens to always be length 4 (oldState,msgCid,newState,gasCost)
	if len(tokens) != 4 {
		return "", fmt.Errorf("invalid gas line, expected 4 tokens, got %d: %s", len(tokens), tokens)
	}

	oldState := tokens[0]
	msgCid := tokens[1]
	// TODO this should be a value
	//newState := tokens[2]
	gasUnits := tokens[3]
	return fmt.Sprintf("\t\"%s-%s\": %s,\n", oldState, msgCid, gasUnits), nil
}

func getGasFilesFromPath(root string) ([]string, error) {
	var gasFiles []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		gasFiles = append(gasFiles, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return gasFiles[1:], nil
}

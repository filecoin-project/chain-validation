package main

import (
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	logging "github.com/ipfs/go-log"

	"github.com/filecoin-project/chain-validation/client"
	"github.com/filecoin-project/chain-validation/client/services"
	"github.com/filecoin-project/chain-validation/suites"
)

func init() {
	logging.SetAllLoggers(logging.LevelInfo)
}

// Env vars
const (
	Env_Host    = "CHAIN_VALIDATION_HOST"
	Env_Post    = "CHAIN_VALIDATION_PORT"
	Env_Timeout = "CHAIN_VALIDATION_TIMEOUT"
)

var (
	host    string
	port    string
	timeout time.Duration
)

func init() {
	host = os.Getenv(Env_Host)
	if host == "" {
		host = "127.0.0.1"
	}
	port = os.Getenv(Env_Post)
	if port == "" {
		port = "8378"
	}
	timeoutStr := os.Getenv(Env_Timeout)
	var err error
	if timeoutStr == "" {
		timeout = 3 * time.Second
	} else {
		timeout, err = time.ParseDuration(timeoutStr)
		if err != nil {
			panic(err)
		}
	}

}

func TestChainValidationMessageSuite(t *testing.T) {
	cfg := client.Config{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}
	handler := services.NewServiceHandler(client.NewRpcClient(cfg))

	for _, testCase := range suites.MessageTestCases() {
		t.Run(caseName(testCase), func(t *testing.T) {
			testCase(t, handler)
		})
	}
}

func TestChainValidationTipSetSuite(t *testing.T) {
	cfg := client.Config{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}
	handler := services.NewServiceHandler(client.NewRpcClient(cfg))
	for _, testCase := range suites.TipSetTestCases() {
		t.Run(caseName(testCase), func(t *testing.T) {
			testCase(t, handler)
		})
	}
}

func caseName(testCase suites.TestCase) string {
	fqName := runtime.FuncForPC(reflect.ValueOf(testCase).Pointer()).Name()
	toks := strings.Split(fqName, ".")
	return toks[len(toks)-1]
}

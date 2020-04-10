package config

import (
	"encoding/json"

	"github.com/filecoin-project/chain-validation/client"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("service/config")

const (
	Method_GetConfig = "ConfigService.Config"
)

type ConfigReply struct {
	TrackGas         bool `json:"trackGas"`
	CheckGas         bool `json:"checkGas"`
	CheckExitCode    bool `json:"checkExitCode"`
	CheckReturnValue bool `json:"checkReturnValue"`
	CheckStateRoot   bool `json:"checkStateRoot"`

	TestSuite []string `json:"testSuite"`
}

func NewConfigService(rpcClient *client.RpcClient) *ConfigService {
	return &ConfigService{rpcClient: rpcClient}
}

type ConfigService struct {
	rpcClient *client.RpcClient
}

func (cs *ConfigService) GetConfig() (*ConfigReply, error) {
	resp, err := cs.rpcClient.Do(Method_GetConfig, nil)
	if err != nil {
		return nil, err
	}
	log.Debugw(Method_GetConfig, "response", resp)

	var out ConfigReply
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

package client

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"time"

	jsonrpc "github.com/gorilla/rpc/v2/json"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("client/rpc")

type Config struct {
	Host string
	Port string

	Timeout time.Duration
}

type RpcClient struct {
	httpclient *http.Client
	config     Config
}

func NewRpcClient(cfg Config) *RpcClient {
	log.Debugw("NewClient", "config", cfg)
	httpclient := &http.Client{Timeout: cfg.Timeout}
	return &RpcClient{httpclient: httpclient, config: cfg}
}

func (c *RpcClient) Do(method string, args interface{}) (json.RawMessage, error) {
	log.Debugw("Do", "method", method, "args", args)

	encReq, err := jsonrpc.EncodeClientRequest(method, args)
	if err != nil {
		return nil, err
	}

	uri := "http://" + net.JoinHostPort(c.config.Host, c.config.Port) + "/rpc"
	req, err := http.NewRequest("POST", uri, bytes.NewReader(encReq))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var out json.RawMessage
	if err := jsonrpc.DecodeClientResponse(resp.Body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

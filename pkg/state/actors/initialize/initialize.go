package initialize

import (
	"github.com/ipfs/go-cid"
)

type ExecParams struct {
	Code   cid.Cid
	Params []byte
}

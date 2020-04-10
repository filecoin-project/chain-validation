package applier

import (
	"encoding/json"

	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/client"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("service/applier")

const (
	Method_ApplyMessage        = "VmWrapperService.ApplyMessage"
	Method_ApplySignedMessage  = "VmWrapperService.ApplySignedMessage"
	Method_ApplyTipSetMessages = "VmWrapperService.ApplyTipSetMessages"
)

func NewApplierService(rpcclient *client.RpcClient) *ApplierService {
	return &ApplierService{rpcClient: rpcclient}
}

type ApplierService struct {
	rpcClient *client.RpcClient
}

type ApplyMessageReply struct {
	Receipt types.MessageReceipt
	Penalty abi.TokenAmount
	Reward  abi.TokenAmount
	Root    cid.Cid
}

type ApplyMessageArgs struct {
	Epoch   abi.ChainEpoch
	Message *types.Message
}

func (as *ApplierService) ApplyMessage(epoch abi.ChainEpoch, msg *types.Message) (*ApplyMessageReply, error) {
	resp, err := as.rpcClient.Do(Method_ApplyMessage, &ApplyMessageArgs{
		Epoch:   epoch,
		Message: msg,
	})
	if err != nil {
		return nil, err
	}
	log.Debugw(Method_ApplyMessage, "response", resp)

	var out ApplyMessageReply
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, err
}

type ApplySignedMessageArgs struct {
	Epoch         abi.ChainEpoch
	SignedMessage *types.SignedMessage
}

func (as *ApplierService) ApplySignedMessage(epoch abi.ChainEpoch, smsg *types.SignedMessage) (*ApplyMessageReply, error) {
	resp, err := as.rpcClient.Do(Method_ApplySignedMessage, &ApplySignedMessageArgs{
		Epoch:         epoch,
		SignedMessage: smsg,
	})
	if err != nil {
		return nil, err
	}
	log.Debugw(Method_ApplySignedMessage, "response", resp)

	var out ApplyMessageReply
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, err
}

type ApplyTipSetMessagesArgs struct {
	Epoch      abi.ChainEpoch
	Blocks     []types.BlockMessagesInfo
	Randomness abi.Randomness
}

type ApplyTipSetMessagesReply struct {
	Receipts []types.MessageReceipt
	Root     cid.Cid
}

func (as *ApplierService) ApplyTipSetMessages(epoch abi.ChainEpoch, blocks []types.BlockMessagesInfo, rand abi.Randomness) (*ApplyTipSetMessagesReply, error) {
	resp, err := as.rpcClient.Do(Method_ApplyTipSetMessages, &ApplyTipSetMessagesArgs{
		Epoch:      epoch,
		Randomness: rand,
		Blocks:     blocks,
	})
	if err != nil {
		return nil, err
	}
	log.Debugw(Method_ApplyTipSetMessages, "response", resp)

	var out ApplyTipSetMessagesReply
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, err
}

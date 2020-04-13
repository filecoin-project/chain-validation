package vmwrapper

import (
	"bytes"
	"encoding/json"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/runtime"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/client"
)

var log = logging.Logger("service/vmwrapper")

const (
	// create a new vm instance
	Method_NewVM = "VmWrapperService.NewVM"

	// state inspection and modification methods methods
	Method_Root          = "VmWrapperService.Root"
	Method_StoreGet      = "VmWrapperService.StoreGet"
	Method_StorePut      = "VmWrapperService.StorePut"
	Method_Actor         = "VmWrapperService.Actor"
	Method_SetActorState = "VmWrapperService.SetActorState"
	Method_CreateActor   = "VmWrapperService.CreateActor"

	// message application methods
	Method_ApplyMessage        = "VmWrapperService.ApplyMessage"
	Method_ApplySignedMessage  = "VmWrapperService.ApplySignedMessage"
	Method_ApplyTipSetMessages = "VmWrapperService.ApplyTipSetMessages"
)

func NewVmWrapperService(client *client.RpcClient) *VmWrapperService {
	return &VmWrapperService{rpcClient: client}
}

type VmWrapperService struct {
	rpcClient *client.RpcClient
}

func (vs *VmWrapperService) New() error {
	resp, err := vs.rpcClient.Do(Method_NewVM, nil)
	if err != nil {
		return err
	}
	log.Debugw(Method_NewVM, "response", resp)
	return nil
}

type RootReply struct {
	Root cid.Cid
}

func (vs *VmWrapperService) Root() (cid.Cid, error) {
	resp, err := vs.rpcClient.Do(Method_Root, nil)
	if err != nil {
		return cid.Undef, err
	}
	log.Debugw(Method_Root, "response", resp)

	var out RootReply
	if err := json.Unmarshal(resp, &out); err != nil {
		return cid.Undef, err
	}
	return out.Root, nil
}

type StoreGetArgs struct {
	Key cid.Cid
}

type StoreGetReply struct {
	Out []byte
}

func (vs *VmWrapperService) StoreGet(key cid.Cid, out runtime.CBORUnmarshaler) error {
	resp, err := vs.rpcClient.Do(Method_StoreGet, &StoreGetArgs{Key: key})
	if err != nil {
		return err
	}
	log.Debugw(Method_StoreGet, "response", resp)

	var tmp StoreGetReply
	if err := json.Unmarshal(resp, &tmp); err != nil {
		return err
	}
	if err := out.UnmarshalCBOR(bytes.NewReader(tmp.Out)); err != nil {
		panic(err)
	}
	return nil
}

type StorePutArgs struct {
	Value []byte
}

type StorePutReply struct {
	Key cid.Cid
}

func (vs *VmWrapperService) StorePut(value runtime.CBORMarshaler) (cid.Cid, error) {
	raw := chain.MustSerialize(value)
	resp, err := vs.rpcClient.Do(Method_StorePut, &StorePutArgs{Value: raw})
	if err != nil {
		return cid.Undef, err
	}
	log.Debugw(Method_StorePut, "response", resp)

	var out StorePutReply
	if err := json.Unmarshal(resp, &out); err != nil {
		return cid.Undef, err
	}

	return out.Key, nil
}

type ActorArgs struct {
	Addr address.Address
}

type ActorReply struct {
	Code       cid.Cid
	Head       cid.Cid
	CallSeqNum uint64
	Balance    big.Int
}

func (vs *VmWrapperService) Actor(addr address.Address) (*ActorReply, error) {
	resp, err := vs.rpcClient.Do(Method_Actor, &ActorArgs{Addr: addr})
	if err != nil {
		return nil, err
	}
	log.Debugw(Method_Actor, "response", resp)

	var out ActorReply
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

type SetActorStateArgs struct {
	Addr    address.Address
	Balance abi.TokenAmount
	State   []byte
}

func (vs *VmWrapperService) SetActorState(addr address.Address, balance abi.TokenAmount, state runtime.CBORMarshaler) (*ActorReply, error) {
	resp, err := vs.rpcClient.Do(Method_SetActorState, &SetActorStateArgs{
		Addr:    addr,
		Balance: balance,
		State:   chain.MustSerialize(state),
	})
	if err != nil {
		return nil, err
	}
	log.Debugw(Method_SetActorState, "response", resp)

	var out ActorReply
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

type CreateActorArgs struct {
	Code    cid.Cid
	Addr    address.Address
	Balance abi.TokenAmount
	State   []byte // TODO see if you can make this CBORBytes instead
}

type CreateActorReply struct {
	Addr  address.Address
	Actor *ActorReply
}

func (vs *VmWrapperService) CreateActor(code cid.Cid, addr address.Address, balance abi.TokenAmount, state runtime.CBORMarshaler) (*CreateActorReply, error) {
	resp, err := vs.rpcClient.Do(Method_CreateActor, &CreateActorArgs{
		Code:    code,
		Addr:    addr,
		Balance: balance,
		State:   chain.MustSerialize(state),
	})
	if err != nil {
		return nil, err
	}
	log.Debugw(Method_CreateActor, "response", resp)

	var out CreateActorReply
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}

	return &out, nil
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

func (vs *VmWrapperService) ApplyMessage(epoch abi.ChainEpoch, msg *types.Message) (*ApplyMessageReply, error) {
	resp, err := vs.rpcClient.Do(Method_ApplyMessage, &ApplyMessageArgs{
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

func (vs *VmWrapperService) ApplySignedMessage(epoch abi.ChainEpoch, smsg *types.SignedMessage) (*ApplyMessageReply, error) {
	resp, err := vs.rpcClient.Do(Method_ApplySignedMessage, &ApplySignedMessageArgs{
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

func (vs *VmWrapperService) ApplyTipSetMessages(epoch abi.ChainEpoch, blocks []types.BlockMessagesInfo, rand abi.Randomness) (*ApplyTipSetMessagesReply, error) {
	resp, err := vs.rpcClient.Do(Method_ApplyTipSetMessages, &ApplyTipSetMessagesArgs{
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

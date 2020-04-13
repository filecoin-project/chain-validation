package services

import (
	"fmt"
	"math/rand"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-crypto"
	"github.com/filecoin-project/specs-actors/actors/abi"
	acrypto "github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	"github.com/minio/blake2b-simd"

	// TODO reimplement these
	lotustypes "github.com/filecoin-project/lotus/chain/types"
	lotuswallet "github.com/filecoin-project/lotus/chain/wallet"

	"github.com/filecoin-project/chain-validation/chain/types"
	"github.com/filecoin-project/chain-validation/client"
	"github.com/filecoin-project/chain-validation/client/services/config"
	"github.com/filecoin-project/chain-validation/client/services/vmwrapper"
	"github.com/filecoin-project/chain-validation/state"
)

var log = logging.Logger("service/handler")

var _ state.VMWrapper = (*ServiceHandler)(nil)
var _ state.Applier = (*ServiceHandler)(nil)
var _ state.Factories = (*ServiceHandler)(nil)

func NewServiceHandler(client *client.RpcClient) *ServiceHandler {
	return &ServiceHandler{
		vm:     vmwrapper.NewVmWrapperService(client),
		config: config.NewConfigService(client),
	}
}

type ServiceHandler struct {
	vm     *vmwrapper.VmWrapperService
	config *config.ConfigService
}

//
// Impl Factories interface
//

func (s *ServiceHandler) NewStateAndApplier() (state.VMWrapper, state.Applier) {
	// sanity
	if s.vm == nil {
		panic("call new service handler first")
	}
	return s, s
}

func (s *ServiceHandler) NewKeyManager() state.KeyManager {
	if s.vm == nil {
		panic("call new service handler first")
	}
	return newKeyManager()
}

func (s *ServiceHandler) NewValidationConfig() state.ValidationConfig {
	if s.vm == nil {
		panic("call new service handler first")
	}

	cfg, err := s.config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	return &configWrapper{cfg: cfg}
}

type configWrapper struct {
	cfg *config.ConfigReply
}

func (c configWrapper) ValidateGas() bool {
	return c.cfg.CheckGas
}

func (c configWrapper) ValidateExitCode() bool {
	return c.cfg.CheckExitCode
}

func (c configWrapper) ValidateReturnValue() bool {
	return c.cfg.CheckReturnValue
}

func (c configWrapper) ValidateStateRoot() bool {
	return c.cfg.CheckStateRoot
}

//
// Impl VMWrapper interface
//

func (s *ServiceHandler) New() {
	if err := s.vm.New(); err != nil {
		log.Fatal(err)
	}
}

func (s *ServiceHandler) Root() cid.Cid {
	root, err := s.vm.Root()
	if err != nil {
		log.Fatal(err)
	}
	return root
}

func (s *ServiceHandler) StoreGet(key cid.Cid, out runtime.CBORUnmarshaler) error {
	return s.vm.StoreGet(key, out)
}

func (s *ServiceHandler) StorePut(value runtime.CBORMarshaler) (cid.Cid, error) {
	return s.vm.StorePut(value)
}

func (s *ServiceHandler) Actor(address address.Address) (state.Actor, error) {
	reply, err := s.vm.Actor(address)
	if err != nil {
		return nil, err
	}
	return &actorWrapper{reply}, nil
}

func (s *ServiceHandler) SetActorState(addr address.Address, balance abi.TokenAmount, state runtime.CBORMarshaler) (state.Actor, error) {
	reply, err := s.vm.SetActorState(addr, balance, state)
	if err != nil {
		return nil, err
	}
	return &actorWrapper{reply}, nil
}

func (s *ServiceHandler) CreateActor(code cid.Cid, addr address.Address, balance abi.TokenAmount, state runtime.CBORMarshaler) (state.Actor, address.Address, error) {
	reply, err := s.vm.CreateActor(code, addr, balance, state)
	if err != nil {
		return nil, address.Undef, err
	}
	return &actorWrapper{reply.Actor}, reply.Addr, nil
}

//
// Impl Applier interface
//

func (s *ServiceHandler) ApplyMessage(epoch abi.ChainEpoch, msg *types.Message) (types.ApplyMessageResult, error) {
	reply, err := s.vm.ApplyMessage(epoch, msg)
	if err != nil {
		return types.ApplyMessageResult{}, err
	}
	return types.ApplyMessageResult{
		Receipt: reply.Receipt,
		Penalty: reply.Penalty,
		Reward:  reply.Reward,
		Root:    reply.Root.String(),
	}, nil

}

func (s *ServiceHandler) ApplySignedMessage(epoch abi.ChainEpoch, msg *types.SignedMessage) (types.ApplyMessageResult, error) {
	reply, err := s.vm.ApplySignedMessage(epoch, msg)
	if err != nil {
		return types.ApplyMessageResult{}, err
	}
	return types.ApplyMessageResult{
		Receipt: reply.Receipt,
		Penalty: reply.Penalty,
		Reward:  reply.Reward,
		Root:    reply.Root.String(),
	}, nil
}

// TODO the RandomnessSource is going to be tricky to do over RPC
func (s *ServiceHandler) ApplyTipSetMessages(epoch abi.ChainEpoch, blocks []types.BlockMessagesInfo, rnd state.RandomnessSource) (types.ApplyTipSetResult, error) {
	reply, err := s.vm.ApplyTipSetMessages(epoch, blocks, nil)
	if err != nil {
		return types.ApplyTipSetResult{}, err
	}
	return types.ApplyTipSetResult{
		Receipts: reply.Receipts,
		Root:     reply.Root.String(),
	}, nil
}

//
// KeyManager
//

// XXX: lazily use the lotus wallet package as I don't feel like reimplementing it
type KeyManager struct {
	// Private keys by address
	keys map[address.Address]*lotuswallet.Key

	// Seed for deterministic secp key generation.
	secpSeed int64
	// Seed for deterministic bls key generation.
	blsSeed int64 // nolint: structcheck
}

func newKeyManager() *KeyManager {
	return &KeyManager{
		keys:     make(map[address.Address]*lotuswallet.Key),
		secpSeed: 0,
	}
}

func (k *KeyManager) NewSECP256k1AccountAddress() address.Address {
	secpKey := k.newSecp256k1Key()
	k.keys[secpKey.Address] = secpKey
	return secpKey.Address
}

func (k *KeyManager) NewBLSAccountAddress() address.Address {
	blsKey := k.newBLSKey()
	k.keys[blsKey.Address] = blsKey
	return blsKey.Address
}

func (k *KeyManager) Sign(addr address.Address, data []byte) (acrypto.Signature, error) {
	ki, ok := k.keys[addr]
	if !ok {
		return acrypto.Signature{}, fmt.Errorf("unknown address %v", addr)
	}
	var sigType acrypto.SigType
	if ki.Type == lotuswallet.KTSecp256k1 {
		sigType = acrypto.SigTypeBLS
		hashed := blake2b.Sum256(data)
		sig, err := crypto.Sign(ki.PrivateKey, hashed[:])
		if err != nil {
			return acrypto.Signature{}, err
		}

		return acrypto.Signature{
			Type: sigType,
			Data: sig,
		}, nil
	} else if ki.Type == lotuswallet.KTBLS {
		panic("lotus validator cannot sign BLS messages")
	} else {
		panic("unknown signature type")
	}

}

func (k *KeyManager) newSecp256k1Key() *lotuswallet.Key {
	randSrc := rand.New(rand.NewSource(k.secpSeed))
	prv, err := crypto.GenerateKeyFromSeed(randSrc)
	if err != nil {
		panic(err)
	}
	k.secpSeed++
	key, err := lotuswallet.NewKey(lotustypes.KeyInfo{
		Type:       lotuswallet.KTSecp256k1,
		PrivateKey: prv,
	})
	if err != nil {
		panic(err)
	}
	return key
}

func (k *KeyManager) newBLSKey() *lotuswallet.Key {
	// FIXME: bls needs deterministic key generation
	//sk := ffi.PrivateKeyGenerate(s.blsSeed)
	// s.blsSeed++
	sk := [32]byte{}
	sk[0] = uint8(k.blsSeed) // hack to keep gas values determinist
	k.blsSeed++
	key, err := lotuswallet.NewKey(lotustypes.KeyInfo{
		Type:       lotuswallet.KTBLS,
		PrivateKey: sk[:],
	})
	if err != nil {
		panic(err)
	}
	return key
}

//
// Actor
//

type actorWrapper struct {
	actor *vmwrapper.ActorReply
}

func (a *actorWrapper) Code() cid.Cid {
	return a.actor.Code
}
func (a *actorWrapper) Head() cid.Cid {
	return a.actor.Head
}
func (a *actorWrapper) CallSeqNum() uint64 {
	return a.actor.CallSeqNum
}
func (a *actorWrapper) Balance() abi.TokenAmount {
	return a.actor.Balance
}

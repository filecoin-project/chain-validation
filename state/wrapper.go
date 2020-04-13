package state

import (
	cid "github.com/ipfs/go-cid"

	address "github.com/filecoin-project/go-address"

	abi "github.com/filecoin-project/specs-actors/actors/abi"
	big "github.com/filecoin-project/specs-actors/actors/abi/big"
	crypto "github.com/filecoin-project/specs-actors/actors/crypto"
	runtime "github.com/filecoin-project/specs-actors/actors/runtime"
)

// VMWrapper abstracts the inspection and mutation of an implementation-specific state tree and storage.
// The interface wraps a single, mutable state.
type VMWrapper interface {
	// Instantiate a new VM
	NewVM()

	// Returns the CID of the root node of the state tree.
	Root() cid.Cid

	// Get the value at key from vm store
	StoreGet(key cid.Cid, out runtime.CBORUnmarshaler) error

	// Put `value` into vm store
	StorePut(value runtime.CBORMarshaler) (cid.Cid, error)

	// Returns the actor state at `address` (or an error if there is none).
	Actor(address address.Address) (Actor, error)

	// Set state on an actor in the state tree. The actors head is the cid of `state`.
	SetActorState(addr address.Address, balance abi.TokenAmount, state runtime.CBORMarshaler) (Actor, error)

	// Installs a new actor in the state tree, going through the init actor when appropriate and returning the ID address of the actor.
	CreateActor(code cid.Cid, addr address.Address, balance abi.TokenAmount, state runtime.CBORMarshaler) (Actor, address.Address, error)
}

// TODO this needs to be implemented by chain validation. Providing these methods over RPC doesn't add a lot of value.
type KeyManager interface {
	// Creates a new secp private key and returns the associated address.
	NewSECP256k1AccountAddress() address.Address

	// Creates a new BLS private key and returns the associated address.
	NewBLSAccountAddress() address.Address

	// Sign data with addr's key.
	Sign(addr address.Address, data []byte) (crypto.Signature, error)
}

// Actor is an abstraction over the actor states stored in the root of the state tree.
type Actor interface {
	Code() cid.Cid
	Head() cid.Cid
	CallSeqNum() uint64
	Balance() big.Int
}

type ValidationConfig interface {
	ValidateGas() bool
	ValidateExitCode() bool
	ValidateReturnValue() bool
	ValidateStateRoot() bool
}

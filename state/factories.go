package state

import "github.com/filecoin-project/specs-actors/actors/runtime"

// Factories wraps up all the implementation-specific integration points.
type Factories interface {
	NewStateAndApplier(syscalls runtime.Syscalls) (VMWrapper, Applier)

	NewKeyManager() KeyManager

	NewValidationConfig() ValidationConfig
}

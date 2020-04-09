package state

// Factories wraps up all the implementation-specific integration points.
type Factories interface {
	NewStateAndApplier() (VMWrapper, Applier)

	NewKeyManager() KeyManager

	NewValidationConfig() ValidationConfig
}

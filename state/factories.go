package state

// Factories wraps up all the implementation-specific integration points.
type Factories interface {
	NewState() VMWrapper

	NewKeyManager() KeyManager

	Applier
}

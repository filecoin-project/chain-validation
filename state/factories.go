package state

// Factories wraps up all the implementation-specific integration points.
type Factories interface {
	NewState() Wrapper

	NewWallet() Wallet

	Applier
}

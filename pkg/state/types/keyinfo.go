package types

// KeyInfo is used for storing keys in KeyStore
type KeyInfo struct {
	Type       string
	PrivateKey []byte
}

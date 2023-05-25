package contract

// ISecretEncryptor - secret encryptor interface
type ISecretEncryptor interface {
	Encrypt([]byte) ([]byte, error)
}

// ISecretDecryptor - secret decryptor interface
type ISecretDecryptor interface {
	Decrypt([]byte) ([]byte, error)
}

// ISecretCrypto - secret crypto interface
type ISecretCrypto interface {
	ISecretEncryptor
	ISecretDecryptor
}

package keeper

import "github.com/itksb/go-secret-keeper/pkg/contract"

// IClientKeeper - client keeper interface
type IClientKeeper interface {
	contract.IKeeper
}

type ISecretEncryptor interface {
	Encrypt([]byte) ([]byte, error)
}

type ISecretDecryptor interface {
	Decrypt([]byte) ([]byte, error)
}

type ISecretCrypto interface {
	ISecretEncryptor
	ISecretDecryptor
}

type GetPrivateKeyFunc func() ([]byte, error)

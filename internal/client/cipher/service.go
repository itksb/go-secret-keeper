package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/itksb/go-secret-keeper/internal/client/keeper"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"io"
)

var _ contract.ISecretCrypto = &CryptoService{}

// CryptoService - crypto service interface
type CryptoService struct {
	getPrivateKeyFunc keeper.GetPrivateKeyFunc
}

// NewCryptoService - constructor. Creates new crypto service
func NewCryptoService(
	getPrivateKeyFunc keeper.GetPrivateKeyFunc,
) *CryptoService {
	return &CryptoService{
		getPrivateKeyFunc: getPrivateKeyFunc,
	}
}

// Encrypt - encrypt data
func (s *CryptoService) Encrypt(bytes []byte) ([]byte, error) {
	key, err := s.getPrivateKeyFunc()
	if err != nil {
		return nil, err
	}

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nonce, nonce, bytes, nil), nil
}

// Decrypt - decrypt data
func (s *CryptoService) Decrypt(bytes []byte) ([]byte, error) {
	key, err := s.getPrivateKeyFunc()
	if err != nil {
		return nil, err
	}

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	nonce, dst := bytes[:nonceSize], bytes[nonceSize:]

	return aesgcm.Open(nil, nonce, dst, nil)
}

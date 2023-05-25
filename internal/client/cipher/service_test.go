package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestEncrypt - test encryption
// These test cases verify that the Encrypt function correctly encrypts the data using AES-GCM encryption.
// It compares the decrypted data with the original data to ensure the encryption
// and decryption process is working correctly.
func TestEncrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef") // 32 bytes key
	data := []byte("hello, world!")

	// Create a mock implementation of GetPrivateKeyFunc
	getPrivateKeyFunc := func() ([]byte, error) {
		return key, nil
	}

	// Create a new CryptoService instance
	cryptoService := NewCryptoService(getPrivateKeyFunc)

	// Encrypt the data
	encrypted, err := cryptoService.Encrypt(data)
	assert.NoError(t, err, "Encryption should not return an error")

	// Create a new AES cipher using the key
	aesBlock, err := aes.NewCipher(key)
	assert.NoError(t, err, "Failed to create AES cipher")

	// Create a new AES-GCM cipher using the AES cipher
	aesGCM, err := cipher.NewGCM(aesBlock)
	assert.NoError(t, err, "Failed to create AES-GCM cipher")

	// Get the nonce size from the AES-GCM cipher
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce and ciphertext from the encrypted data
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	// Decrypt the ciphertext using the AES-GCM cipher
	decrypted, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	assert.NoError(t, err, "Decryption should not return an error")

	// Check if the decrypted data matches the original data
	assert.Equal(t, data, decrypted, "Decrypted data should match the original data")
}

// TestCryptoService_Decrypt - test decryption
// These test cases verify that the Decrypt function correctly decrypts the encrypted data using AES-GCM decryption.
// It compares the decrypted data with the original data to ensure the encryption and decryption process is working correctly.
func TestCryptoService_Decrypt(t *testing.T) {

	key := []byte("0123456789abcdef0123456789abcdef") // 32 bytes key
	originalData := []byte("hello, world!")

	// Create a mock implementation of GetPrivateKeyFunc
	getPrivateKeyFunc := func() ([]byte, error) {
		return key, nil
	}

	// Create a new CryptoService instance
	cryptoService := NewCryptoService(getPrivateKeyFunc)

	// Encrypt the original data
	encrypted, err := cryptoService.Encrypt(originalData)
	assert.NoError(t, err, "Encryption should not return an error")

	// Decrypt the encrypted data
	decrypted, err := cryptoService.Decrypt(encrypted)
	assert.NoError(t, err, "Decryption should not return an error")

	// Check if the decrypted data matches the original data
	assert.Equal(t, originalData, decrypted, "Decrypted data should match the original data")

}

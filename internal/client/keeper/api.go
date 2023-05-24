package keeper

import (
	"context"
	"github.com/itksb/go-secret-keeper/pkg/contract"
)

var _ contract.IKeeperRepository = &APIKeeper{}

// APIKeeper -  keeper server api wrapper
type APIKeeper struct {
	crypto contract.ISecretCrypto
	l      contract.IApplicationLogger
}

// NewAPIKeeper - constructor. Creates new keeper server api instance
func NewAPIKeeper(
	crypto contract.ISecretCrypto,
	l contract.IApplicationLogger,
) *APIKeeper {
	return &APIKeeper{
		crypto: crypto,
		l:      l,
	}
}

// SaveSecret - save or update user secret
func (a *APIKeeper) SaveSecret(
	ctx context.Context,
	userID contract.UserID,
	secret contract.IUserSecretItem,
) error {

	var err error

	secretDTO := secret.DTO()
	secretDTO.EncryptedData, err = a.crypto.Encrypt(secretDTO.EncryptedData)
	if err != nil {
		a.l.Errorf("failed to encrypt secret data for user %s: %s", userID, err)
		return err
	}

	secretDTO.EncryptedMeta, err = a.crypto.Encrypt(secretDTO.EncryptedMeta)
	if err != nil {
		a.l.Errorf("failed to encrypt secret meta for user %s: %s", userID, err)
		return err
	}

	panic("implement me - grpc call to server")

}

// GetAllSecrets - get all user secrets
func (a *APIKeeper) GetAllSecrets(
	ctx context.Context,
	userID contract.UserID,
) ([]contract.IUserSecretItem, error) {

	// retreives all secrets from the server and decrypts them using crypto ...

	//TODO implement me
	panic("implement me")
}

// Delete - delete user secret
func (a *APIKeeper) Delete(
	ctx context.Context,
	userID contract.UserID,
	id contract.SecretID,
) error {
	//TODO implement me
	panic("implement me")
}

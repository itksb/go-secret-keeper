package keeper

import (
	"context"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"time"
)

var _ IClientKeeper = &ClientKeeper{}

// ClientKeeper - client keeper
type ClientKeeper struct {
	repo              contract.IKeeperRepository
	l                 contract.IApplicationLogger
	crypto            ISecretCrypto
	getPrivateKeyFunc GetPrivateKeyFunc
}

// NewClientKeeper - constructor. Creates new client keeper
func NewClientKeeper(
	repo contract.IKeeperRepository,
	l contract.IApplicationLogger,
	crypto ISecretCrypto,
	getPrivateKeyFunc GetPrivateKeyFunc,
) *ClientKeeper {
	return &ClientKeeper{
		repo:              repo,
		l:                 l,
		crypto:            crypto,
		getPrivateKeyFunc: getPrivateKeyFunc,
	}
}

// SaveSecret - save secret for the user
func (c *ClientKeeper) SaveSecret(ctx context.Context, userID contract.UserID, secret contract.IUserSecretItem) error {
	// network request under the hood. So timeout is twice bigger than server keeper timeout
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err := c.repo.SaveSecret(ctx2, userID, secret)
	if err != nil {
		c.l.Errorf("failed to save secret for user %s: %s", userID, err)
		return err
	}
	return nil
}

// GetAllSecrets - get all secrets for user
func (c *ClientKeeper) GetAllSecrets(ctx context.Context, userID contract.UserID) ([]contract.IUserSecretItem, error) {
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	items, err := c.repo.GetAllSecrets(ctx2, userID)
	if err != nil {
		c.l.Errorf("failed to get all secrets for user %s: %s", userID, err)
		return nil, err
	}
	return items, nil
}

func (c *ClientKeeper) Delete(ctx context.Context, userID contract.UserID, id contract.SecretID) error {
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	err := c.repo.Delete(ctx2, userID, id)
	if err != nil {
		c.l.Errorf("failed to delete secret %s for user %s: %s", id, userID, err)
		return err
	}
	return nil
}

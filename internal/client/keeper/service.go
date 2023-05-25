package keeper

import (
	"context"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"time"
)

var _ IClientKeeper = &ClientKeeper{}

// ClientKeeper - client keeper
type ClientKeeper struct {
	api contract.IKeeperRepository
	l   contract.IApplicationLogger
}

// NewClientKeeper - constructor. Creates new client keeper
func NewClientKeeper(
	repo contract.IKeeperRepository,
	l contract.IApplicationLogger,
) *ClientKeeper {
	return &ClientKeeper{
		api: repo,
		l:   l,
	}
}

// SaveSecret - save secret for the user
func (c *ClientKeeper) SaveSecret(
	ctx context.Context,
	userID contract.UserID,
	secret contract.IUserSecretItem,
) (contract.IUserSecretItem, error) {
	// network request under the hood. So timeout is twice bigger than server keeper timeout
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	savedSecret, err := c.api.SaveSecret(ctx2, userID, secret)
	if err != nil {
		c.l.Errorf("failed to save secret for user %s: %s", userID, err)
		return nil, err
	}
	return savedSecret, nil
}

// GetAllSecrets - get all secrets for user
func (c *ClientKeeper) GetAllSecrets(
	ctx context.Context,
	userID contract.UserID,
) ([]contract.IUserSecretItem, error) {
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	items, err := c.api.GetAllSecrets(ctx2, userID)
	if err != nil {
		c.l.Errorf("failed to get all secrets for user %s: %s", userID, err)
		return nil, err
	}
	return items, nil
}

// Delete - delete secret for user
func (c *ClientKeeper) Delete(
	ctx context.Context,
	userID contract.UserID,
	id contract.SecretID,
) error {
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	err := c.api.Delete(ctx2, userID, id)
	if err != nil {
		c.l.Errorf("failed to delete secret %s for user %s: %s", id, userID, err)
		return err
	}
	return nil
}

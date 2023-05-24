package keeper

import (
	"context"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"time"
)

// IServerKeeper - server keeper interface
type IServerKeeper interface {
	contract.IKeeper
}

var _ IServerKeeper = &ServerKeeper{}

// ServerKeeper - server keeper
type ServerKeeper struct {
	repo contract.IKeeperRepository
	l    contract.IApplicationLogger
}

// NewServerKeeper - constructor. Creates new server keeper
func NewServerKeeper(
	repo contract.IKeeperRepository,
	l contract.IApplicationLogger,
) *ServerKeeper {
	return &ServerKeeper{
		repo: repo,
		l:    l,
	}
}

// SaveSecret - save secret for the user
func (s *ServerKeeper) SaveSecret(
	ctx context.Context,
	userID contract.UserID,
	secret contract.IUserSecretItem,
) error {

	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := s.repo.SaveSecret(ctx2, userID, secret)
	if err != nil {
		s.l.Errorf("failed to save secret for user %s: %s", userID, err)
		return err
	}
	return nil
}

// GetAllSecrets - get all secrets for user
func (s *ServerKeeper) GetAllSecrets(
	ctx context.Context,
	userID contract.UserID,
) ([]contract.IUserSecretItem, error) {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	items, err := s.repo.GetAllSecrets(ctx2, userID)
	if err != nil {
		s.l.Errorf("failed to get all secrets for user %s: %s", userID, err)
		return nil, err
	}
	return items, nil
}

// Delete - delete secret for the user
func (s *ServerKeeper) Delete(
	ctx context.Context,
	userID contract.UserID,
	id contract.SecretID,
) error {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := s.repo.Delete(ctx2, userID, id)
	if err != nil {
		s.l.Errorf("failed to delete secret %s for user %s: %s", id, userID, err)
		return err
	}
	return nil
}

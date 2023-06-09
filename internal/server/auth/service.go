package auth

import (
	"context"
	"github.com/itksb/go-secret-keeper/pkg/contract"
)

var _ contract.IAuthService = &AuthService{}

// AuthService - Auth service
type AuthService struct {
	tokenProvider contract.ITokenProvider
	repo          contract.IAuthRepository
	passHasher    contract.IPassHasher
}

// NewAuthService - create new auth service
func NewAuthService(
	tokenProvider contract.ITokenProvider,
	repo contract.IAuthRepository,
	passHasher contract.IPassHasher,
) *AuthService {
	return &AuthService{
		tokenProvider: tokenProvider,
		repo:          repo,
		passHasher:    passHasher,
	}
}

// SignUp - sign up new user
func (a *AuthService) SignUp(
	ctx context.Context,
	login,
	password string,
) (contract.IAccount, string, error) {
	passwordHash, err := a.passHasher.HashPassword([]byte(password))
	if err != nil {
		return nil, "", err
	}
	acc, err := a.repo.Create(ctx, login, passwordHash)
	if err != nil {
		return nil, "", err
	}
	token, err := a.tokenProvider.GenerateToken(ctx, acc.GetID())
	return acc, token, err
}

// SignIn - sign in existing user
func (a *AuthService) SignIn(
	ctx context.Context,
	login,
	password string,
) (contract.IAccount, string, error) {
	passwordHash, err := a.passHasher.HashPassword([]byte(password))
	if err != nil {
		return nil, "", err
	}

	acc, err := a.repo.Find(ctx, login, passwordHash)
	if err != nil {
		return nil, "", err
	}
	token, err := a.tokenProvider.GenerateToken(ctx, acc.GetID())

	return acc, token, err
}

package auth

import (
	"context"
	"github.com/itksb/go-secret-keeper/pkg/contract"
)

var _ IClientAuthService = &ClientAuthService{}

// ClientAuthService - client auth service
type ClientAuthService struct {
}

// NewClientAuthService - create new client auth service
func NewClientAuthService() *ClientAuthService {
	return &ClientAuthService{}
}

// SignUp - sign up
func (s *ClientAuthService) SignUp(
	ctx context.Context,
	login,
	password string,
) (
	acc contract.IAccount,
	token string,
	err error,
) {
	panic("implement me")
}

// SignIn - sign in
func (s *ClientAuthService) SignIn(
	ctx context.Context,
	login,
	password string,
) (
	acc contract.IAccount,
	token string,
	err error,
) {
	panic("implement me")
}

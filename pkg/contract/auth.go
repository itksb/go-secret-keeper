package contract

import (
	"context"
	"errors"
)

// IAccount - Account interface
type IAccount interface {
	GetID() string
	GetLogin() string
	SetPasswordHash(string)
}

// IAuthService - Auth service interface
type ITokenProvider interface {
	GenerateToken(ctx context.Context, userID string) (string, error)
	ValidateToken(ctx context.Context, token string) (userID string, err error)
}

// IAuthService - Auth service interface
type IAuthService interface {
	SignUp(ctx context.Context, login, password string) (IAccount, error)
	SignIn(ctx context.Context, login, password string) (IAccount, error)
}

var ErrDuplicateAccount = errors.New("duplicate account")

// IAuthRepository - Auth repository interface
type IAuthRepository interface {
	Create(ctx context.Context, login, password string) (IAccount, error)
	Find(ctx context.Context, login, password string) (IAccount, error)
}

type IPassHasher interface {
	HashPassword(password []byte) (string, error)
	ComparePassword(hashedPwd string, plainPwd []byte) (bool, error)
}

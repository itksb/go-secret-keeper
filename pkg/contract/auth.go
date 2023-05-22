package contract

import "context"

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

// IAuthRepository - Auth repository interface
type IAuthRepository interface {
	Create(ctx context.Context, login, password string) (IAccount, error)
	Find(ctx context.Context, login, password string) (IAccount, error)
}

package keeper

import (
	"context"
)

// use grpc

// id user_id contentType content

// account:
// id login password

type UserSecret []byte

type IUserSecret interface {
	GetType() string
}

type SecretID string

type IKeeper interface {
	SaveSecret(ctx context.Context, secret IUserSecret) error
	GetAllSecrets(ctx context.Context) ([]IUserSecret, error)
	Delete(ctx context.Context, id SecretID) error
}

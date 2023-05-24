package contract

import (
	"context"
)

type (
	// UserSecret - user secret
	UserSecretData []byte
	// SecretID - secret id
	SecretID string
	// UserID - user id
	UserID string
	// UserSecretMeta - user secret meta information such as data holder name, etc.
	UserSecretMeta string
)

type UserSecretType string

const (
	// UserSecretTypeBase - base secret type
	UserSecretTypeBase UserSecretType = "base"
	// UserSecretTypeTextData - text data secret type
	UserSecretTypeTextData UserSecretType = "textdata"
	// UserSecretLoginPasswd - login password secret type
	UserSecretLoginPasswd UserSecretType = "loginpasswd"
)

// IUserSecretItem - user secret interface
type IUserSecretItem interface {
	GetID() SecretID
	GetType() UserSecretType
	GetData() UserSecretData
	GetMeta() UserSecretMeta
}

// IKeeper - keeper interface
type IKeeper interface {
	// SaveSecret - save secret for the user
	SaveSecret(ctx context.Context, userID UserID, secret IUserSecretItem) error
	// GetAllSecrets - get all secrets for user
	GetAllSecrets(ctx context.Context, userID UserID) ([]IUserSecretItem, error)
	// Delete - delete secret for the user
	Delete(ctx context.Context, userID UserID, id SecretID) error
}

// IKeeperRepository - keeper repository interface
type IKeeperRepository interface {
	// SaveSecret - save secret for the user
	SaveSecret(ctx context.Context, userID UserID, secret IUserSecretItem) error
	// GetAllSecrets - get all secrets for user
	GetAllSecrets(ctx context.Context, userID UserID) ([]IUserSecretItem, error)
	// Delete - delete secret for the user
	Delete(ctx context.Context, userID UserID, id SecretID) error
}
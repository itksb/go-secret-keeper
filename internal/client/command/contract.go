package command

import "github.com/itksb/go-secret-keeper/pkg/contract"

// ICommand - command interface
type ICommand interface {
	Execute() error
}

// SecretItemProcessorFunc -  secret item processor function callback
type SecretItemProcessorFunc func(secret contract.IUserSecretItem) error

// SecretsProcessorFunc - list of secretS processor function callback
type SecretsProcessorFunc func(secrets []contract.IUserSecretItem) error

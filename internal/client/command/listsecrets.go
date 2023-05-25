package command

import "github.com/itksb/go-secret-keeper/pkg/contract"

var _ ICommand = &ListSecretsCommand{}

// ListSecretsCommand - list secrets command
type ListSecretsCommand struct {
	l contract.IApplicationLogger
}

// NewListSecretsCommand - create new list secrets command
func NewListSecretsCommand(
	l contract.IApplicationLogger,
) *ListSecretsCommand {
	return &ListSecretsCommand{
		l: l,
	}
}

// Execute - execute command
func (c *ListSecretsCommand) Execute() error {
	return nil
}

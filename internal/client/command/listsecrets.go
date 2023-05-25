package command

import (
	"context"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"time"
)

var _ ICommand = &ListSecretsCommand{}

// ListSecretsCommand - list secrets command
type ListSecretsCommand struct {
	l           contract.IApplicationLogger
	keeper      contract.IKeeper
	userID      string
	processFunc SecretsProcessorFunc
}

// NewListSecretsCommand - create new list secrets command
func NewListSecretsCommand(
	l contract.IApplicationLogger,
	keeper contract.IKeeper,
	userID string,
	processFunc SecretsProcessorFunc,
) *ListSecretsCommand {
	return &ListSecretsCommand{
		l:           l,
		keeper:      keeper,
		userID:      userID,
		processFunc: processFunc,
	}
}

// Execute - execute command
func (c *ListSecretsCommand) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	secrets, err := c.keeper.GetAllSecrets(ctx, contract.UserID(c.userID))
	if err != nil {
		c.l.Errorf("failed to get all secrets for user: %s. Err: %s", c.userID, err.Error())
		return err
	}

	if c.processFunc != nil {
		err = c.processFunc(secrets)
		if err != nil {
			c.l.Errorf("failed to process secrets for user: %s. Err: %s", c.userID, err.Error())
			return err
		}
	}

	return nil
}

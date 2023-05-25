package command

import (
	"context"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"time"
)

var _ ICommand = &DeleteSecretCommand{}

// DeleteSecretCommand - delete secret command
type DeleteSecretCommand struct {
	l        contract.IApplicationLogger
	keeper   contract.IKeeper
	userID   string
	secretID string
}

// NewDeleteSecretCommand - create new delete secret command
func NewDeleteSecretCommand(
	l contract.IApplicationLogger,
	keeper contract.IKeeper,
	userID string,
	secretID string,
) *DeleteSecretCommand {
	return &DeleteSecretCommand{
		l:        l,
		keeper:   keeper,
		userID:   userID,
		secretID: secretID,
	}
}

// Execute - execute command
func (c *DeleteSecretCommand) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.keeper.Delete(
		ctx,
		contract.UserID(c.userID),
		contract.SecretID(c.secretID),
	)
	if err != nil {
		c.l.Errorf("failed to delete secret for user: %s. Err: %s", c.userID, err.Error())
		return err
	}

	return nil
}

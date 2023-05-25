package command

import (
	"context"
	"errors"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"time"
)

var _ ICommand = &SaveSecretCommand{}

// SaveSecretCommand - save secret command
type SaveSecretCommand struct {
	l                 contract.IApplicationLogger
	keeper            contract.IKeeper
	userID            string
	secret            contract.IUserSecretItem
	secretProcessFunc SecretItemProcessorFunc
}

// NewSaveSecretCommand - execute command
func NewSaveSecretCommand(
	l contract.IApplicationLogger,
	keeper contract.IKeeper,
	userID string,
	secret contract.IUserSecretItem,
) *SaveSecretCommand {
	return &SaveSecretCommand{
		l:      l,
		keeper: keeper,
		userID: userID,
		secret: secret,
	}
}

// Execute - execute command
func (c *SaveSecretCommand) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	secretSaved, err := c.keeper.SaveSecret(ctx, contract.UserID(c.userID), c.secret)
	if err != nil {
		c.l.Errorf("failed to save secret for user: %s. Err: %s", c.userID, err.Error())
		return err
	}

	if c.secretProcessFunc != nil {
		err = c.secretProcessFunc(secretSaved)
		if err != nil {
			c.l.Errorf("failed to process secret callback for user: %s. Err: %s", c.userID, err.Error())
			return errors.New("failed to process secret callback")
		}
	}

	return nil
}

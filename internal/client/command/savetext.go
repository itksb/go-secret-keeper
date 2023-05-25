package command

import (
	"context"
	"errors"
	secret2 "github.com/itksb/go-secret-keeper/internal/client/keeper/secret"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"github.com/itksb/go-secret-keeper/pkg/keeper/secret"
	"time"
)

var _ ICommand = &SaveTextCommand{}

// SaveSecretCommand - save secret command
type SaveTextCommand struct {
	l      contract.IApplicationLogger
	keeper contract.IKeeper
	userID string

	secretID string
	text     string
	meta     string

	secretProcessFunc SecretItemProcessorFunc
}

// NewSaveTextCommand - execute command
func NewSaveTextCommand(
	l contract.IApplicationLogger,
	keeper contract.IKeeper,
	userID string,
	secretID string,
	text string,
	meta string,
	secretItemProcessorFunc SecretItemProcessorFunc,
) *SaveTextCommand {
	return &SaveTextCommand{
		l:        l,
		keeper:   keeper,
		userID:   userID,
		secretID: secretID,

		meta:              meta,
		secretProcessFunc: secretItemProcessorFunc,
	}
}

// SaveTextCommandFactory - create new save secret command factory
// using currying technic for dependency injection
func SaveTextCommandFactory(
	l contract.IApplicationLogger,
	keeper contract.IKeeper,
) func(
	userID string,
	secretID string,
	text string,
	meta string,
	secretItemProcessorFunc SecretItemProcessorFunc,
) *SaveTextCommand {
	return func(
		userID string,
		secretID string,
		text string,
		meta string,
		secretItemProcessorFunc SecretItemProcessorFunc,
	) *SaveTextCommand {
		return NewSaveTextCommand(
			l,
			keeper,
			userID,
			secretID,
			text,
			meta,
			secretItemProcessorFunc,
		)
	}
}

// Execute - execute command
func (c *SaveTextCommand) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	entity, err := secret.CreateSecretItem(
		contract.SecretID(c.secretID),
		contract.UserSecretData{},
		contract.UserSecretMeta(c.meta),
		contract.UserSecretTypeTextData,
	)

	if err != nil {
		c.l.Errorf("failed to create secret for user: %s. Err: %s", c.userID, err.Error())
		return err
	}

	secretItem, ok := entity.(*secret.TextDataSecretItem)
	if !ok {
		c.l.Errorf("failed to cast secret for user: %s. Err: %s", c.userID, err.Error())
		return err
	}

	packer := secret2.TextDataSecretItemPacker{
		Entity: *secretItem,
	}

	err = packer.Write(c.text)
	if err != nil {
		c.l.Errorf("failed to pack secret for user: %s. Err: %s", c.userID, err.Error())
		return err
	}

	secretItem = &packer.Entity

	secretSaved, err := c.keeper.SaveSecret(ctx, contract.UserID(c.userID), secretItem)
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

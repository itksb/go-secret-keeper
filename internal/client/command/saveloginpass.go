package command

import (
	"context"
	"errors"
	secret2 "github.com/itksb/go-secret-keeper/internal/client/keeper/secret"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"github.com/itksb/go-secret-keeper/pkg/keeper/secret"
	"time"
)

var _ ICommand = &SaveLoginPassCommand{}

// SaveSecretCommand - save secret command
type SaveLoginPassCommand struct {
	l      contract.IApplicationLogger
	keeper contract.IKeeper
	userID string

	secretID string
	login    string
	password string
	meta     string

	secretProcessFunc SecretItemProcessorFunc
}

// NewSaveSecretCommand - execute command
func NewSaveSecretCommand(
	l contract.IApplicationLogger,
	keeper contract.IKeeper,
	userID string,
	secretID string,
	login string,
	password string,
	meta string,
	secretItemProcessorFunc SecretItemProcessorFunc,
) *SaveLoginPassCommand {
	return &SaveLoginPassCommand{
		l:                 l,
		keeper:            keeper,
		userID:            userID,
		secretID:          secretID,
		login:             login,
		password:          password,
		meta:              meta,
		secretProcessFunc: secretItemProcessorFunc,
	}
}

// SaveSecretCommandFactory - create new save secret command factory
// using currying technic for dependency injection
func SaveSecretCommandFactory(
	l contract.IApplicationLogger,
	keeper contract.IKeeper,
) func(
	userID string,
	secretID string,
	login string,
	password string,
	meta string,
	secretItemProcessorFunc SecretItemProcessorFunc,
) *SaveLoginPassCommand {
	return func(
		userID string,
		secretID string,
		login string,
		password string,
		meta string,
		secretItemProcessorFunc SecretItemProcessorFunc,
	) *SaveLoginPassCommand {
		return NewSaveSecretCommand(
			l,
			keeper,
			userID,
			secretID,
			login,
			password,
			meta,
			secretItemProcessorFunc,
		)
	}
}

// Execute - execute command
func (c *SaveLoginPassCommand) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	entity, err := secret.CreateSecretItem(
		contract.SecretID(c.secretID),
		contract.UserSecretData{},
		contract.UserSecretMeta(c.meta),
		contract.UserSecretLoginPasswd,
	)

	if err != nil {
		c.l.Errorf("failed to create secret for user: %s. Err: %s", c.userID, err.Error())
		return err
	}

	secretItem, ok := entity.(*secret.LoginPasswdSecretItem)
	if !ok {
		c.l.Errorf("failed to cast secret for user: %s. Err: %s", c.userID, err.Error())
		return err
	}

	packer := secret2.LoginPasswdSecretItemPacker{
		Entity: *secretItem,
	}

	err = packer.Write(c.login, c.password)
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

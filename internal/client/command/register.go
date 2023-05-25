package command

import (
	"context"
	"errors"
	"github.com/itksb/go-secret-keeper/internal/client/session"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"time"
)

var _ ICommand = &RegisterCommand{}

// RegisterCommand - register command
type RegisterCommand struct {
	auth     contract.IAuthService
	session  session.ISession
	l        contract.IApplicationLogger
	login    string
	password string
}

// NewRegisterCommand - execute command
func NewRegisterCommand(
	auth contract.IAuthService,
	l contract.IApplicationLogger,
	session session.ISession,
	login string,
	password string,
) *RegisterCommand {
	return &RegisterCommand{
		auth:     auth,
		session:  session,
		l:        l,
		login:    login,
		password: password,
	}
}

// RegisterCmdAbstractFabric - abstract fabric for register command
// using currying technic for dependency injection
func RegisterCmdAbstractFabric(
	auth contract.IAuthService,
	l contract.IApplicationLogger,
) func(
	session session.ISession,
	login string,
	password string,
) *RegisterCommand {
	return func(
		session session.ISession,
		login string,
		password string,
	) *RegisterCommand {
		return NewRegisterCommand(auth, l, session, login, password)
	}
}

// Execute - execute command
func (c *RegisterCommand) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	account, token, err := c.auth.SignUp(ctx, c.login, c.password)
	if err != nil {
		c.l.Errorf("failed to sign up for login: %s. Err: %s", c.login, err.Error())
		return err
	}

	err = c.session.SetValue(session.AccountKey, account)
	if err != nil {
		c.l.Errorf("failed to set account to session", err)
		return errors.New("failed to set account to session")
	}

	err = c.session.SetValue(session.TokenKey, token)
	if err != nil {
		c.l.Errorf("failed to set token to session, err: %s", err.Error())
		return errors.New("failed to set token to session")
	}

	return nil
}

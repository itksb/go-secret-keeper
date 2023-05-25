package command

import (
	"context"
	"errors"
	"github.com/itksb/go-secret-keeper/internal/client/session"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"time"
)

var _ ICommand = &LoginCommand{}

// LoginCommand - login command
type LoginCommand struct {
	auth     contract.IAuthService
	session  session.ISession
	l        contract.IApplicationLogger
	login    string
	password string
}

// NewLoginCommand - create new login command
func NewLoginCommand(
	auth contract.IAuthService,
	l contract.IApplicationLogger,
	session session.ISession,
	login string,
	password string,
) *LoginCommand {
	return &LoginCommand{
		auth:     auth,
		session:  session,
		l:        l,
		login:    login,
		password: password,
	}
}

// Execute - execute command
func (c *LoginCommand) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	account, token, err := c.auth.SignIn(ctx, c.login, c.password)
	if err != nil {
		c.l.Errorf("failed to sign for login: %s. Err: %s", c.login, err.Error())
		return errors.New("invalid login or password")
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

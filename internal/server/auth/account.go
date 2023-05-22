package auth

import "github.com/itksb/go-secret-keeper/pkg/contract"

var _ contract.IAccount = &Account{}

// Account - User
type Account struct {
	ID           string
	Login        string
	PasswordHash string
}

// Account - GetID
func (a *Account) GetID() string {
	return a.ID
}

// Account - GetLogin
func (a *Account) GetLogin() string {
	return a.Login
}

// Account - GetPasswordHash
func (a *Account) SetPasswordHash(h string) {
	a.PasswordHash = h
}

package auth

import "github.com/itksb/go-secret-keeper/pkg/contract"

var _ contract.IAccount = &Account{}

// Account - User
type Account struct {
	ID           string
	Login        string
	PasswordHash string
}

// GetID get account id
func (a *Account) GetID() string {
	return a.ID
}

// GetLogin get account login
func (a *Account) GetLogin() string {
	return a.Login
}

// SetPasswordHash set password hash
func (a *Account) SetPasswordHash(h string) {
	a.PasswordHash = h
}

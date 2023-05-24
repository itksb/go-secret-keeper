package secret

import "github.com/itksb/go-secret-keeper/pkg/contract"

// LoginPasswdSecretItem - login passwd secret item
type LoginPasswdSecretItem struct {
	BaseSecretItem
}

// GetType - get secret type
func (i *LoginPasswdSecretItem) GetType() contract.UserSecretType {
	return contract.UserSecretLoginPasswd
}

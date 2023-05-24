package secret

import "github.com/itksb/go-secret-keeper/pkg/contract"

// TextDataSecretItem - text data secret item
type TextDataSecretItem struct {
	BaseSecretItem
}

// GetType - get secret type
func (i *TextDataSecretItem) GetType() contract.UserSecretType {
	return contract.UserSecretTypeTextData
}

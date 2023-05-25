package secret

import "github.com/itksb/go-secret-keeper/pkg/contract"

// CreateSecretItem - create secret item fabric method
// returns error ErrUnknownSecretType if unknown secret type
func CreateSecretItem(
	id contract.SecretID,
	data contract.UserSecretData,
	meta contract.UserSecretMeta,
	secretType contract.UserSecretType,
) (contract.IUserSecretItem, error) {

	base := BaseSecretItem{
		id:   id,
		data: data,
		meta: meta,
	}

	switch secretType {
	case contract.UserSecretTypeTextData:
		return &TextDataSecretItem{
			base,
		}, nil
	case contract.UserSecretLoginPasswd:
		return &LoginPasswdSecretItem{
			base,
		}, nil

	default:
		return nil, ErrUnknownSecretType
	}
}

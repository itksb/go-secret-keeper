package secret

import (
	"github.com/itksb/go-secret-keeper/pkg/contract"
)

var _ contract.IUserSecretItem = &BaseSecretItem{}

// BaseSecretItem - base secret item
type BaseSecretItem struct {
	id   contract.SecretID
	data contract.UserSecretData
	meta contract.UserSecretMeta
}

// NewBaseSecretItem - create new base secret item
func (b *BaseSecretItem) GetMeta() contract.UserSecretMeta {
	return b.meta
}

// GetID - get secret id
func (b *BaseSecretItem) GetID() contract.SecretID {
	return b.id
}

// GetType - get secret type
func (*BaseSecretItem) GetType() contract.UserSecretType {
	return contract.UserSecretTypeBase
}

// GetData - get secret data
func (b *BaseSecretItem) GetData() contract.UserSecretData {
	return b.data
}

// DTO - convert to DTO
func (b *BaseSecretItem) DTO() contract.UserSecretItemDTO {
	return contract.UserSecretItemDTO{
		ID:            b.id,
		Type:          b.GetType(),
		EncryptedData: b.data,
		EncryptedMeta: []byte(b.meta),
	}
}

// FromDTO - create from DTO
func (b *BaseSecretItem) FromDTO(dto contract.UserSecretItemDTO) (contract.IUserSecretItem, error) {
	return CreateSecretItem(
		dto.ID,
		dto.EncryptedData,
		contract.UserSecretMeta(dto.EncryptedMeta),
		dto.Type,
	)
}

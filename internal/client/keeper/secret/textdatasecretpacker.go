package secret

import (
	"errors"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"github.com/itksb/go-secret-keeper/pkg/keeper/secret"
)

// TextDataSecretItemPacker - text data secret item packer
type TextDataSecretItemPacker struct {
	Entity secret.TextDataSecretItem
}

// Read - read secret item
func (p *TextDataSecretItemPacker) Read() (
	text string,
	err error,
) {
	return string(p.Entity.GetData()), nil
}

// Write - write secret item
func (p *TextDataSecretItemPacker) Write(text string) error {
	var entity contract.IUserSecretItem
	var err error
	entity, err = p.Entity.FromDTO(
		contract.UserSecretItemDTO{
			ID:            p.Entity.GetID(),
			Type:          p.Entity.GetType(),
			EncryptedData: []byte(text),
			EncryptedMeta: []byte(p.Entity.GetMeta()),
		},
	)

	if err != nil {
		return err
	}

	item, ok := entity.(*secret.TextDataSecretItem)
	if ok {
		p.Entity = *item
		return nil
	} else {
		return errors.New("invalid entity type")
	}

	return nil
}

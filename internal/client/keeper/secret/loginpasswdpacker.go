package secret

import (
	"encoding/json"
	"errors"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"github.com/itksb/go-secret-keeper/pkg/keeper/secret"
)

// LoginPasswdSecretItemPacker - login passwd secret item packer
type LoginPasswdSecretItemPacker struct {
	Entity secret.LoginPasswdSecretItem
}

type dataStruct struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (p *LoginPasswdSecretItemPacker) Read() (
	login string,
	passwd string,
	err error,
) {
	data := p.Entity.GetData()
	d := &dataStruct{}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return "", "", err
	}
	return d.Login, d.Password, nil
}

// Write - write secret item
func (p *LoginPasswdSecretItemPacker) Write(
	login string,
	passwd string,
) error {
	d := &dataStruct{
		Login:    login,
		Password: passwd,
	}
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}

	var entity contract.IUserSecretItem
	entity, err = p.Entity.FromDTO(
		contract.UserSecretItemDTO{
			ID:            p.Entity.GetID(),
			Type:          p.Entity.GetType(),
			EncryptedData: []byte(data),
			EncryptedMeta: []byte(p.Entity.GetMeta()),
		},
	)

	if err != nil {
		return err
	}

	item, ok := entity.(*secret.LoginPasswdSecretItem)
	if ok {
		p.Entity = *item
		return nil
	} else {
		return errors.New("invalid entity type")
	}

}

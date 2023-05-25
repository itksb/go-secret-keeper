package keeper

import (
	"context"
	"database/sql"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	secret2 "github.com/itksb/go-secret-keeper/pkg/keeper/secret"
	"github.com/jmoiron/sqlx"
)

// SecretsRepo - secrets repository
type SecretsRepo struct {
	db *sqlx.DB
	l  contract.IApplicationLogger
}

// NewSecretsRepo - create new secrets repository
func NewSecretsRepo(
	db *sqlx.DB,
	l contract.IApplicationLogger,
) *SecretsRepo {
	return &SecretsRepo{
		db: db,
		l:  l,
	}
}

// SaveSecret - save secret for the user
func (r *SecretsRepo) SaveSecret(
	ctx context.Context,
	userID contract.UserID,
	secret contract.IUserSecretItem,
) (contract.IUserSecretItem, error) {
	var query string
	var row *sql.Row

	if secret.GetID() == "" {
		query = `INSERT INTO secrets(user_id, type, data, meta) VALUES ($1, $2, $3, $4) RETURNING id`
		row = r.db.QueryRowContext(
			ctx,
			query,
			userID,
			secret.GetType(),
			secret.GetData(),
			secret.GetMeta(),
		)
	} else {
		query = `UPDATE secrets SET type = $1, data = $2, meta = $3 WHERE id = $4 RETURNING id`
		row = r.db.QueryRowContext(
			ctx,
			query,
			secret.GetType(),
			secret.GetData(),
			secret.GetMeta(),
			secret.GetID(),
		)
	}

	var returningID string
	err := row.Scan(&returningID)
	if err != nil {
		r.l.Errorf("%s", err.Error())
		return secret, err
	}
	dto := secret.DTO() // .ID = returningID
	dto.ID = contract.SecretID(returningID)
	secret, err = secret.FromDTO(dto)
	if err != nil {
		return secret, err
	}

	return secret, nil
}

// GetAllSecrets - get all secrets for user
func (r *SecretsRepo) GetAllSecrets(
	ctx context.Context,
	userID contract.UserID,
) ([]contract.IUserSecretItem, error) {
	query := `SELECT id, type, data, meta FROM secrets WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.l.Errorf("%s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []contract.IUserSecretItem
	for rows.Next() {
		var dto contract.UserSecretItemDTO
		err = rows.Scan(&dto.ID, &dto.Type, &dto.EncryptedData, &dto.EncryptedMeta)
		if err != nil {
			r.l.Errorf("%s", err.Error())
			return nil, err
		}
		var err2 error

		secret, err2 := secret2.CreateSecretItem(
			dto.ID,
			dto.EncryptedData,
			contract.UserSecretMeta(dto.EncryptedMeta),
			dto.Type,
		)
		if err2 != nil {
			r.l.Errorf("%s", err2.Error())
			return nil, err2
		}
		result = append(result, secret)
	}

	return result, nil
}

// Delete - delete secret for the user
func (r *SecretsRepo) Delete(
	ctx context.Context,
	userID contract.UserID,
	id contract.SecretID,
) error {
	query := `DELETE FROM secrets WHERE user_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, id)
	if err != nil {
		r.l.Errorf("%s", err.Error())
		return err
	}
	return nil
}

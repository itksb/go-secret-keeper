package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"github.com/jmoiron/sqlx"
)

var _ contract.IAuthRepository = &AuthRepository{}

// AuthRepository - Auth repository
type AuthRepository struct {
	db *sqlx.DB
	l  contract.IApplicationLogger
}

// NewAuthRepository - create new auth repository
func NewAuthRepository(
	db *sqlx.DB,
	l contract.IApplicationLogger,
) *AuthRepository {
	return &AuthRepository{
		db: db,
		l:  l,
	}
}

// Create - create new user
func (a *AuthRepository) Create(
	ctx context.Context,
	login,
	password string,
) (contract.IAccount, error) {
	query := `INSERT INTO auth_accounts( login, password) VALUES ($1, $2) RETURNING id`
	row := a.db.QueryRowContext(ctx, query, login, password)
	result := &Account{
		ID:    "",
		Login: login,
	}

	var returningID string
	err := row.Scan(&returningID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//query does not return id, so duplicate conflict, need to retrieve id from db
			row2 := a.db.QueryRowContext(ctx, `SELECT id FROM auth_accounts WHERE login = $1`, login)
			err = row2.Scan(&returningID)
			if err != nil {
				a.l.Errorf("%s", err.Error())
				return result, err
			}
			result.ID = returningID
			return result, fmt.Errorf("%w", contract.ErrDuplicateAccount)
		}
	}
	result.ID = returningID

	return result, nil
}

// Find - find user by login and password
func (a *AuthRepository) Find(
	ctx context.Context,
	login,
	password string,
) (contract.IAccount, error) {
	//query does not return id, so duplicate conflict, need to retrieve id from db
	row := a.db.QueryRowContext(
		ctx,
		`SELECT id FROM auth_accounts WHERE login = $1 and password = $2 `,
		login,
		password,
	)

	var returningID string

	result := &Account{
		ID:    "",
		Login: login,
	}

	err := row.Scan(&returningID)
	if err != nil {
		a.l.Errorf("%s", err.Error())
		return result, err
	}
	result.ID = returningID

	return result, nil
}

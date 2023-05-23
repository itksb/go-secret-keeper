package db

import (
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const dbDriverName = "postgres"

func NewPostgresDbPool(
	dsn string,
	l contract.IApplicationLogger,
) (*sqlx.DB, error) {
	db, err := sqlx.Connect(dbDriverName, dsn)
	if err != nil {
		l.Errorf("connection to the db: %s", err.Error())
		return nil, err
	}

	return db, nil
}

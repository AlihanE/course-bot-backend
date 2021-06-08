package db

import (
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/stdlib"
)

func Connect(logger hclog.Logger, driver, connString string) (*sqlx.DB, error) {
	var db *sqlx.DB

	db, err := sqlx.Open(driver, connString)
	if err != nil {
		logger.Error("failed to open connection", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error("failed to check connection", err)
		return nil, err
	}

	return db, nil
}

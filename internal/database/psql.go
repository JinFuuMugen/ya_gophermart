package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDatabase(config string) error {
	var err error
	DB, err = sql.Open("pgx", config)
	if err != nil {
		return fmt.Errorf("cannot connect to database: %w", err)
	}

	ctx := context.Background()

	_, err = DB.ExecContext(ctx, "SELECT 1") //TODO: change init query

	if err != nil {
		return fmt.Errorf("cannot create table: %w", err)
	}

	return nil
}

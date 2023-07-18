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

	_, err = DB.ExecContext(ctx, "") //TODO: change init query

	if err != nil {
		return fmt.Errorf("cannot create table: %w", err)
	}

	return nil
}

func StoreOrder(orderNum string, user string) error {

	ctx := context.Background()

	errCh := make(chan error, 1)

	go func() {

		_, err := DB.ExecContext(ctx, "INSERT INTO orders (order, user) VALUES ($1, $2);", orderNum, user) //TODO: might need changes
		if err != nil {
			errCh <- fmt.Errorf("error storing order in the database: %w", err)
			return
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

//func GetOrders(user string) error {
//
//	DB.ExecContext(ctx, "SELECT * FROM orders WHERE user = $uname;", user) //TODO: might need changes
//
//}

func StoreUser(login string, password string) error {

	ctx := context.Background()

	errCh := make(chan error, 1)

	go func() {

		_, err := DB.ExecContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2);", login, password) //TODO: might need changes
		if err != nil {
			errCh <- fmt.Errorf("error storing user in the database: %w", err)
			return
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

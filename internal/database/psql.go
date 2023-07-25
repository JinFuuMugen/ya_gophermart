package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/models"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

var DB *sql.DB

func InitDatabase(config string) error {
	var err error
	DB, err = sql.Open("pgx", config)
	if err != nil {
		return fmt.Errorf("cannot connect to database: %w", err)
	}

	ctx := context.Background()

	_, err = DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			login TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			PRIMARY KEY (login)
		)
	`)
	if err != nil {
		return fmt.Errorf("cannot create users table: %w", err)
	}

	_, err = DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			number numeric UNIQUE,
			login TEXT,
			uploaded TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("cannot create orders table: %w", err)
	}

	return nil
}

func CheckOrder(orderNum string, user string) (int, error) {
	errCh := make(chan error, 1)
	codeCh := make(chan int, 1)

	go func() {

		rows, err := DB.Query("SELECT login FROM orders WHERE number = $1", orderNum)
		if err != nil {
			errCh <- fmt.Errorf("error checking order: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var prevLogin string

			err = rows.Scan(&prevLogin)
			if err != nil {
				errCh <- fmt.Errorf("error scanning rows: %w", err)
			}
			if prevLogin == user {
				codeCh <- 200
			} else if prevLogin != user && prevLogin != "" {
				codeCh <- 409
			}
			codeCh <- 202
		}

		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return 0, err
	case code := <-codeCh:
		return code, nil
	}
}

func StoreOrder(orderNum string, user string) error {

	ctx := context.Background()

	errCh := make(chan error, 1)

	go func() {
		curentTime := time.Now().Format(time.RFC3339)
		t, err := time.Parse(time.RFC3339, curentTime)
		if err != nil {
			errCh <- fmt.Errorf("error parsing datetime: %w", err)
			return
		}
		_, err = DB.ExecContext(ctx, "INSERT INTO orders (number, login, uploaded) VALUES ($1, $2, $3);", orderNum, user, t) //TODO: might need changes
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

func UserAuth(authLogin string, authPass string) (bool, error) {

	errCh := make(chan error, 1)
	boolCh := make(chan bool, 1)

	go func() {

		rows, err := DB.Query("SELECT * FROM users WHERE login = $1", authLogin)
		if err != nil {
			errCh <- fmt.Errorf("error checking user: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var login string
			var pass string

			err = rows.Scan(&login, &pass)
			if err != nil {
				errCh <- fmt.Errorf("error scanning rows: %w", err)
			}
			if login == authLogin && pass == authPass {
				boolCh <- true
				return
			}
			boolCh <- false
		}

		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return false, err
	case flag := <-boolCh:
		return flag, nil
	}
}

func CheckLoginTaken(user string) (bool, error) {

	errCh := make(chan error, 1)
	boolCh := make(chan bool, 1)

	go func() {

		rows, err := DB.Query("SELECT COUNT(*) FROM users WHERE login = $1", user)
		if err != nil {
			errCh <- fmt.Errorf("error checking user: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var cnt int

			err = rows.Scan(&cnt)
			if err != nil {
				errCh <- fmt.Errorf("error scanning rows: %w", err)
			}
			if cnt == 1 {
				boolCh <- true
				return
			}
			boolCh <- false
		}

		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return false, err
	case flag := <-boolCh:
		return flag, nil
	}
}

func GetOrdersDB(user string) ([]models.Order, error) {
	errCh := make(chan error, 1)
	ordCh := make(chan []models.Order, 1)

	go func() {

		orders := make([]models.Order, 0)

		rows, err := DB.Query("SELECT number, uploaded FROM orders WHERE login = $1 ORDER BY uploaded", user) //TODO: might need changes
		if err != nil {
			errCh <- fmt.Errorf("error checking user: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var number int
			var date time.Time

			err := rows.Scan(&number, &date)
			if err != nil {
				errCh <- fmt.Errorf("error scanning rows: %w", err)
			}
			orders = append(orders, models.Order{Number: number, Dateadd: date.Format(time.RFC3339)})
		}
		ordCh <- orders
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return nil, err
	case ord := <-ordCh:
		return ord, nil
	}
}

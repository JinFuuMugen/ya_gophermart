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
	rows, err := DB.Query("SELECT login FROM orders WHERE number = $1", orderNum)
	if err != nil {
		return 0, fmt.Errorf("error checking order: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var prevLogin string

		err = rows.Scan(&prevLogin)
		if err != nil {
			return 0, fmt.Errorf("error scanning rows: %w", err)
		}
		if prevLogin == user {
			return 200, nil
		} else if prevLogin != user && prevLogin != "" {
			return 409, nil
		}
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error scanning rows: %w", err)
	}

	return 202, nil
}

func StoreOrder(orderNum string, user string) error {
	ctx := context.Background()

	currentTime := time.Now().Format(time.RFC3339)
	t, err := time.Parse(time.RFC3339, currentTime)
	if err != nil {
		return fmt.Errorf("error parsing datetime: %w", err)
	}

	_, err = DB.ExecContext(ctx, "INSERT INTO orders (number, login, uploaded) VALUES ($1, $2, $3);", orderNum, user, t)
	if err != nil {
		return fmt.Errorf("error storing order in the database: %w", err)
	}

	return nil
}

func StoreUser(login string, password string) error {
	ctx := context.Background()

	_, err := DB.ExecContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2);", login, password)
	if err != nil {
		return fmt.Errorf("error storing user in the database: %w", err)
	}

	return nil
}

func UserAuth(authLogin string, authPass string) (bool, error) {
	rows, err := DB.Query("SELECT * FROM users WHERE login = $1", authLogin)
	if err != nil {
		return false, fmt.Errorf("error checking user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var login string
		var pass string

		err = rows.Scan(&login, &pass)
		if err != nil {
			return false, fmt.Errorf("error scanning rows: %w", err)
		}
		if login == authLogin && pass == authPass {
			return true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("error scanning rows: %w", err)
	}

	return false, nil
}

func CheckLoginTaken(user string) (bool, error) {
	rows, err := DB.Query("SELECT COUNT(*) FROM users WHERE login = $1", user)
	if err != nil {
		return false, fmt.Errorf("error checking user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cnt int
		err = rows.Scan(&cnt)
		if err != nil {
			return false, fmt.Errorf("error scanning rows: %w", err)
		}
		if cnt == 1 {
			return true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("error scanning rows: %w", err)
	}

	return false, nil
}

func GetOrdersDB(user string) ([]models.Order, error) {
	orders := make([]models.Order, 0)

	rows, err := DB.Query("SELECT number, uploaded FROM orders WHERE login = $1 ORDER BY uploaded", user)
	if err != nil {
		return nil, fmt.Errorf("error checking user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var number string
		var date time.Time

		err := rows.Scan(&number, &date)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}
		orders = append(orders, models.Order{Number: number, Dateadd: date.Format(time.RFC3339)})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning rows: %w", err)
	}

	return orders, nil
}

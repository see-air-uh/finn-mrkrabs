package data

import (
	"context"
	"database/sql"
	"errors"
	"log"

	// "errors"

	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Transaction: Transaction{},
	}
}

type Models struct {
	Transaction Transaction
}

type Transaction struct {
	TransactionID     int     `json:"-"`
	UserID            string  `json:"id"`
	TransactionAmount float32 `json:"transactionAmount"`
}

func (t *Transaction) GetUserBalance(email string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select SUM(TransactionAmount) From mrkrabs.Transactions
	where Username = $1
	group by Username`

	var totalBalance float32

	row := db.QueryRowContext(ctx, query, email)
	err := row.Scan(&totalBalance)

	if err != nil {
		return 0, err
	}

	return totalBalance, nil
}

func (t *Transaction) UpdateBalance(username string, transactionAmount float32) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `insert into mrkrabs.Transactions (Username, TransactionAmount) values
	($1,$2)`

	balance, err := t.GetUserBalance(username)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, err
	}

	log.Println(balance + transactionAmount)
	if (balance + transactionAmount) < 0 {
		return 0, errors.New("error. can not decrement balance below zero")
	}

	_, err = db.ExecContext(ctx, query, username, transactionAmount)
	if err != nil {
		return 0, err
	}

	return balance + transactionAmount, nil
}

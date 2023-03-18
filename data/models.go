package data

import (
	"context"
	"database/sql"
	"errors"

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
	Transaction      Transaction
	RecurringPayment RecurringPayment
	PaymentHistory   PaymentHistory
}

type Transaction struct {
	TransactionID          int     `json:"-"`
	UserID                 string  `json:"id"`
	TransactionAmount      float32 `json:"transactionAmount"`
	TransactionName        string  `json:"transactionName"`
	TransactionDescription string  `json:"transactionDescription"`
}

type RecurringPayment struct {
	PaymentID          int     `json:"-"`
	UserName           string  `json:"username"`
	PaymentAmount      float32 `json:"amount"`
	PaymentName        string  `json:"paymentName"`
	PaymentDescription string  `json:"paymentDescription"`
	PaymentDate        string  `json:"paymentDate"`
}

type PaymentHistory struct {
	PaymentHistoryID     int    `json:"-"`
	PaymentID            int    `json:"paymentID"`
	PaymentHistoryDate   string `json:"paymentHistoryDate"`
	PaymentHistoryStatus bool   `json:"paymentHistoryStatus"`
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

func (t *Transaction) UpdateBalance(username string, transactionAmount float32, transactionName string, transactionDescription string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `insert into mrkrabs.Transactions (Username, TransactionAmount, TransactionName, TransactionDescription) values
	($1,$2,$3,$4)`

	balance, err := t.GetUserBalance(username)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, err
	}
	if (balance + transactionAmount) < 0 {
		return 0, errors.New("error. can not decrement balance below zero")
	}

	_, err = db.ExecContext(ctx, query, username, transactionAmount, transactionName, transactionDescription)
	if err != nil {
		return 0, err
	}

	return balance + transactionAmount, nil
}

func (t *Transaction) GetAllTransactions(username string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select TransactionID, username, transactionamount, transactionname, transactiondescription from mrkrabs.Transactions where Username = $1`

	rows, err := db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var trans Transaction
		if err := rows.Scan(&trans.TransactionID, &trans.UserID, &trans.TransactionAmount, &trans.TransactionName, &trans.TransactionDescription); err != nil {
			return transactions, err
		}
		transactions = append(transactions, trans)
	}
	if err = rows.Err(); err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (t *RecurringPayment) GetAllReccurringPayments() ([]RecurringPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT paymentid, username, paymentamount, paymentname, paymentdescription, paymentdate
	FROM foreman.recurring_payment`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurring_payments []RecurringPayment

	for rows.Next() {
		var recurring RecurringPayment
		if err := rows.Scan(&recurring.PaymentID, &recurring.UserName, &recurring.PaymentAmount, &recurring.PaymentName, &recurring.PaymentDescription, &recurring.PaymentDate); err != nil {
			return recurring_payments, err
		}
		recurring_payments = append(recurring_payments, recurring)
	}
	if err = rows.Err(); err != nil {
		return recurring_payments, err
	}
	return recurring_payments, nil
}

func (t *RecurringPayment) GetReccurringPayments(username string) ([]RecurringPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT paymentid, username, paymentamount, paymentname, paymentdescription, paymentdate
	FROM foreman.recurring_payment WHERE username = $1`

	rows, err := db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurring_payments []RecurringPayment

	for rows.Next() {
		var recurring RecurringPayment
		if err := rows.Scan(&recurring.PaymentID, &recurring.UserName, &recurring.PaymentAmount, &recurring.PaymentName, &recurring.PaymentDescription, &recurring.PaymentDate); err != nil {
			return recurring_payments, err
		}
		recurring_payments = append(recurring_payments, recurring)
	}
	if err = rows.Err(); err != nil {
		return recurring_payments, err
	}
	return recurring_payments, nil
}

func (t *RecurringPayment) AddReccurringPayment(username string, paymentAmount float32, paymentName string, paymentDescription string, paymentDate string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO foreman.recurring_payment(
		username, paymentamount, paymentname, paymentdescription, paymentdate)
		VALUES ($1, $2, $3, $4, $5);`

	_, err := db.ExecContext(ctx, query, username, paymentAmount, paymentName, paymentDescription, paymentDate)

	if err != nil {
		return 0, err
	}

	return 1, err
}

func (t *PaymentHistory) GetPaymentHistory(paymentID int) ([]PaymentHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT paymenthistoryid, paymentid, paymenthistorydate, paymenthistorystatus
				FROM foreman.payment_history WHERE paymentid = $1`

	rows, err := db.QueryContext(ctx, query, paymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []PaymentHistory

	for rows.Next() {
		var payment PaymentHistory
		if err := rows.Scan(&payment.PaymentHistoryID, &payment.PaymentID, &payment.PaymentHistoryDate, &payment.PaymentHistoryStatus); err != nil {
			return payments, err
		}
		payments = append(payments, payment)
	}
	if err = rows.Err(); err != nil {
		return payments, err
	}
	return payments, nil
}

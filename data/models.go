package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	Transaction      Transaction
	RecurringPayment RecurringPayment
	PaymentHistory   PaymentHistory
	Account          Account
	Category         Category
	Debt             Debt
}

type Transaction struct {
	TransactionID          int     `json:"transaction_id"`
	UserID                 string  `json:"user_id"`
	TransactionAmount      float32 `json:"transactionAmount"`
	TransactionName        string  `json:"transactionName"`
	TransactionDescription string  `json:"transactionDescription"`
	TransactionCategory    string  `json:"transactionCategory"`
}

type Debt struct {
	DebtID            int     `json:"debtID"`
	UserID            string  `json:"user_id"`
	TotalOwing        float32 `json:"total_owing"`
	TotalDebtPayments float32 `json:"total_payments"`
	Name string `json:"name"`
}
type DebtPayment struct {
	PaymentID     int `json:"payment_id"`
	TransactionID int `json:"transaction_id"`
}

type Category struct {
	TransactionCategory string `json:"transactionCategory"`
	Username            string `json:"username"`
}

type Account struct {
	AccountID   int    `json:"id"`
	AccountName string `json:"accountname"`
	Email       string `json:"email"`
	IsPrimary   bool   `json:"isprimary"`
}

type RecurringPayment struct {
	PaymentID          int     `json:"paymentid"`
	UserName           string  `json:"username"`
	PaymentAmount      float32 `json:"amount"`
	PaymentName        string  `json:"paymentName"`
	PaymentDescription string  `json:"paymentDescription"`
	PaymentDate        string  `json:"paymentDate"`
	PaymentType        string  `json:"paymentType"`
}

type PaymentHistory struct {
	PaymentHistoryID     int    `json:"paymenthistoryid"`
	PaymentID            int    `json:"paymentID"`
	PaymentHistoryDate   string `json:"paymentHistoryDate"`
	PaymentHistoryStatus bool   `json:"paymentHistoryStatus"`
}

func (t *Account) GetUserAccounts(email string) ([]Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT *
				FROM mrkrabs.Account WHERE username = $1`

	rows, err := db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.AccountID, &account.AccountName, &account.Email, &account.IsPrimary); err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}
	if err = rows.Err(); err != nil {
		return accounts, err
	}
	return accounts, nil
}

func (t *Account) AddAccount(email string, account_name string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO mrkrabs.Account (
		accountname, username, isprimary)
		VALUES ($1, $2, $3);`

	_, err := db.ExecContext(ctx, query, account_name, email, true)

	if err != nil {
		return 0, err
	}

	return 1, err
}

func (t *Account) AddUserToAccount(email string, account_name string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO mrkrabs.Account (
		accountname, username, isprimary)
		VALUES ($1, $2, $3);`

	_, err := db.ExecContext(ctx, query, account_name, email, false)

	if err != nil {
		return 0, err
	}

	return 1, err
}

func (t *Transaction) GetUserBalance(email string, account string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select SUM(TransactionAmount) From mrkrabs.Transactions
	where Username = $1 and accountname = $2
	group by Username`

	var totalBalance float32

	row := db.QueryRowContext(ctx, query, email, account)
	err := row.Scan(&totalBalance)

	if err != nil {
		return 0, err
	}

	return totalBalance, nil
}

func (t *Transaction) UpdateTransactionCategory(username string, account string, transactionID int, category string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
	update mrkrabs.transactions
	set category = $1
	where transactionid = $2 and username = $3 and accountname = $3
	`
	_, err := db.QueryContext(ctx, query, category, transactionID, username, account)
	if err != nil {
		return err
	}
	return nil
}

func (t *Transaction) GetAllCategories(username string, account string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
	select distinct category from mrkrabs.transactions
  where username = $1 and accountname = $2
	`

	var categories []string

	rows, err := db.QueryContext(ctx, query, username, account)
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return categories, nil
		}
		categories = append(categories, s)
	}
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (t *Transaction) UpdateBalance(username string, account string, transactionAmount float32, transactionName string, transactionDescription string, transactionCategory string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `insert into mrkrabs.Transactions (Username, AccountName, TransactionAmount, TransactionName, TransactionDescription, Category) values
	($1,$2,$3,$4,$5,$6)`

	balance, err := t.GetUserBalance(username, account)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, err
	}
	if (balance + transactionAmount) < 0 {
		return 0, errors.New("error. can not decrement balance below zero")
	}

	_, err = db.ExecContext(ctx, query, username, account, transactionAmount, transactionName, transactionDescription, transactionCategory)
	if err != nil {
		return 0, err
	}

	return balance + transactionAmount, nil
}
func (t *Transaction) GetAllTransactionsOfCategory(username, account string, category string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select TransactionID, username, transactionamount, transactionname, transactiondescription, category from mrkrabs.Transactions where Username = $1 and category = $2 and accountname = $3`

	rows, err := db.QueryContext(ctx, query, username, category, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var trans Transaction
		if err := rows.Scan(&trans.TransactionID, &trans.UserID, &trans.TransactionAmount, &trans.TransactionName, &trans.TransactionDescription, &trans.TransactionCategory); err != nil {
			return transactions, err
		}
		transactions = append(transactions, trans)
	}
	if err = rows.Err(); err != nil {
		return transactions, err
	}
	return transactions, nil
}
func (t *Transaction) GetAllTransactions(username string, account string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select TransactionID, username, transactionamount, transactionname, transactiondescription, category from mrkrabs.Transactions where Username = $1 and accountname = $2`

	rows, err := db.QueryContext(ctx, query, username, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var trans Transaction
		if err := rows.Scan(&trans.TransactionID, &trans.UserID, &trans.TransactionAmount, &trans.TransactionName, &trans.TransactionDescription, &trans.TransactionCategory); err != nil {
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

func (t *RecurringPayment) GetReccurringPayments(username string, account string) ([]RecurringPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT paymentid, username, paymentamount, paymentname, paymentdescription, paymentdate, paymenttype
	FROM foreman.recurring_payment WHERE username = $1 and accountname = $2`

	rows, err := db.QueryContext(ctx, query, username, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurring_payments []RecurringPayment

	for rows.Next() {
		var recurring RecurringPayment
		if err := rows.Scan(&recurring.PaymentID, &recurring.UserName, &recurring.PaymentAmount, &recurring.PaymentName, &recurring.PaymentDescription, &recurring.PaymentDate, &recurring.PaymentType); err != nil {
			return recurring_payments, err
		}
		recurring_payments = append(recurring_payments, recurring)
	}
	if err = rows.Err(); err != nil {
		return recurring_payments, err
	}
	return recurring_payments, nil
}

func (t *RecurringPayment) AddReccurringPayment(username string, account string, paymentAmount float32, paymentName string, paymentDescription string, paymentDate string, paymentType string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO foreman.recurring_payment(
		username, accountname, paymentamount, paymentname, paymentdescription, paymentdate, paymenttype)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := db.ExecContext(ctx, query, username, account, paymentAmount, paymentName, paymentDescription, paymentDate, paymentType)

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
func (d *Debt) GetAllDebts(userID string, account string) ([]Debt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
	SELECT
  Debt.DebtID,
  Debt.UserID,
  Debt.TotalOwing,
  COALESCE(SUM(transactions.TransactionAmount), 0) * -1 AS TotalDebtPayments,
  Debt.Name
FROM
  mrkrabs.Debt
  LEFT JOIN mrkrabs.DebtPayment
    ON Debt.DebtID = DebtPayment.DebtID
  LEFT JOIN mrkrabs.transactions
    ON DebtPayment.TransactionID = transactions.TransactionID
WHERE
  Debt.UserID = $1
  AND Debt.accountname = $2
GROUP BY
  Debt.DebtID,
  Debt.UserID,
  Debt.TotalOwing,
  Debt.Name;`
	rows, err := db.QueryContext(ctx, query, userID, account)
	if err != nil {
	log.Println("Here")

		return nil, err
	}
	defer rows.Close()
	var debts []Debt
	for rows.Next() {
		var debt Debt
		if err := rows.Scan(&debt.DebtID, &debt.UserID, &debt.TotalOwing, &debt.TotalDebtPayments, &debt.Name); err != nil {
			return debts, err
		}
		debts = append(debts, debt)
	}
	if err = rows.Err(); err != nil {
	log.Println("There")
		return debts, err
	}
	return debts, nil
}
func (d *Debt) CreateDebt(userID string, account string, totalOwing float32, name string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO mrkrabs.Debt (UserID, AccountName, TotalOwing, Name)
	VALUES ($1, $2, $3, $4)
	RETURNING DebtID;
	`

	var debtID int

	row := db.QueryRowContext(ctx, query, userID, account, totalOwing, name)
	err := row.Scan(&debtID)
	if err != nil {
		return -1, err
	}
	return debtID, nil
}

func (d *Debt) GetDebtByID(debtID int, userID string, account string) (Debt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
	SELECT
  Debt.DebtID,
  Debt.UserID,
  Debt.TotalOwing,
  COALESCE(SUM(transactions.TransactionAmount), 0) * -1 AS TotalDebtPayments,
	Debt.Name
FROM
  mrkrabs.Debt
  LEFT JOIN mrkrabs.DebtPayment
    ON Debt.DebtID = DebtPayment.DebtID
  LEFT JOIN mrkrabs.transactions
    ON DebtPayment.TransactionID = transactions.TransactionID
WHERE
  Debt.DebtID = $1
	AND Debt.UserID = $2
	AND Debt.AccountName = $3
GROUP BY
  Debt.DebtID,
  Debt.UserID,
  Debt.TotalOwing,
	Debt.Name;`

	var debt Debt
	row := db.QueryRowContext(ctx, query, debtID, userID, account)
	err := row.Scan(&debt.DebtID, &debt.UserID, &debt.TotalOwing, &debt.TotalDebtPayments, &debt.Name)
	if err != nil {
		return debt, err
	}

	return debt, nil
}
func (d *Debt) MakeDebtPayment(userID string, account string, debtID int, amount float32) (Debt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	if amount > 0 {
		amount = amount * -1
	}
	fmt.Sprintf("amt:", amount)
	var t *Transaction
	var debt Debt
	var transactionID int

	// check if debt exists
	debt, err := d.GetDebtByID(debtID, userID, account)
	if err != nil {
		return debt, err
	}

	// create transaction
	query := `insert into mrkrabs.Transactions (Username, AccountName, TransactionAmount, TransactionName, TransactionDescription, Category) values
	($1,$2,$3,$4,$5,$6)
	RETURNING TransactionID`
	balance, err := t.GetUserBalance(userID, account)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return debt, err
	}
	log.Println(balance)
	if (balance + amount) < 0 {
		return debt, errors.New("error. can not decrement balance below zero")
	}
	row := db.QueryRowContext(ctx, query, userID, account, amount, fmt.Sprintf("balance payment for debt %d", debtID), "", "Debt")
	err = row.Scan(&transactionID)
	if err != nil {
		return debt, err
	}

	// insert transaction with debt
	query = `insert into mrkrabs.DebtPayment (TransactionID, DebtID)
	values ($1,$2)`
	_, err = db.QueryContext(ctx, query, transactionID, debtID)
	if err != nil {
		return debt, err
	}

	debt, err = d.GetDebtByID(debtID, userID, account)
	if err != nil {
		return debt, err
	}
	return debt, nil
}

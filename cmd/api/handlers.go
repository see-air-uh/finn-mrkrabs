package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *Config) GetBalance(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")

	balance, err := app.Models.Transaction.GetUserBalance(u, account)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//TODO: Add logging
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Grabbed user balance for user %s", u),
		Data:    balance,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
func (app *Config) UpdateTransactionCategory(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")

	var requestPayload struct {
		TransactionID       int    `json:"transactionID"`
		TransactionCategory string `json:"transactionCategory"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.Transaction.UpdateTransactionCategory(u, account, requestPayload.TransactionID, requestPayload.TransactionCategory)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusAccepted, jsonResponse{
		Error:   true,
		Message: "updated category",
	})
}
func (app *Config) GetCategories(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")
	categories, err := app.Models.Transaction.GetAllCategories(u, account)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("successfully grabbed all categories for %s", u),
		Data:    categories,
	})
}
func (app *Config) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")
	var requestPayload struct {
		// Username          string  `json:"username"`
		TransactionAmount      float32 `json:"transactionAmount"`
		TransactionName        string  `json:"transactionName"`
		TransactionDescription string  `json:"transactionDescription"`
		TransactionCategory    string  `json:"transactionCategory"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	balance, err := app.Models.Transaction.UpdateBalance(u, account, requestPayload.TransactionAmount, requestPayload.TransactionName, requestPayload.TransactionDescription, requestPayload.TransactionCategory)
	if err != nil {

		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Updated user balance for user %s", u),
		Data:    balance,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}
func (app *Config) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")

	transactions, err := app.Models.Transaction.GetAllTransactions(u, account)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Retrieved transaction data for user %s", u),
		Data:    transactions,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}
func (app *Config) GetAllTransactionsOfCategory(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")
	c := chi.URLParam(r, "category")

	transactions, err := app.Models.Transaction.GetAllTransactionsOfCategory(u, account, c)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Retrieved transaction data for user %s", u),
		Data:    transactions,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}
func (app *Config) GetReccurringPayments(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")

	balance, err := app.Models.RecurringPayment.GetReccurringPayments(u, account)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//TODO: Add logging
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Grabbed recurring payments for user %s", u),
		Data:    balance,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
func (app *Config) GetAllReccurringPayments(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")

	balance, err := app.Models.RecurringPayment.GetAllReccurringPayments()
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//TODO: Add logging
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Grabbed recurring payments for user %s", u),
		Data:    balance,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
func (app *Config) AddReccurringPayment(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")
	var requestPayload struct {
		PaymentAmount      float32 `json:"amount"`
		PaymentName        string  `json:"paymentName"`
		PaymentDescription string  `json:"paymentDescription"`
		PaymentDate        string  `json:"paymentDate"`
		PaymentType        string  `json:"paymentType"`
		PaymentFrequency   string  `json:"paymentFrequency"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	balance, err := app.Models.RecurringPayment.AddReccurringPayment(u, account, requestPayload.PaymentAmount, requestPayload.PaymentName, requestPayload.PaymentDescription, requestPayload.PaymentDate, requestPayload.PaymentType, requestPayload.PaymentFrequency)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Added recurring payment for user %s", u),
		Data:    balance,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	r_id := chi.URLParam(r, "recurring_id")

	recurring_id, err := strconv.Atoi(r_id)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	transactions, err := app.Models.PaymentHistory.GetPaymentHistory(recurring_id)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Retrieved payment history for user %s", u),
		Data:    transactions,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")

	accounts, err := app.Models.Account.GetUserAccounts(u)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Retrieved accounts for user %s", u),
		Data:    accounts,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetAllDebts(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")

	debts, err := app.Models.Debt.GetAllDebts(u, account)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Retrieved debts for user %s", u),
		Data:    debts,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) CreateDebt(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")
	var debtPayload struct {
		TotalOwing float32 `json:"total_owing"`
	}
	err := app.readJSON(w, r, &debtPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	debt, err := app.Models.Debt.CreateDebt(u, account, debtPayload.TotalOwing)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Created debt for user %s that has id %d", u, debt),
		Data:    debt,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetDebtByID(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")
	debtIDString := chi.URLParam(r, "debtID")

	debtID, err := strconv.Atoi(debtIDString)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	debt, err := app.Models.Debt.GetDebtByID(debtID, u, account)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Retrieved debt for user %s with id %s", u, debtIDString),
		Data:    debt,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) MakeDebtPayment(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")
	debtIDString := chi.URLParam(r, "debtID")

	debtID, err := strconv.Atoi(debtIDString)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	var debtPayload struct {
		Amount float32 `json:"amount"`
	}
	err = app.readJSON(w, r, &debtPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	debt, err := app.Models.Debt.MakeDebtPayment(u, account, debtID, debtPayload.Amount)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Made debt payment for user %s with id %s", u, debtIDString),
		Data:    debt,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) AddAccount(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")

	_, err := app.Models.Account.AddAccount(u, account)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Added account for user %s", u),
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) AddUserToAccount(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	account := chi.URLParam(r, "account")
	u2 := chi.URLParam(r, "user2")

	_, err := app.Models.Account.AddUserToAccount(u2, account)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Added user to account for user %s", u),
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

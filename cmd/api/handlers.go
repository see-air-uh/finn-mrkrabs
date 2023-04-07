package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *Config) GetBalance(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")

	balance, err := app.Models.Transaction.GetUserBalance(u)
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
func (app *Config) UpdateTransactionCategory( w http.ResponseWriter, r *http.Request){
	u := chi.URLParam(r, "user")

	var requestPayload struct {
		TransactionID int `json:"transactionID"`
		TransactionCategory string `json:"transactionCategory"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.Transaction.UpdateTransactionCategory(u,requestPayload.TransactionID, requestPayload.TransactionCategory)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusAccepted, jsonResponse {
		Error: true,
		Message: "updated category",
	})
}
func (app *Config) GetCategories(w http.ResponseWriter, r *http.Request){
	u := chi.URLParam(r, "user")
	categories, err := app.Models.Transaction.GetAllCategories(u)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonResponse {
		Error: false,
		Message: fmt.Sprintf("successfully grabbed all categories for %s", u),
		Data: categories,
	})
}
func (app *Config) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	var requestPayload struct {
		// Username          string  `json:"username"`
		TransactionAmount      float32 `json:"transactionAmount"`
		TransactionName        string  `json:"transactionName"`
		TransactionDescription string  `json:"transactionDescription"`
		TransactionCategory 	 string  `json:"transactionCategory"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	balance, err := app.Models.Transaction.UpdateBalance(u, requestPayload.TransactionAmount, requestPayload.TransactionName, requestPayload.TransactionDescription, requestPayload.TransactionCategory)
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

	transactions, err := app.Models.Transaction.GetAllTransactions(u)
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
func (app *Config) GetAllTransactionsOfCategory(w http.ResponseWriter, r *http.Request){
	u := chi.URLParam(r, "user")
	c := chi.URLParam(r, "category")

	transactions, err := app.Models.Transaction.GetAllTransactionsOfCategory(u,c)
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

	balance, err := app.Models.RecurringPayment.GetReccurringPayments(u)
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
	var requestPayload struct {
		PaymentAmount      float32 `json:"amount"`
		PaymentName        string  `json:"paymentName"`
		PaymentDescription string  `json:"paymentDescription"`
		PaymentDate        string  `json:"paymentDate"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	balance, err := app.Models.RecurringPayment.AddReccurringPayment(u, requestPayload.PaymentAmount, requestPayload.PaymentName, requestPayload.PaymentDescription, requestPayload.PaymentDate)
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
	var requestPayload struct {
		PaymentID int `json:"paymentID"`
	}

	transactions, err := app.Models.PaymentHistory.GetPaymentHistory(requestPayload.PaymentID)
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
func (app *Config) GetAllDebts(w http.ResponseWriter, r* http.Request){
	u := chi.URLParam(r, "user")

	debts, err := app.Models.Debt.GetAllDebts(u)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("Retrieved debts for user %s",u),
		Data: debts,
	}
	app.writeJSON(w, http.StatusAccepted,payload)
}

func (app *Config) CreateDebt(w http.ResponseWriter, r* http.Request){
	u := chi.URLParam(r, "user")
	var debtPayload struct {
		TotalOwing float32 `json:"total_owing"`
	}
	err := app.readJSON(w, r, &debtPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	debt, err := app.Models.Debt.CreateDebt(u, debtPayload.TotalOwing)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("Created debt for user %s that has id %d",u, debt),
		Data: debt,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetDebtByID(w http.ResponseWriter, r* http.Request){
	u := chi.URLParam(r, "user")
	debtIDString := chi.URLParam(r, "debtID")

	debtID, err := strconv.Atoi(debtIDString)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	debt, err := app.Models.Debt.GetDebtByID(debtID, u)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("Retrieved debt for user %s with id %s",u, debtIDString),
		Data: debt,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) MakeDebtPayment(w http.ResponseWriter, r* http.Request){
	u := chi.URLParam(r, "user")
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
	
	debt, err := app.Models.Debt.MakeDebtPayment(u, debtID,debtPayload.Amount)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return 
	}
	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("Made debt payment for user %s with id %s",u, debtIDString),
		Data: debt,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}
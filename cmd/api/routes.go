package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/balance/{user}/{account}", app.GetBalance)
	mux.Post("/balance/{user}/{account}", app.UpdateBalance)

	mux.Get("/recurring/{user}/{account}", app.GetReccurringPayments)
	mux.Post("/recurring/add/{user}/{account}", app.AddReccurringPayment)
	mux.Post("/recurring/history/{user}", app.GetPaymentHistory)

	mux.Get("/accounts/{user}", app.GetUserAccounts)
	mux.Post("/accounts/add/{user}/{account}", app.AddAccount)
	mux.Post("/accounts/add_user/{user}/{account}/{user2}", app.AddUserToAccount)

	mux.Get("/transaction/{user}/{account}", app.GetAllTransactions)
	mux.Post("/transaction/{user}/{account}/category", app.UpdateTransactionCategory)
	mux.Get("/transaction/{user}/{account}/category", app.GetCategories)
	mux.Get("/transaction/{user}/{account}/category/{category}", app.GetAllTransactionsOfCategory)

	mux.Get("/debt/{user}/{account}", app.GetAllDebts)
	mux.Post("/debt/{user}/{account}", app.CreateDebt)
	mux.Get("/debt/{user}/{account}/{debtID}", app.GetDebtByID)
	mux.Post("/debt/{user}/{account}/{debtID}", app.MakeDebtPayment)

	return mux
}

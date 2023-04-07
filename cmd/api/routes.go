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

	mux.Get("/balance/{user}", app.GetBalance)
	mux.Post("/balance/{user}", app.UpdateBalance)
	mux.Get("/transaction/{user}", app.GetAllTransactions)
	mux.Get("/recurring/{user}", app.GetReccurringPayments)
	mux.Post("/recurring/add/{user}", app.AddReccurringPayment)
	mux.Get("/recurring/history/{user}", app.GetPaymentHistory)

	mux.Post("/transaction/{user}/category", app.UpdateTransactionCategory)
	mux.Get("/transaction/{user}/category", app.GetCategories)
	mux.Get("/transaction/{user}/category/{category}",app.GetAllTransactionsOfCategory)

	mux.Get("/debt/{user}", app.GetAllDebts)
	mux.Post("/debt/{user}",app.CreateDebt)
	mux.Get("/debt/{user}/{debtID}", app.GetDebtByID)

	return mux
}

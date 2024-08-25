package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	mux.Get("/", app.Home)
	mux.Get("/admin/virtual-terminal", app.Auth(app.VirtualTerminal))

	//mux.Post("/virtual-terminal-payment-succeeded", app.VirtualTerminalPaymentSucceeded)
	//mux.Get("/virtual-terminal-receipt", app.VirtualTerminalReceipt)

	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)

	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/widget/{id}", app.ChargeOnce)
	mux.Get("/receipt", app.Receipt)

	mux.Get("/plans/gold", app.GoldPlan)
	mux.Get("/receipt/gold", app.GoldPlanReceipt)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return app.Session.LoadAndSave(mux)
}

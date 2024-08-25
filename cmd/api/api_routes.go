package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

	// todo: fix problem admin middleware
	//mux.Route("/api/admin", func(r chi.Router) {
	//	mux.Use(app.Auth)
	//
	//	mux.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
	//		writer.WriteHeader(http.StatusOK)
	//		writer.Write([]byte("OK"))
	//	})
	//
	//	mux.Post("/virtual-terminal-successed", app.VirtualTerminalSuccessed)
	//})

	mux.Post("/api/payment-intent", app.GetPaymentIntent)
	mux.Post("/api/admin/virtual-terminal-successed", app.VirtualTerminalSuccessed)
	mux.Get("/api/widgets/{id}", app.GetWidgetByID)
	mux.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)
	mux.Post("/api/authenticate", app.CreateAuthToken)
	mux.Post("/api/is_authenticated", app.CheckAuthentication)

	return mux
}

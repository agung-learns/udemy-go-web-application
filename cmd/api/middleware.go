package main

import "net/http"

func (app *application) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := app.authenticateToken(r)
		if err != nil {
			_ = app.invalidCredentials(w)
			return
		}
		next.ServeHTTP(w, r)
	}
}

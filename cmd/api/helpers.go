package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, payload interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dst)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("data must only contain a single JSON object")
	}

	return nil
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = err.Error()

	our, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(our)

	return nil
}

func (app *application) invalidCredentials(w http.ResponseWriter) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = true
	payload.Message = "invalid authentication credentials"

	return app.writeJSON(w, http.StatusUnauthorized, payload)
}

func (app *application) passwordMatches(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

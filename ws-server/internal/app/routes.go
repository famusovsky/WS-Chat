package app

import (
	"net/http"
)

// routes - создание ServeMux для сервера чата.
//
// Возвращает: ServeMux.
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.connect)

	return mux
}

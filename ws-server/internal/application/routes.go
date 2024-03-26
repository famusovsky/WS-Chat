package application

import (
	"net/http"
)

// routes - создание ServeMux для сервера чата.
//
// Возвращает: ServeMux.
func (app *App) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.connect)

	return mux
}

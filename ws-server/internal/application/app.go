package application

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"
	"ws-server/internal/postgres"

	"github.com/gorilla/websocket"
)

// App - структура, описывающая приложение.
type App struct {
	addr     string
	infoLog  *log.Logger
	errorLog *log.Logger
	db       postgres.MessagesDBProcessor
	upgrader websocket.Upgrader
	users    map[*User]struct{}
	mu       sync.Mutex
	srvr     *http.Server
}

// CreateApp - создание модели приложения.
//
// Принимает: адрес, логгер информации, логгер ошибок, БД, указатель на необходимость пересоздать БД.
//
// Возвращает: сервер, ошибку.
func CreateApp(addr string, infoLog *log.Logger, errorLog *log.Logger, db *sql.DB, overrideTables bool) (*App, error) {
	repo, err := postgres.Get(db, overrideTables)
	if err != nil {
		return nil, err
	}

	res := &App{
		addr:     addr,
		infoLog:  infoLog,
		errorLog: errorLog,
		db:       repo,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Пропускаем любой запрос
			},
		},
		users: make(map[*User]struct{}),
		mu:    sync.Mutex{},
	}
	res.srvr = &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  res.routes(),
	}

	return res, nil
}

// Run - запуск сервера.
func (app *App) Run() {
	app.infoLog.Printf("App started on address: %s\n", app.addr)

	go func() {
		err := app.srvr.ListenAndServe()
		app.errorLog.Fatal(err)
	}()
}

// Shutdown - изящное отключение сервера.
func (app *App) Shutdown() {
	app.errorLog.Fatal(app.srvr.Shutdown(context.Background()))
}

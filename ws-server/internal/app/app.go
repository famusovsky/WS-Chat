package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"ws-server/internal/postgres"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

// AppModel - структура, описывающая модель приложения.
type AppModel struct {
	addr     string
	infoLog  *log.Logger
	errorLog *log.Logger
	db       postgres.MessagesDBProcessor
}

// AppModel - структура, описывающая приложение.
type application struct {
	*AppModel
	upgrader websocket.Upgrader
	users    map[*User]struct{}
	mu       sync.Mutex
}

// CreateModel - создание модели приложения.
//
// Принимает: адрес, логгер информации, логгер ошибок, БД, указатель на необходимость пересоздать БД.
//
// Возвращает: модель приложения, ошибку.
func CreateModel(addr string, infoLog *log.Logger, errorLog *log.Logger, db *sql.DB, overrideTables bool) (*AppModel, error) {
	repo, err := postgres.Get(db, overrideTables)
	if err != nil {
		return nil, err
	}

	return &AppModel{
		addr:     addr,
		infoLog:  infoLog,
		errorLog: errorLog,
		db:       repo,
	}, nil
}

// Run - запуск приложения.
func (model *AppModel) Run() {
	app := &application{
		AppModel: model,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Пропускаем любой запрос
			},
		},
		users: make(map[*User]struct{}),
		mu:    sync.Mutex{},
	}

	srvr := &http.Server{
		Addr:     app.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("App started on address: %s\n", app.addr) //

	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		}
	})

	go func() {
		err := srvr.ListenAndServe()
		app.errorLog.Fatal(err)
	}()

	if err := eg.Wait(); err != nil {
		model.infoLog.Printf("gracefully shutting down the server: %v\n", err)
	}

	err := srvr.Shutdown(context.Background())
	app.errorLog.Fatal(err)
}

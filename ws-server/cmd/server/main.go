package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"ws-server/internal/application"
	"ws-server/pkg/database"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP address")
	overrideTables := flag.Bool("override_tables", false, "Override tables in database")
	dsn := flag.String("dsn", "", "dsn for the db")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

	var db *sql.DB
	var err error
	if *dsn == "" {
		db, err = database.OpenViaEnvVars("postgres")
	} else {
		db, err = database.OpenViaDsn(*dsn, "postgres")
	}
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app, err := application.CreateApp(*addr, infoLog, errorLog, db, *overrideTables)
	if err != nil {
		errorLog.Fatal(err)
	}

	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	eg := new(errgroup.Group)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		}
	})

	app.Run()

	if err := eg.Wait(); err != nil {
		infoLog.Printf("gracefully shutting down the server: %v\n", err)
		app.Shutdown()
	}
}

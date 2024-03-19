package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"ws-server/internal/app"
	"ws-server/pkg/database"

	_ "github.com/lib/pq"
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

	model, err := app.CreateModel(*addr, infoLog, errorLog, db, *overrideTables)
	if err != nil {
		errorLog.Fatal(err)
	}
	model.Run()
}

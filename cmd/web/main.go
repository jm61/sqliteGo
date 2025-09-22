package main

import (
	"chinook/internal/models"
	"crypto/tls"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	users          *models.UserModel
	employees      *models.EmployeeModel
	albums         *models.AlbumModel
}

var dict = []string{}
var records [][]string
var dataMap (map[string]string)

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address (default is 8080), port must start at 1024")
	flag.Parse()
	dsn := "ch.db"

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))

	app := &application{
		logger: logger,
	}

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	db, err := app.openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	dict, records = createDictionary()
	dataMap = createMap()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 1 * time.Hour

	app = &application{
		logger:         logger,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		users:          &models.UserModel{DB: db},
		employees:      &models.EmployeeModel{DB: db},
		albums:         &models.AlbumModel{DB: db},
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "addr", srv.Addr)
	fmt.Printf("https://127.0.0.1%s\n", srv.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func (app *application) openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	app.logger.Info("Connected to DB")
	return db, nil
}

func createDictionary() ([]string, [][]string) {
	file, err := os.Open("artists.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, row := range records {
		dict = append(dict, row...)
	}
	return dict, records
}

func createMap() map[string]string {
	dataMap = make(map[string]string)

	for _, pair := range records {
		if len(pair) == 2 {
			dataMap[pair[0]] = pair[1]
		}
	}

	fmt.Println("Datamap: ", len(dataMap))
	return dataMap
}

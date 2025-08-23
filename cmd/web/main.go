package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mike-jacks/snippetbox/config"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	var cfg config.Config

	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.DSN, "dsn", "postgres://web:web@127.0.0.1:5433/snippetbox?sslmode=disable", "Postgres DB dsn")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "If program is verbose")
	flag.BoolVar(&cfg.Trace, "trace", false, "If add trace to logger")
	flag.Parse()

	db, err := openDB(cfg.DSN, l)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := config.NewApplication(cfg, db, l)

	if app.Config().Verbose {
		app.Logger().Info("Verbose Logging Enabled")
	}

	app.Logger().Info("starting server", "addr", fmt.Sprintf("http://localhost%s", app.Config().Addr))

	err = http.ListenAndServe(app.Config().Addr, app.Routes())
	app.Logger().Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string, l *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	l.Info("Database ping success!")

	return db, nil
}

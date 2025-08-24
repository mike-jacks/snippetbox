package config

import (
	"database/sql"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/mike-jacks/snippetbox/internal/models"
)

type ApplicaitonInterface interface {
	ServerError(w http.ResponseWriter, r *http.Request, err error)
	ClientError(w http.ResponseWriter, status int)
	Logger() *slog.Logger
	Config() Config
	Snippets() *models.SnippetModel
}

type Application struct {
	config     Config
	logger     *slog.Logger
	fileServer http.Handler
	handler    *Handler
	snippets   *models.SnippetModel
}

func NewApplication(config Config, db *sql.DB, logger *slog.Logger) *Application {
	fileServer := http.FileServer(http.Dir(config.StaticDir))

	app := &Application{
		config:     config,
		logger:     logger,
		fileServer: fileServer,
		snippets: &models.SnippetModel{
			DB: db,
		},
	}

	app.handler = &Handler{
		app: app,
	}

	return app
}

func (app *Application) Logger() *slog.Logger {
	return app.logger
}

func (app *Application) Config() Config {
	return app.config
}

func (app *Application) FileServer() http.Handler {
	return app.fileServer
}

func (app *Application) Handler() *Handler {
	return app.handler
}

func (app *Application) Snippets() *models.SnippetModel {
	return app.snippets
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	logFields := map[string]string{
		"error": err.Error(),
	}
	if app.Config().Verbose {
		logFields["method"] = r.Method
		logFields["uri"] = r.URL.RequestURI()
	}

	if app.Config().Trace {
		logFields["trace"] = string(debug.Stack())
	}

	fields := make([]interface{}, 0, len(logFields)*2)
	if value, exists := logFields["error"]; exists {
		fields = append(fields, "error", value)
	}
	if value, exists := logFields["method"]; exists {
		fields = append(fields, "method", value)
	}
	if value, exists := logFields["uri"]; exists {
		fields = append(fields, "uri", value)
	}
	if value, exists := logFields["trace"]; exists {
		fields = append(fields, "trace", value)
	}

	app.Logger().Error("Internal server error occured", fields...)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", app.Handler().Home())
	mux.Handle("GET /snippet/view/{id}", app.Handler().SnippetView())
	mux.Handle("GET /snippet/create", app.Handler().SnippetCreate())
	mux.Handle("POST /snippet/create", app.Handler().SnippetCreatePost())
	mux.Handle("GET /static/", app.FileServer())

	return mux
}

package config

import (
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
)

type ApplicaitonInterface interface {
	ServerError(w http.ResponseWriter, r *http.Request, err error)
	ClientError(w http.ResponseWriter, status int)
	Logger() *slog.Logger
	Config() Config
}

type Application struct {
	config     Config
	logger     *slog.Logger
	fileServer http.Handler
	handler    *Handler
}

func NewApplication(config Config) *Application {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	fileServer := http.FileServer(http.Dir(config.StaticDir))

	app := &Application{
		config:     config,
		logger:     logger,
		fileServer: fileServer,
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
	mux.Handle("GET /static/", http.StripPrefix("/static/", app.FileServer()))

	return mux
}

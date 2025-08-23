package config

import (
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
)

type Application struct {
	Config     Config
	Logger     *slog.Logger
	FileServer http.Handler
}

type Config struct {
	Addr      string
	StaticDir string
	Verbose   bool
	Trace     bool
}

func NewApplication(config Config) *Application {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	fileServer := http.FileServer(http.Dir(config.StaticDir))

	return &Application{
		Config:     config,
		Logger:     logger,
		FileServer: fileServer,
	}
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	logFields := map[string]string{
		"error": err.Error(),
	}
	if app.Config.Verbose {
		logFields["method"] = r.Method
		logFields["uri"] = r.URL.RequestURI()
	}

	if app.Config.Trace {
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

	app.Logger.Error("Internal server error occured", fields...)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

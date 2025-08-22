package config

import (
	"log/slog"
	"net/http"
	"os"
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


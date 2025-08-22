package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mike-jacks/snippetbox/config"
)

func main() {
	var cfg config.Config

	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "If program is verbose")
	flag.Parse()

	app := config.NewApplication(cfg)

	if app.Config.Verbose {
		app.Logger.Info("Verbose Logging Enabled")
	}

	mux := http.NewServeMux()
	registerMux(mux, app)

	app.Logger.Info("starting server", "addr", fmt.Sprintf("http://localhost%s", app.Config.Addr))

	err := http.ListenAndServe(app.Config.Addr, mux)
	app.Logger.Error(err.Error())
	os.Exit(1)
}

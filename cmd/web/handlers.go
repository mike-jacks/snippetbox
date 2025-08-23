package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/mike-jacks/snippetbox/config"
)

func Home(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/pages/home.tmpl",
			"./ui/html/partials/nav.tmpl",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ServerError(w,r,err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.ServerError(w,r,err)
		}
	})
}

func SnippetView(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			if err != nil { app.Logger.Error(err.Error()) }
			http.NotFound(w,r)
			return
		}
		fmt.Fprintf(w, "Display a specific sinppet with ID %d...", id)
	})
}

func SnippetCreate(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Display a form for creating a new snippet...")
	})
}

func SnippetCreatePost(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Save a new snippet...")
	})
}


func registerMux(mux *http.ServeMux, app *config.Application) {
	mux.Handle("GET /{$}", Home(app))
	mux.Handle("GET /snippet/view/{id}", SnippetView(app))
	mux.Handle("GET /snippet/create", SnippetCreate(app))
	mux.Handle("POST /snippet/create", SnippetCreatePost(app))
	mux.Handle("GET /static/", http.StripPrefix("/static/", app.FileServer))
}

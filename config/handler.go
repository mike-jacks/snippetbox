package config

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	app ApplicaitonInterface
}

func NewHandler(app ApplicaitonInterface) *Handler {
	return &Handler{app: app}
}

func (h *Handler) Home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/pages/home.tmpl",
			"./ui/html/partials/nav.tmpl",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			h.app.ServerError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			h.app.ServerError(w, r, err)
		}
	})
}

func (h *Handler) SnippetView() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			if err != nil {
				h.app.Logger().Error(err.Error())
			}
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Display a specific sinppet with ID %d...", id)
	})
}

func (h *Handler) SnippetCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Display a form for creating a new snippet...")
	})
}

func (h *Handler) SnippetCreatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Save a new snippet...")
	})
}


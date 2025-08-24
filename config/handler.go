package config

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/mike-jacks/snippetbox/internal/models"
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

		snippets, err := h.app.Snippets().Latest()
		if err != nil {
			h.app.ServerError(w, r, err)
			return
		}

		for _, snippet := range snippets {
			fmt.Fprintf(w, "%+v\n", &snippet)
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
		snippet, err := h.app.Snippets().Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				h.app.ServerError(w, r, err)
			}
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/view.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			h.app.ServerError(w, r, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", snippet)
		if err != nil {
			h.app.ServerError(w, r, err)
		}
	})
}

func (h *Handler) SnippetCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Display a form for creating a new snippet...")
	})
}

func (h *Handler) SnippetCreatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := "0 snail"
		content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
		expires := 7

		id, err := h.app.Snippets().Insert(title, content, expires)
		if err != nil {
			h.app.ServerError(w, r, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	})
}

package router

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/michaeldebetaz/unilike/internal/env"
	"github.com/michaeldebetaz/unilike/internal/middleware"
)

func Start() error {
	url := fmt.Sprintf("%s:%s", env.ORIGIN(), env.PORT())

	msg := fmt.Sprintf("Server is running at %s", url)
	slog.Info(msg)

	http.HandleFunc("GET /", route(homeHandler))
	http.HandleFunc("GET /faculties", route(facultiesHandler))
	http.HandleFunc("GET /scrapper", route(scrapperHandler))

	return http.ListenAndServe(url, nil)
}

func route(f http.HandlerFunc) http.HandlerFunc {
	return middleware.Chain(f, middleware.Logging())
}

func serveTemplate(w http.ResponseWriter, templ *template.Template, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := templ.ExecuteTemplate(w, "layout", data); err != nil {
		text := fmt.Sprintf("Error executing template: %v", err)
		http.Error(w, text, http.StatusInternalServerError)
		return
	}
}

func templPath(name string) string {
	fp := filepath.Join("templates", name)
	return fp
}

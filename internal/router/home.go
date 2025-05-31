package router

import (
	"fmt"
	"html/template"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	layout := templPath("layout.html")
	home := templPath("home.html")

	templ, err := template.ParseFiles(layout, home)
	if err != nil {
		text := fmt.Sprintf("Error parsing template files: %v", err)
		http.Error(w, text, http.StatusInternalServerError)
		return
	}

	serveTemplate(w, templ, nil)
}

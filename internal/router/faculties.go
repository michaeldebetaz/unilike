package router

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/michaeldebetaz/unilike/internal/db"
)

func facultiesHandler(w http.ResponseWriter, r *http.Request) {
	layout := templPath("layout.html")
	faculties := templPath("faculties.html")

	templ, err := template.ParseFiles(layout, faculties)
	if err != nil {
		text := fmt.Sprintf("Error parsing template files: %v", err)
		http.Error(w, text, http.StatusInternalServerError)
		return
	}

	data, err := db.LoadFromJson()
	if err != nil {
		text := fmt.Sprintf("Error loading data from JSON: %v", err)
		http.Error(w, text, http.StatusInternalServerError)
		return
	}

	if len(data.Faculties) < 1 {
		http.Error(w, "No faculties found", http.StatusNotFound)
		return
	}

	serveTemplate(w, templ, data.Faculties)
}

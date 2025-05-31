package router

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/michaeldebetaz/unilike/internal/cache"
	"github.com/michaeldebetaz/unilike/internal/db"
	"github.com/michaeldebetaz/unilike/internal/env"
	"github.com/michaeldebetaz/unilike/internal/parser"
	"github.com/michaeldebetaz/unilike/internal/scrapper"
)

func scrapperHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Starting scraping process...")

	cache, err := cache.Load()
	if err != nil {
		log.Fatalf("Failed to load cache: %v", err)
	}

	facultiesUrl := env.BASE_PATH() + "index.php?v_langue=fr&v_isinterne="
	facultiesHtml, err := scrapper.GetHtml(facultiesUrl, cache)
	if err != nil {
		log.Fatalf("Failed to get faculties page html: %v", err)
	}

	faculties, err := parser.ExtractFaculties(facultiesHtml)
	if err != nil {
		log.Fatalf("Failed to extract faculties: %v", err)
	}

	data := db.Db{}

FACULTIES:
	for _, faculty := range faculties {
		slog.Info("Faculty", "ueid", faculty.Ueid, "name", faculty.Name)

		facultyHtml, err := scrapper.GetHtml(faculty.Url, cache)
		if err != nil {
			slog.Warn("Failed to get faculty page html:", "faculty url", faculty.Url, "error", err)
			continue FACULTIES
		}
		// faculty.Html = facultyHtml

		programs, err := parser.ExtractPrograms(facultyHtml)
		if err != nil {
			log.Fatalf("Failed to extract programs: %v", err)
		}

	PROGRAMS:
		for _, program := range programs {
			slog.Info("Program", "etapeId1", program.EtapeId1, "name", program.Name)

			program.Filename = fmt.Sprintf("%s_%s", faculty.Filename, program.Filename)

			programHtml, err := scrapper.GetHtml(program.Url, cache)
			if err != nil {
				slog.Warn("Program page not found:", "program url", program.Url, "error", err)
				continue PROGRAMS
			}
			program.Html = programHtml

			courses, err := parser.ExtractCourses(programHtml)
			if err != nil {
				log.Fatalf("Failed to extract courses: %v", err)
			}

		COURSES:
			for _, course := range courses {
				slog.Info("Course", "enstyId", course.EnstyId, "name", course.Name)

				course.Filename = fmt.Sprintf("%s_%s", program.Filename, course.Filename)

				html, err := scrapper.GetHtml(course.Url, cache)
				if err != nil {
					slog.Warn("Failed to get course page html:", "course url", course.Url, "error", err)
					continue COURSES
				}
				course.Html = html

				teachers, err := parser.ExtractCourseTeachers(course.Html)
				if err != nil {
					slog.Warn("Failed to extract course teachers:", "course url", course.Url, "error", err)
					continue COURSES
				}

				course.Teachers = teachers

				program.Courses = append(program.Courses, course)
			}
			faculty.Programs = append(faculty.Programs, program)
		}
		data.Faculties = append(data.Faculties, faculty)
	}

	w.Header().Set("Content-Type", "text/plain")

	if err := data.SaveAsJson(); err != nil {
		msg := fmt.Sprintf("Failed to save data to JSON: %v", err)
		slog.Error(msg)
		http.Error(w, fmt.Sprintf("Failed to save data to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Scraping process completed successfully with %d faculties.", len(data.Faculties))
	slog.Info(msg)
	fmt.Fprintln(w, msg)
}

package scrapper

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/michaeldebetaz/unilscrap/internal/cache"
	"github.com/michaeldebetaz/unilscrap/internal/db"
	"github.com/michaeldebetaz/unilscrap/internal/env"
	"github.com/michaeldebetaz/unilscrap/internal/parser"
)

func Scrape() {
	slog.Info("Starting scraping process...")

	cache, err := cache.Load()
	if err != nil {
		log.Fatalf("Failed to load cache: %v", err)
	}

	facultiesUrl := env.BASE_PATH() + "index.php?v_langue=fr&v_isinterne="
	facultiesHtml, err := getHtml(facultiesUrl, cache)
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

		facultyHtml, err := getHtml(faculty.Url, cache)
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
			slog.Info("Program", "etapeId-1", program.EtapeId1, "name", program.Name)

			program.Filename = fmt.Sprintf("%s_%s", faculty.Filename, program.Filename)

			programHtml, err := getHtml(program.Url, cache)
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

				html, err := getHtml(course.Url, cache)
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

	if err := data.SaveAsJson(); err != nil {
		msg := fmt.Sprintf("Failed to save data to JSON: %v", err)
		slog.Error(msg)
	}

	slog.Info("Scraping process completed successfully.")
}

func getHtml(url string, cache cache.Cache) (string, error) {
	if html, ok := cache.Get(url); ok {
		slog.Info("Cache hit", "url", url)
		return html, nil
	}

	fmt.Printf("Visiting %s\n", url)

	res, err := http.Get(url)
	if err != nil {
		err := fmt.Errorf("Error while doing Get request: %v", err)
		return "", err
	}
	defer res.Body.Close()

	fmt.Printf("Response status: %s\n", res.Status)

	if res.StatusCode > 299 {
		err := fmt.Errorf("Error: %s", res.Status)
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		err := fmt.Errorf("Error while reading body: %v", err)
		return "", err
	}

	html := string(body)
	cache.Set(url, html)

	return html, nil
}

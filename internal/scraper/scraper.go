package scraper

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/michaeldebetaz/unilscrap/internal/cache"
	"github.com/michaeldebetaz/unilscrap/internal/db"
	"github.com/michaeldebetaz/unilscrap/internal/env"
	"github.com/michaeldebetaz/unilscrap/internal/parser"
)

type facultyResult struct {
	index   int
	faculty db.Faculty
}

func Scrape() {
	slog.Info("Starting scraping process...")

	cache, err := cache.Load()
	if err != nil {
		log.Fatalf("Failed to load cache: %v", err)
	}

	facultiesUrl := env.BASE_PATH() + "index.php?v_langue=fr&v_isinterne="
	facultiesHtml, err := fetchHtml(facultiesUrl, cache)
	if err != nil {
		log.Fatalf("Failed to get faculties page html: %v", err)
	}

	faculties := fetchFaculties(facultiesHtml, cache)

	for i, faculty := range faculties {
		programs := fetchPrograms(faculty, cache)
		faculties[i].Programs = programs
	}

	for i, faculty := range faculties {
		for j, program := range faculty.Programs {
			courses := fecthCourses(program, cache)
			faculties[i].Programs[j].Courses = courses
		}
	}

	data := db.Data{Faculties: faculties}
	data.Debug()

	data.SaveAsJson()

	cache.Save()

	slog.Info("Scraping process completed successfully.")
}

func fetchFaculties(html string, cache cache.Cache) []db.Faculty {
	faculties, err := parser.ExtractFaculties(html)
	if err != nil {
		log.Fatalf("Failed to extract faculties: %v", err)
	}

	var wg sync.WaitGroup

	for i, faculty := range faculties {
		wg.Add(1)
		go func(idx int, faculty db.Faculty) {
			defer wg.Done()
			slog.Info("Faculty", "ueid", faculty.Ueid, "name", faculty.Name)

			facultyHtml, err := fetchHtml(faculty.Url, cache)
			if err != nil {
				slog.Warn("Failed to get faculty page html:", "faculty url", faculty.Url, "error", err)
			}
			faculty.Html = facultyHtml

			programs, err := parser.ExtractPrograms(faculty.Html)
			if err != nil {
				log.Fatalf("Failed to extract programs: %v", err)
			}
			faculty.Programs = programs

			faculties[idx] = faculty
		}(i, faculty)
	}

	wg.Wait()

	return faculties
}

func fetchPrograms(faculty db.Faculty, cache cache.Cache) []db.Program {
	programs, err := parser.ExtractPrograms(faculty.Html)
	if err != nil {
		log.Fatalf("Failed to extract programs for faculty %s: %v", faculty.Ueid, err)
	}

	var wg sync.WaitGroup

	for i, program := range programs {
		wg.Add(1)
		go func(idx int, program db.Program) {
			defer wg.Done()

			slog.Info("Program", "etapeId-1", program.EtapeId1, "name", program.Name)

			programHtml, err := fetchHtml(program.Url, cache)
			if err != nil {
				slog.Warn("Program page not found:", "program url", program.Url, "error", err)
			}
			program.Html = programHtml

			program.Courses, err = parser.ExtractCourses(programHtml)
			if err != nil {
				slog.Warn("Failed to extract courses:", "program url", program.Url, "error", err)
			}

			program.Filename = fmt.Sprintf("%s_%s", faculty.Filename, program.Filename)

			programs[idx] = program
		}(i, program)
	}

	wg.Wait()

	return programs
}

func fecthCourses(program db.Program, cache cache.Cache) []db.Course {
	var wg sync.WaitGroup

	for i, course := range program.Courses {
		wg.Add(1)
		go func(idx int, course db.Course) {
			defer wg.Done()

			slog.Info("Course", "enstyId", course.EnstyId, "name", course.Name)

			courseHtml, err := fetchHtml(course.Url, cache)
			if err != nil {
				slog.Warn("Failed to get course page html:", "course url", course.Url, "error", err)
			}
			course.Html = courseHtml

			teachers, err := parser.ExtractCourseTeachers(courseHtml)
			if err != nil {
				slog.Warn("Failed to extract course teachers:", "course url", course.Url, "error", err)
			}

			course.Teachers = teachers
			course.Filename = fmt.Sprintf("%s_%s", program.Filename, course.Filename)

			program.Courses[idx] = course
		}(i, course)
	}

	wg.Wait()

	return program.Courses
}

func fetchHtml(url string, cache cache.Cache) (string, error) {
	if html, ok := cache.Get(url); ok {
		slog.Info("Cache hit", "url", url)
		return html, nil
	}

	slog.Info("Fetching HTML", "url", url)

	res, err := http.Get(url)
	if err != nil {
		err := fmt.Errorf("Error while doing Get request: %v", err)
		return "", err
	}
	defer res.Body.Close()

	slog.Info("Response status", "status", res.Status)

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

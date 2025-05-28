package main

import (
	"fmt"
	"log"

	"github.com/michaeldebetaz/unilike/internal/db"
	"github.com/michaeldebetaz/unilike/internal/env"
	"github.com/michaeldebetaz/unilike/internal/parser"
	"github.com/michaeldebetaz/unilike/internal/scrapper"
)

func main() {
	env := env.GetEnv("BASE_PATH")
	cache, err := scrapper.LoadCache()
	if err != nil {
		log.Fatalf("Failed to load cache: %v", err)
	}
	defer cache.Save()

	facultiesUrl := env.BASE_PATH + "index.php?v_langue=fr&v_isinterne="

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
		fmt.Printf("Faculty: %s - %s\n", faculty.Ueid, faculty.Name)

		facultyHtml, err := scrapper.GetHtml(faculty.Url, cache)
		if err != nil {
			fmt.Printf("Faculty page not found: %s\n", faculty.Url)
			continue FACULTIES
		}
		faculty.Html = facultyHtml

		programs, err := parser.ExtractPrograms(facultyHtml)
		if err != nil {
			log.Fatalf("Failed to extract programs: %v", err)
		}

	PROGRAMS:
		for _, program := range programs {
			fmt.Printf("Program: %s - %s\n", program.EtapeId1, program.Name)

			program.Filename = fmt.Sprintf("%s_%s", faculty.Filename, program.Filename)

			programHtml, err := scrapper.GetHtml(program.Url, cache)
			if err != nil {
				fmt.Printf("Program page not found: %s\n", program.Url)
				continue PROGRAMS
			}
			program.Html = programHtml

			courses, err := parser.ExtractCourses(programHtml)
			if err != nil {
				log.Fatalf("Failed to extract courses: %v", err)
			}

		COURSES:
			for _, course := range courses {
				fmt.Printf("Course: %s - %s\n", course.EnstyId, course.Name)
				course.Filename = fmt.Sprintf("%s_%s", program.Filename, course.Filename)

				html, err := scrapper.GetHtml(course.Url, cache)
				if err != nil {
					fmt.Printf("Course page not found: %s\n", course.Url)
					continue COURSES
				}
				course.Html = html

				teachers, err := parser.ExtractCourseTeachers(course.Html)
				if err != nil {
					fmt.Printf("Failed to extract course details: %v\n", err)
					continue COURSES
				}

				course.Teachers = teachers

				program.Courses = append(program.Courses, course)
			}
			faculty.Programs = append(faculty.Programs, program)
		}
		data.Faculties = append(data.Faculties, faculty)
	}

	data.Debug()
}

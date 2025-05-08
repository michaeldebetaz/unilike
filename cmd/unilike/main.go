package main

import (
	"fmt"
	"log"

	"github.com/michaeldebetaz/unilike/internal/db"
	"github.com/michaeldebetaz/unilike/internal/env"
	"github.com/michaeldebetaz/unilike/internal/parser"
	"github.com/michaeldebetaz/unilike/internal/scrapper"
	"github.com/michaeldebetaz/unilike/internal/writer"
)

func main() {
	env := env.GetEnv("BASE_PATH")

	facultiesUrl := env.BASE_PATH + "index.php?v_langue=fr&v_isinterne="
	facultiesHtml, err := scrapper.GetHtml(facultiesUrl)
	if err != nil {
		log.Fatalf("Failed to get faculties page html: %v", err)
	}

	faculties, err := parser.ExtractFaculties(facultiesHtml)
	if err != nil {
		log.Fatalf("Failed to extract faculty hrefs: %v", err)
	}

	data := db.Db{}

	for _, faculty := range faculties {
		facultyHtml, err := scrapper.GetHtml(faculty.Url.String())
		if err != nil {
			log.Fatalf("Failed to get faculty page html: %v", err)
		}

		programs, err := parser.ExtractPrograms(facultyHtml)
		if err != nil {
			log.Fatalf("Failed to extract programs: %v", err)
		}

		for _, program := range programs {
			programHtml, err := scrapper.GetHtml(program.Url.String())
			if err != nil {
				log.Fatalf("Failed to get program page html: %v", err)
			}

			order := fmt.Sprintf("%001d", program.Order)
			fileName := faculty.FileName + "_" + order + "_" + program.FileName + ".html"
			if err := writer.ToFile(fileName, programHtml); err != nil {
				log.Fatalf("Failed to write program html to file: %v", err)
			}

			courses, err := parser.ExtractCourses(programHtml)
			if err != nil {
				log.Fatalf("Failed to extract courses: %v", err)
			}

			program.Courses = courses

			faculty.Programs = append(faculty.Programs, program)
		}

		data.Faculties = append(data.Faculties, faculty)
	}

	data.Debug()
}

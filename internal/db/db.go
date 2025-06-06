package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/michaeldebetaz/unilscrap/internal/writer"
)

type Data struct {
	Faculties []Faculty `json:"faculties"`
}

type Faculty struct {
	Order    int       `json:"order"`
	Ueid     string    `json:"ueid"`
	Name     string    `json:"name"`
	Filename string    `json:"filename"`
	Url      string    `json:"url"`
	Html     string    `json:"html"`
	Programs []Program `json:"programs"`
}

type Program struct {
	Order          int      `json:"order"`
	SemPosSelected int      `json:"semposselected"`
	EtapeId1       string   `json:"etapeid1"`
	Name           string   `json:"name"`
	Filename       string   `json:"filename"`
	Url            string   `json:"url"`
	Html           string   `json:"html"`
	Courses        []Course `json:"courses"`
}

type Course struct {
	Order    int    `json:"order"`
	EnstyId  string `json:"enstyid"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
	Html     string `json:"html"`
	Teachers string `json:"teachers"`
}

func (db *Data) Debug() {
	for _, faculty := range db.Faculties {
		fmt.Printf("Faculty: %s - %s\n", faculty.Ueid, faculty.Name)
		fmt.Printf("  URL: %s\n", faculty.Url)
		for _, program := range faculty.Programs {
			fmt.Printf("  Program: %s - %s\n", program.EtapeId1, program.Name)
			fmt.Printf("    URL: %s\n", program.Url)
			for _, course := range program.Courses {
				fmt.Printf("    Course: %s - %s\n", course.EnstyId, course.Name)
				fmt.Printf("      Teachers: %s\n", course.Teachers)
				fmt.Printf("      URL: %s\n", course.Url)
			}
		}
	}
}

func (db *Data) SaveAsJson() {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		slog.Error("Failed to marschal db to JSON", "error", err.Error())
	}

	filename := "db.json"
	if err := writer.ToFile(filename, string(data)); err != nil {
		slog.Error("Failed to write db to file", "filename", filename, "error", err.Error())
	}

	slog.Info("Data saved as JSON.")
}

func LoadFromJson() (Data, error) {
	db := Data{}

	filePath := "db.json"

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return db, nil
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return db, fmt.Errorf("Error reading db from file: %v", err)
	}
	if err = json.Unmarshal(bytes, &db); err != nil {
		return db, fmt.Errorf("Error unmarshalling db: %v", err)
	}

	return db, nil
}

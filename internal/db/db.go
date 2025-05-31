package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/michaeldebetaz/unilike/internal/writer"
)

type Db struct {
	Faculties []Faculty `json:"faculties"`
}

type Faculty struct {
	Order    int    `json:"order"`
	Ueid     string `json:"ueid"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
	// Html     string    `json:"html"`
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

func (db *Db) Debug() {
	for _, faculty := range db.Faculties {
		fmt.Printf("Faculty: %s - %s\n", faculty.Ueid, faculty.Name)
		for _, program := range faculty.Programs {
			fmt.Printf("  Program: %s - %s\n", program.EtapeId1, program.Name)
			for _, course := range program.Courses {
				fmt.Printf("    Course: %s - %s\n", course.EnstyId, course.Name)
				fmt.Printf("      Teachers: %s\n", course.Teachers)
				fmt.Printf("      URL: %s\n", course.Url)
			}
		}
	}
}

func (db *Db) SaveAsJson() error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal db to JSON: %w", err)
	}

	filename := "db.json"
	if err := writer.ToFile(filename, string(data)); err != nil {
		return fmt.Errorf("failed to write db to file %s: %w", filename, err)
	}

	return nil
}

func LoadFromJson() (Db, error) {
	db := Db{}

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

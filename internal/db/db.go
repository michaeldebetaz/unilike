package db

import "fmt"

type Db struct {
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

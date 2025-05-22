package db

import (
	"encoding/json"
	"fmt"
)

type Db struct {
	Faculties []Faculty `json:"faculties"`
}

type Faculty struct {
	Order     int       `json:"order"`
	Ueid      string    `json:"ueid"`
	Name      string    `json:"name"`
	ClassName string    `json:"class_name"`
	FileName  string    `json:"file_name"`
	Url       string    `json:"url"`
	Programs  []Program `json:"programs"`
}

type Program struct {
	Order          int      `json:"order"`
	SemPosSelected int      `json:"sem_pos_selected"`
	EtapeId1       string   `json:"etape_id_1"`
	Name           string   `json:"name"`
	FileName       string   `json:"file_name"`
	Url            string   `json:"url"`
	Courses        []Course `json:"courses"`
}

type Course struct {
	Order        int    `json:"order"`
	SemesterName string `json:"semester_name"`
	Ueid         string `json:"ueid"`
	Name         string `json:"name"`
	Teacher      string `json:"teacher"`
}

func (d *Db) Debug() {
	bytes, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data:", err)
	}
	fmt.Println(string(bytes))
}

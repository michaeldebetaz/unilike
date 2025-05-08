package parser

import (
	"github.com/michaeldebetaz/unilike/internal/db"
)

func ExtractCourses(html string) ([]db.Course, error) {
	courses := []db.Course{}

	// node, err := parseHtml(html)
	// if err != nil {
	// 	err := fmt.Errorf("Failed to parse HTML: %v", err)
	// 	return courses, err
	// }
	//
	// id := "UniDocContent"
	// node, err = getElementById(node, id)
	// if err != nil {
	// 	err := fmt.Errorf("Failed to find element with id %s: %v", id, err)
	// 	return courses, err
	// }

	return courses, nil
}

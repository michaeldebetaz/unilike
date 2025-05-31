package parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/michaeldebetaz/unilike/internal/assert"
	"github.com/michaeldebetaz/unilike/internal/db"
	"github.com/michaeldebetaz/unilike/internal/env"
)

func ExtractCourses(html string) ([]db.Course, error) {
	courses := []db.Course{}

	htmlNode, err := parseHtml(html)
	if err != nil {
		err := fmt.Errorf("Failed to parse HTML: %v", err)
		return courses, err
	}

	uniDocContentDivNode, err := htmlNode.id("UniDocContent")
	if err != nil {
		err := fmt.Errorf("Failed to find element with id UniDocContent: %v", err)
		return courses, err
	}

	uniDocContentUlNodes := uniDocContentDivNode.tags("ul")
	if uniDocContentUlNodes.len() < 1 {
		return courses, nil
	}

	uniDocContentUlNode := uniDocContentUlNodes.first()
	uniDocContentLiNodes := uniDocContentUlNode.tags("li")

	order := 1

	for _, uniDocContentLiNode := range *uniDocContentLiNodes {
		course := db.Course{}

		UniDocContentANode := uniDocContentLiNode.tags("a").first()

		name := UniDocContentANode.innerText()
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, " -")
		name = strings.ReplaceAll(name, "-[", "- [")
		course.Name = name

		onClick := UniDocContentANode.attributeValue("onclick")

		re := regexp.MustCompile(`window\.open\('([^']+)'`)
		matches := re.FindStringSubmatch(onClick)

		href := assert.At(matches, 1)
		u, err := url.Parse(env.BASE_PATH() + href)
		if err != nil {
			err := fmt.Errorf("Failed to parse URL: %v", err)
			return courses, err
		}
		course.Url = u.String()

		values := u.Query()
		course.EnstyId = values.Get("v_enstyid")

		course.Order = order
		order++

		index := fmt.Sprintf("03%d", order)
		course.Filename = fmt.Sprintf("%s_%s", index, course.EnstyId)

		courses = append(courses, course)
	}

	return courses, nil
}

func ExtractCourseTeachers(html string) (string, error) {
	htmlNode, err := parseHtml(html)
	if err != nil {
		err := fmt.Errorf("Failed to parse HTML: %v", err)
		return "", err
	}

	uniDocContentDivNode, err := htmlNode.id("UniDocContent")
	if err != nil {
		err := fmt.Errorf("Failed to find element with id UniDocContent: %v", err)
		return "", err
	}

	uniDocContentPNode := uniDocContentDivNode.tags("p").first()
	teachers := uniDocContentPNode.innerText()
	teachers = strings.TrimSpace(teachers)
	teachers = strings.ReplaceAll(teachers, "Responsables(s):", "Responsable(s) :")
	teachers = strings.TrimSuffix(teachers, "Intervenant(s): -")
	teachers = strings.ReplaceAll(teachers, "Intervenant(s):", "; Intervenant(s) :")

	return teachers, nil
}

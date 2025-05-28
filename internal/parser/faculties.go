package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/michaeldebetaz/unilike/internal/db"
	"github.com/michaeldebetaz/unilike/internal/env"
)

func ExtractFaculties(html string) ([]db.Faculty, error) {
	faculties := []db.Faculty{}

	htmlNode, err := parseHtml(html)
	if err != nil {
		err := fmt.Errorf("Failed to parse HTML: %v", err)
		return faculties, err
	}

	id := "UniDocContent"
	uniDocContentDivNode, err := htmlNode.id(id)
	if err != nil {
		err := fmt.Errorf("Failed to find element with id %s: %v", id, err)
		return faculties, err
	}

	factListDivNode := uniDocContentDivNode.classes("fac-list").first()

	// Get the url for each faculty
	tag := "a"
	facListANodes := factListDivNode.tags(tag)
	if facListANodes.len() < 1 {
		err := fmt.Errorf("Failed to find elements with tag %s", tag)
		return faculties, err
	}

	for i, aNode := range *facListANodes {
		order := i + 1
		index := fmt.Sprintf("%03d", order)

		h5Node := aNode.tags("h5").first()
		name := h5Node.innerText()

		href := aNode.attributeValue("href")
		url, error := url.Parse(env.BASE_PATH + href)
		if error != nil {
			err := fmt.Errorf("Failed to parse URL: %v", error)
			return faculties, err
		}

		values := url.Query()
		ueid := values.Get("v_ueid")

		classNames := strings.Split(aNode.attributeValue("class"), " ")
		lastIndex := len(classNames) - 1
		className := classNames[lastIndex]

		faculty := db.Faculty{
			Order:    order,
			Ueid:     ueid,
			Name:     name,
			Filename: fmt.Sprintf("%s_faculty_%s_%s", index, ueid, className),
			Url:      url.String(),
		}

		faculties = append(faculties, faculty)
	}

	return faculties, nil
}

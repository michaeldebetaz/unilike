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

	node, err := parseHtml(html)
	if err != nil {
		err := fmt.Errorf("Failed to parse HTML: %v", err)
		return faculties, err
	}

	id := "UniDocContent"
	node, err = getElementById(node, id)
	if err != nil {
		err := fmt.Errorf("Failed to find element with id %s: %v", id, err)
		return faculties, err
	}

	class := "fac-list"
	node = getElementsByClass(node, class)[0]

	// Get the url for each faculty
	tag := "a"
	nodes := getElementsByTag(node, tag)
	if len(nodes) == 0 {
		err := fmt.Errorf("Failed to find elements with tag %s", tag)
		return faculties, err
	}

	for i, n := range nodes {
		order := i + 1
		index := fmt.Sprintf("%03d", order)

		nodes := getElementsByTag(n, "h5")
		node = nodes[0]
		name := getInnerText(node)

		href := getAttributeValue(n, "href")
		url, error := url.Parse(env.BASE_PATH + href)
		if error != nil {
			err := fmt.Errorf("Failed to parse URL: %v", error)
			return faculties, err
		}

		values := url.Query()
		ueid := values.Get("v_ueid")

		classNames := strings.Split(getAttributeValue(n, "class"), " ")
		lastIndex := len(classNames) - 1
		className := classNames[lastIndex]

		faculty := db.Faculty{
			Order:     order,
			Ueid:      ueid,
			Name:      name,
			ClassName: className,
			FileName:  fmt.Sprintf("%s_faculty_%s_%s", index, ueid, className),
			Url:       url,
		}

		faculties = append(faculties, faculty)
	}

	return faculties, nil
}

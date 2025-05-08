package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/michaeldebetaz/unilike/internal/db"
	"github.com/michaeldebetaz/unilike/internal/env"
)

func ExtractPrograms(html string) ([]db.Program, error) {
	programs := []db.Program{}

	node, err := parseHtml(html)
	if err != nil {
		err := fmt.Errorf("Failed to parse HTML: %v", err)
		return programs, err
	}

	id := "UniDocContent"
	node, err = getElementById(node, id)
	if err != nil {
		err := fmt.Errorf("Failed to find element with id %s: %v", id, err)
		return programs, err
	}

	class := "listeEtapes"
	node = getElementsByClass(node, class)[0]

	// Get the name of each curriculum
	tag := "tr"
	trNodes := getElementsByTag(node, tag)

	for _, n := range trNodes {
		program := db.Program{}

		class := "tdNomEtape"
		nameNode := getElementsByClass(n, class)
		if len(nameNode) == 0 {
			continue
		}
		name := getInnerText(nameNode[0])
		program.Name = name

		class = "liens"
		node = getElementsByClass(n, class)[0]

		tag := "tr"
		trNodes = getElementsByTag(node, tag)

		for _, n := range trNodes {
		}

		tag = "a"
		aNodes := getElementsByTag(node, tag)

		for _, aNode := range aNodes {
			href := getAttributeValue(aNode, "href")

			if strings.HasPrefix(href, "listeCours.php") {
				url, error := url.Parse(env.BASE_PATH + href)
				if error != nil {
					err := fmt.Errorf("Failed to parse URL: %v", error)
					return programs, err
				}

				semPosSelected := url.Query().Get("v_semposselected")
				etapeId1 := url.Query().Get("v_etapeid1")

				program.Url = url
				program.EtapeId1 = etapeId1
				program.FileName = fmt.Sprintf("%s", etapeId1)

				programs = append(programs, program)
			}
		}
	}

	return programs, nil
}

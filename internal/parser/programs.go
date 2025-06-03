package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/michaeldebetaz/unilscrap/internal/db"
	"github.com/michaeldebetaz/unilscrap/internal/env"
)

func ExtractPrograms(html string) ([]db.Program, error) {
	programs := []db.Program{}

	htmlNode, err := parseHtml(html)
	if err != nil {
		err := fmt.Errorf("Failed to parse HTML: %v", err)
		return programs, err
	}

	id := "UniDocContent"
	uniDocContentDivNode, err := htmlNode.id(id)
	if err != nil {
		err := fmt.Errorf("Failed to find element with id %s: %v", id, err)
		return programs, err
	}

	listeEtapesNode := uniDocContentDivNode.classes("listeEtapes").first()
	listeEtapesTrNodes := listeEtapesNode.tags("tr")

	order := 0

	etapeTitle := ""
	_ = etapeTitle

	for _, listeEtapesTrNode := range *listeEtapesTrNodes {
		program := db.Program{}

		etapeTitleTdNodes := listeEtapesTrNode.classes("etapeTitle")
		if etapeTitleTdNodes.len() > 0 {
			etapeTitle = etapeTitleTdNodes.first().innerText()
		}

		nomEtapeTdNodes := listeEtapesTrNode.classes("tdNomEtape")

		if nomEtapeTdNodes.len() < 1 {
			continue
		}

		nomEtape := nomEtapeTdNodes.first().innerText()
		liensTableNode := listeEtapesTrNode.classes("liens").first()
		liensTrNodes := liensTableNode.tags("tr")

		for _, liensTrNode := range *liensTrNodes {
			formNodes := liensTrNode.tags("form")

			// For programs with searchable date
			if formNodes.len() > 0 {
				inputNodes := formNodes.first().tags("input")

				v := url.Values{}

				for _, inputNode := range *inputNodes {
					name := inputNode.attributeValue("name")
					value := inputNode.attributeValue("value")

					if strings.HasPrefix(name, "etape_") {
						v.Set("v_date", value)
					}

					if name == "v_ueid" || name == "v_langue" || name == "v_etapeid1" {
						v.Set(name, value)
					}
				}

				date := v.Get("v_date")
				order++
				program.Order = order

				program.Name = etapeTitle + " - " + nomEtape + " - " + date
				u, err := url.Parse(env.BASE_PATH() + "listeCours.php?" + v.Encode())
				if err != nil {
					err := fmt.Errorf("Failed to parse URL: %v", err)
					return programs, err
				}

				program.Url = u.String()

				etapeId1 := v.Get("v_etapeid1")
				program.EtapeId1 = etapeId1

				index := fmt.Sprintf("%03d", order)

				program.Filename = fmt.Sprintf("%s_program_%s", index, etapeId1)
				program.SemPosSelected = -1

				programs = append(programs, program)

			}

			tag := "td"

			liensTdNodes := liensTrNode.tags(tag)

			if liensTdNodes.len() > 0 {
				semester := liensTdNodes.first().innerText()

				liensANodes := liensTrNode.tags("a")

				for _, aNode := range *liensANodes {
					href := aNode.attributeValue("href")

					if strings.HasPrefix(href, "listeCours.php") {
						order++
						program.Order = order

						program.Name = etapeTitle + " - " + nomEtape + " - " + semester

						u, error := url.Parse(env.BASE_PATH() + href)
						if error != nil {
							err := fmt.Errorf("Failed to parse URL: %v", error)
							return programs, err
						}
						program.Url = u.String()

						semPosSelectedStr := u.Query().Get("v_semposselected")
						semPosSelected, err := strconv.Atoi(semPosSelectedStr)
						if err != nil {
							err := fmt.Errorf("Failed to convert semPosSelected to int: %v", err)
							return programs, err
						}
						program.SemPosSelected = semPosSelected

						etapeId1 := u.Query().Get("v_etapeid1")
						program.EtapeId1 = etapeId1

						index := fmt.Sprintf("%03d", order)
						program.Filename = fmt.Sprintf("%s_%s", index, etapeId1)

						programs = append(programs, program)
					}
				}
			}
		}
	}

	return programs, nil
}

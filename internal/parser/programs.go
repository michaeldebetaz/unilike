package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/michaeldebetaz/unilike/internal/db"
	"github.com/michaeldebetaz/unilike/internal/env"
)

func ExtractPrograms(faculty db.Faculty, html string) ([]db.Program, error) {
	fmt.Println("Faculty:", faculty.Ueid, "-", faculty.Name)

	programs := []db.Program{}

	htmlNode, err := parseHtml(html)
	if err != nil {
		err := fmt.Errorf("Failed to parse HTML: %v", err)
		return programs, err
	}

	id := "UniDocContent"
	uniDocContentDivNode, err := getElementById(htmlNode, id)
	if err != nil {
		err := fmt.Errorf("Failed to find element with id %s: %v", id, err)
		return programs, err
	}

	class := "listeEtapes"
	listeEtapesNode := getFirstNode(getElementsByClass(uniDocContentDivNode, class), "listeEtapes")

	tag := "tr"
	listeEtapesTrNodes := getElementsByTag(listeEtapesNode, tag)

	order := 1

	etapeTitle := ""
	_ = etapeTitle

	for _, listeEtapesTrNode := range listeEtapesTrNodes {
		program := db.Program{}

		class := "etapeTitle"
		etapeTitleTdNodes := getElementsByClass(listeEtapesTrNode, class)
		if len(etapeTitleTdNodes) > 0 {
			etapeTitle = getInnerText(getFirstNode(etapeTitleTdNodes, "etapeTitle"))
		}

		class = "tdNomEtape"
		nomEtapeTdNodes := getElementsByClass(listeEtapesTrNode, class)

		if len(nomEtapeTdNodes) < 1 {
			continue
		}

		nomEtape := getInnerText(getFirstNode(nomEtapeTdNodes, "tdNomEtape"))

		class = "liens"
		liensTableNode := getFirstNode(getElementsByClass(listeEtapesTrNode, class), "liens")

		tag = "tr"
		liensTrNodes := getElementsByTag(liensTableNode, tag)

		for _, liensTrNode := range liensTrNodes {
			tag = "form"
			formNodes := getElementsByTag(liensTrNode, tag)

			// For programs with searchable date
			if len(formNodes) > 0 {
				formNode := getFirstNode(formNodes, "form")

				tag := "input"
				inputNodes := getElementsByTag(formNode, tag)

				v := url.Values{}

				for _, inputNode := range inputNodes {
					name := getAttributeValue(inputNode, "name")
					value := getAttributeValue(inputNode, "value")

					if strings.HasPrefix(name, "etape_") {
						v.Set("v_date", value)
					}

					if name == "v_ueid" || name == "v_langue" || name == "v_etapeid1" {
						v.Set(name, value)
					}
				}

				date := v.Get("v_date")
				program.Order = order
				order++

				program.Name = etapeTitle + " - " + nomEtape + " - " + date
				u, err := url.Parse(env.BASE_PATH + "listeCours.php?" + v.Encode())
				if err != nil {
					err := fmt.Errorf("Failed to parse URL: %v", err)
					return programs, err
				}

				fmt.Println("Program:", program.Name)

				program.Url = u.String()

				etapeId1 := v.Get("v_etapeid1")
				program.EtapeId1 = etapeId1

				index := fmt.Sprintf("%03d", order)

				program.FileName = fmt.Sprintf("%s_%s_%s", faculty.FileName, index, etapeId1)
				program.SemPosSelected = -1

				programs = append(programs, program)

			}

			tag := "td"

			liensTdNodes := getElementsByTag(liensTrNode, tag)

			if len(liensTdNodes) > 0 {
				semester := getInnerText(getFirstNode(liensTdNodes, "liensTdNode"))

				tag = "a"
				liensANodes := getElementsByTag(liensTrNode, tag)

				for _, aNode := range liensANodes {
					href := getAttributeValue(aNode, "href")

					if strings.HasPrefix(href, "listeCours.php") {
						program.Order = order
						order++

						program.Name = etapeTitle + " - " + nomEtape + " - " + semester

						u, error := url.Parse(env.BASE_PATH + href)
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
						program.FileName = fmt.Sprintf("%s_%s_%s", faculty.FileName, index, etapeId1)

						programs = append(programs, program)
					}
				}
			}
		}
	}

	return programs, nil
}

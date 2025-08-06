package scraper

import (
	"encoding/csv"
	"log"

	"github.com/gocolly/colly/v2"
)

type VillageData struct {
	Index      int // start from 1
	PostalCode string
	Name       string
	AreaCode   string

	DistrictAreaCode string
	CityAreaCode     string
	ProvinceAreaCode string
}

const (
	VILLAGE_INDEX = iota
	VILLAGE_POSTAL_CODE
	VILLAGE_NAME
	VILLAGE_AREA_CODE
)

func ScrapeVillage(e *colly.HTMLElement, writer *csv.Writer) {
	villageCount := e.DOM.Children().Length() - 1

	e.ForEach("tr", func(i int, h *colly.HTMLElement) {
		// skip last row
		if i >= villageCount {
			return
		}
		child := h.DOM.Children()

		name := child.Eq(VILLAGE_NAME).Text()
		index := child.Eq(VILLAGE_INDEX).Text()
		postalCode := child.Eq(VILLAGE_POSTAL_CODE).Text()
		areaCode := child.Eq(VILLAGE_AREA_CODE).Text()

		err := writer.Write([]string{
			index,
			name,
			postalCode,
			areaCode,
		})

		if err != nil {
			log.Fatalln("Failed to write CSV row:", err)
		}
	})

	writer.Flush()
}

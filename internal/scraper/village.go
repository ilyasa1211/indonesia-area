package scraper

import (
	"log"
	"strconv"

	"github.com/gocolly/colly/v2"
)

type VillageData struct {
	Index      int    `json:"index"` // start from 1
	PostalCode string `json:"postal_code"`
	Name       string `json:"name"`
	AreaCode   string `json:"area_code"`

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

func ScrapeVillage(c *colly.Collector, url string, selector string) []VillageData {
	out := make([]VillageData, 0)

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		villageCount := e.DOM.Children().Length()

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

			indexInt, err := strconv.Atoi(index)

			if err != nil {
				log.Fatalln("Failed to convert index to int:", err)
			}

			out = append(out, VillageData{
				Index:      indexInt,
				Name:       name,
				PostalCode: postalCode,
				AreaCode:   areaCode,
			})
		})

	})

	c.Visit(url)

	return out
}

package scraper

import (
	"encoding/csv"
	"log"

	"github.com/gocolly/colly/v2"
	"github.com/ilyasa1211/indonesia-area/utils"
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

func ScrapeVillage(c *colly.Collector, url string, selector string, writer *csv.Writer) {
	utils.SetProperHeader(c)
	utils.SetErrorHandling(c)

	defer writer.Flush()

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

	})

	c.Visit(url)
}

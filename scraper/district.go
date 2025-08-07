package scraper

import (
	"encoding/csv"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/ilyasa1211/indonesia-area/utils"
)

type DistrictData struct {
	Index        int
	Name         string
	PostalCode   []string // could be more than one
	TotalVillage int
	TotalIsland  *int
	AreaCode     string

	CityAreaCode     string
	ProvinceAreaCode string
}

const (
	DISTRICT_INDEX = iota
	DISTRICT_NAME
	DISTRICT_POSTAL_CODE
	DISTRICT_TOTAL_VILLAGE
	DISTRICT_TOTAL_ISLAND
	DISTRICT_AREA_CODE
)

func ScrapeDistrict(c *colly.Collector, url string, selector string, writer *csv.Writer) {
	utils.SetProperHeader(c)
	utils.SetErrorHandling(c)

	defer writer.Flush()

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		districtCount := e.DOM.Children().Length()

		e.ForEach("tr", func(i int, h *colly.HTMLElement) {
			// skip last row
			if i >= districtCount {
				return
			}

			child := h.DOM.Children()

			name := child.Eq(DISTRICT_NAME).Text()
			index := child.Eq(DISTRICT_INDEX).Text()
			totalVillage := strings.ReplaceAll(child.Eq(DISTRICT_TOTAL_VILLAGE).Text(), ".", "")
			totalIsland := strings.ReplaceAll(child.Eq(DISTRICT_TOTAL_ISLAND).Text(), ".", "")
			postalCode := strings.Split(child.Eq(DISTRICT_POSTAL_CODE).Text(), " - ")
			areaCode := child.Eq(DISTRICT_AREA_CODE).Text()

			err := writer.Write([]string{
				(index),
				name,
				(totalVillage),
				(totalIsland),
				areaCode,
				strings.Join(postalCode, "-"),
			})

			if err != nil {
				log.Fatalln("Failed to write CSV row:", err)
			}
		})

	})

	c.Visit(url)
}

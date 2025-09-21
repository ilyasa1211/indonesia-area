package scraper

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type DistrictData struct {
	Index        int      `json:"index"` // start from 1
	Name         string   `json:"name"`
	PostalCode   []string `json:"postal_code"` // could be more than one
	TotalVillage int      `json:"total_village"`
	TotalIsland  int      `json:"total_island"`
	AreaCode     string   `json:"area_code"`

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

func ScrapeDistrict(c *colly.Collector, url string, selector string) []DistrictData {
	out := make([]DistrictData, 0)

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

			indexInt, err := strconv.Atoi(index)

			if err != nil {
				log.Fatalln("Failed to write CSV row:", err)
			}

			postalCodes := make([]string, 0, len(postalCode))
			postalCodes = append(postalCodes, postalCode...)

			totalVillageInt, err := strconv.Atoi(totalVillage)

			if err != nil {
				totalVillageInt = 0
			}

			totalIslandInt, err := strconv.Atoi(totalIsland)

			if err != nil {
				totalIslandInt = 0
			}

			out = append(out, DistrictData{
				Index:        indexInt,
				Name:         name,
				PostalCode:   postalCodes,
				TotalVillage: totalVillageInt,
				TotalIsland:  totalIslandInt,
				AreaCode:     areaCode,
			})
		})

	})

	c.Visit(url)

	return out
}

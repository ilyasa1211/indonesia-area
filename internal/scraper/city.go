package scraper

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type CityType string

type CityData struct {
	Index int      `json:"index"` // start from 1
	Type  CityType `json:"type"`  // Kabupaten / Kota
	Name  string   `json:"name"`
	// AllocationPostalCode
	// RangePostalCode
	TotalDistrict int    `json:"total_district"`
	TotalVillage  int    `json:"total_village"`
	TotalIsland   int    `json:"total_island"`
	AreaCode      string `json:"area_code"`

	ProvinceAreaCode string
}

const (
	CITY_INDEX = iota
	CITY_TYPE
	CITY_NAME
	CITY_ALLOCATION_POSTAL_CODE
	CITY_RANGE_POSTAL_CODE
	CITY_TOTAL_DISTRICT
	CITY_TOTAL_VILLAGE
	CITY_TOTAL_ISLAND
	CITY_AREA_CODE
)

func ScrapeCity(c *colly.Collector, url string, selector string) []CityData {
	out := make([]CityData, 0)

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		cityCount := e.DOM.Children().Length()

		e.ForEach("tr", func(i int, h *colly.HTMLElement) {
			// skip last row
			if i >= cityCount {
				return
			}

			child := h.DOM.Children()

			name := child.Eq(CITY_NAME).Text()
			index := child.Eq(CITY_INDEX).Text()
			cityType := child.Eq(CITY_TYPE).Text()
			totalDistrict := child.Eq(CITY_TOTAL_DISTRICT).Text()
			totalVillage := strings.ReplaceAll(child.Eq(CITY_TOTAL_VILLAGE).Text(), ".", "") // 1.000 -> 1000
			totalIsland := strings.ReplaceAll(child.Eq(CITY_TOTAL_ISLAND).Text(), ".", "")   // 1.000 -> 1000
			areaCode := child.Eq(CITY_AREA_CODE).Text()

			indexInt, err := strconv.Atoi(index)

			if err != nil {
				log.Fatalln("failed to covert index to int: ", err)
			}

			totalDistrictInt, err := strconv.Atoi(totalDistrict)

			if err != nil {
				totalDistrictInt = 0
			}

			totalVillageInt, err := strconv.Atoi(totalVillage)

			if err != nil {
				totalVillageInt = 0
			}

			totalIslandInt, err := strconv.Atoi(totalIsland)

			if err != nil {
				totalIslandInt = 0
			}

			out = append(out, CityData{
				Index:         indexInt,
				Type:          CityType(cityType),
				Name:          name,
				TotalDistrict: totalDistrictInt,
				TotalVillage:  totalVillageInt,
				TotalIsland:   totalIslandInt,
				AreaCode:      areaCode,
			})
		})
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("Finished", r.Request.URL)
	})

	c.Visit(url)

	return out
}

package scraper

import (
	"encoding/csv"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

type CityType string

type CityData struct {
	Index int      // start from 1
	Type  CityType // Kabupaten / Kota
	Name  string
	// AllocationPostalCode
	// RangePostalCode
	TotalDistrict int
	TotalVillage  int
	TotalIsland   *int
	AreaCode      string

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

func ScrapeCity(e *colly.HTMLElement, writer *csv.Writer) {
	cityCount := e.DOM.Children().Length() - 1

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
		totalVillage := strings.ReplaceAll(child.Eq(CITY_TOTAL_VILLAGE).Text(), ".", "")
		totalIsland := strings.ReplaceAll(child.Eq(CITY_TOTAL_ISLAND).Text(), ".", "")

		areaCode := child.Eq(CITY_AREA_CODE).Text()

		err := writer.Write([]string{
			(index),
			cityType,
			name,
			(totalDistrict),
			(totalVillage),
			(totalIsland),
			areaCode,
		})

		if err != nil {
			log.Fatalln("failed to write city csv: ", err)
		}
	})

	writer.Flush()
}

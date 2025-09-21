package scraper

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ProvinceData struct {
	Index int    `json:"index"` // start from 1
	Name  string `json:"name"`
	// OneDigitPostalCode 2
	// AllocationPostalCode 23xxx-24xxx
	// RangePostalCode 23111-24794
	TotalRegencyAndCity int    `json:"total_regency_and_city"`
	TotalRegency        int    `json:"total_regency"`
	TotalCity           int    `json:"total_city"`
	TotalDistrict       int    `json:"total_district"`
	TotalVillage        int    `json:"total_village"`
	TotalIsland         int    `json:"total_island"`
	AreaCode            string `json:"area_code"`
}

const (
	PROVINCE_INDEX = iota
	PROVINCE_NAME
	PROVINCE_ONE_DIGIT_POSTAL_CODE
	PROVINCE_ALLOCATION_POSTAL_CODE
	PROVINCE_RANGE_POSTAL_CODE
	PROVINCE_TOTAL_CITY_AND_REGENCY
	PROVINCE_TOTAL_CITY
	PROVINCE_TOTAL_REGENCY
	PROVINCE_TOTAL_DISTRICT
	PROVINCE_TOTAL_VILLAGE
	PROVINCE_TOTAL_ISLAND
	PROVINCE_AREA_CODE
)

func ScrapeProvince(c *colly.Collector, url string, selector string) []ProvinceData {
	out := make([]ProvinceData, 0)

	c.OnHTML(selector, func(h *colly.HTMLElement) {
		provinceCount := h.DOM.Children().Length() - 1

		h.ForEach("tr", func(i int, h *colly.HTMLElement) {
			// Skip the last row and header
			if i >= provinceCount {
				return
			}

			child := h.DOM.Children()

			index := child.Eq(PROVINCE_INDEX).Text()
			provinceName := child.Eq(PROVINCE_NAME).Text()
			totalRegCity := child.Eq(PROVINCE_TOTAL_CITY_AND_REGENCY).Text()
			totalReg := child.Eq(PROVINCE_TOTAL_REGENCY).Text()
			totalCity := child.Eq(PROVINCE_TOTAL_CITY).Text()
			totalDistrict := child.Eq(PROVINCE_TOTAL_DISTRICT).Text()
			totalVillage := strings.ReplaceAll(child.Eq(PROVINCE_TOTAL_VILLAGE).Text(), ".", "")
			totalIsland := strings.ReplaceAll(child.Eq(PROVINCE_TOTAL_ISLAND).Text(), ".", "")
			areaCode := child.Eq(PROVINCE_AREA_CODE).Text()

			indexInt, err := strconv.Atoi(index)

			if err != nil {
				log.Fatalln("Failed to convert index to int:", err)
			}

			totalRegencyAndCityInt, err := strconv.Atoi(totalRegCity)

			if err != nil {
				totalRegencyAndCityInt = 0
			}

			totalRegencyInt, err := strconv.Atoi(totalReg)

			if err != nil {
				totalRegencyInt = 0
			}

			totalCityInt, err := strconv.Atoi(totalCity)

			if err != nil {
				totalCityInt = 0
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

			out = append(out, ProvinceData{
				Index:               indexInt,
				Name:                provinceName,
				TotalRegencyAndCity: totalRegencyAndCityInt,
				TotalRegency:        totalRegencyInt,
				TotalCity:           totalCityInt,
				TotalDistrict:       totalDistrictInt,
				TotalVillage:        totalVillageInt,
				TotalIsland:         totalIslandInt,
				AreaCode:            areaCode,
			})
		})
	})

	c.Visit(url)

	return out
}

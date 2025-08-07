package scraper

import (
	"encoding/csv"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/ilyasa1211/indonesia-area/utils"
)

type ProvinceData struct {
	Index int // start from 1
	Name  string
	// OneDigitPostalCode 2
	// AllocationPostalCode 23xxx-24xxx
	// RangePostalCode 23111-24794
	TotalRegencyAndCity int
	TotalRegency        int
	TotalCity           int
	TotalDistrict       int
	TotalVillage        int
	TotalIsland         int
	AreaCode            string
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

func ScrapeProvince(c *colly.Collector, url string, selector string, writer *csv.Writer) {
	utils.SetProperHeader(c)
	utils.SetErrorHandling(c)

	defer writer.Flush()

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

			// Write data to CSV
			err := writer.Write([]string{
				(index),
				provinceName,
				(totalRegCity),
				(totalReg),
				(totalCity),
				(totalDistrict),
				(totalVillage),
				totalIsland,
				areaCode,
			})

			if err != nil {
				log.Fatalln("Failed to write CSV row:", err)
			}
		})
	})

	c.Visit(url)
}

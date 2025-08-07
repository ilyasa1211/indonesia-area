package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"

	"github.com/gocolly/colly/v2"
	"github.com/ilyasa1211/indonesia-area/config"
	"github.com/ilyasa1211/indonesia-area/scraper"
	"github.com/ilyasa1211/indonesia-area/types"
	"github.com/ilyasa1211/indonesia-area/utils"
)

func main() {
	co := colly.NewCollector(func(c *colly.Collector) {
	})

	provinceFile := utils.OpenFile("out/csv/provinces.csv")
	cityFile := utils.OpenFile("out/csv/cities.csv")
	districtFile := utils.OpenFile("out/csv/districts.csv")
	villageFile := utils.OpenFile("out/csv/villages.csv")

	defer provinceFile.Close()
	defer cityFile.Close()
	defer districtFile.Close()
	defer villageFile.Close()

	provinceWriter := csv.NewWriter(provinceFile)
	cityWriter := csv.NewWriter(cityFile)
	districtWriter := csv.NewWriter(districtFile)
	villageWriter := csv.NewWriter(villageFile)

	co.OnError(func(r *colly.Response, err error) {
		/**
		* for some unkown reason, it sometimes return error.
		* it could be from cloudflare's protection.
		* retry if error
		 */
		if r.StatusCode == 403 {
			log.Println("error with code forbidden, will retry")
			r.Request.Visit(r.Request.URL.String())

			return
		}
		log.Println("Something went wrong:", r.StatusCode, r, err)
	})

	provinceCount := scraper.GetDataCount(co.Clone(), utils.GetURL(types.FETCH_TYPE_PROVINCE, 1, 1).String(), config.HTML_DATA_COUNT_SELECTOR)
	cityCount := scraper.GetDataCount(co.Clone(), utils.GetURL(types.FETCH_TYPE_CITY, 1, 1).String(), config.HTML_DATA_COUNT_SELECTOR)
	districtCount := scraper.GetDataCount(co.Clone(), utils.GetURL(types.FETCH_TYPE_DISTRICT, 1, 1).String(), config.HTML_DATA_COUNT_SELECTOR)
	villageCount := scraper.GetDataCount(co.Clone(), utils.GetURL(types.FETCH_TYPE_VILLAGE, 1, 1).String(), config.HTML_DATA_COUNT_SELECTOR)

	fmt.Println([]int{
		int(provinceCount),
		int(cityCount),
		int(districtCount),
		int(villageCount),
	})

	perPage := config.MAX_DATA_PER_PAGE

	p := math.Ceil(float64(provinceCount) / float64(perPage))
	c := math.Ceil(float64(cityCount) / float64(perPage))
	d := math.Ceil(float64(districtCount) / float64(perPage))
	v := math.Ceil(float64(villageCount) / float64(perPage))

	for i := 1; i <= int(p); i++ {
		scraper.ScrapeProvince(co.Clone(), utils.GetURL(types.FETCH_TYPE_PROVINCE, uint(i), uint(perPage)).String(), config.HTML_DATA_SELECTOR, provinceWriter)
	}
	for i := 1; i <= int(c); i++ {
		scraper.ScrapeCity(co.Clone(), utils.GetURL(types.FETCH_TYPE_CITY, uint(i), uint(perPage)).String(), config.HTML_DATA_SELECTOR, cityWriter)
	}
	for i := 1; i <= int(d); i++ {
		scraper.ScrapeDistrict(co.Clone(), utils.GetURL(types.FETCH_TYPE_DISTRICT, uint(i), uint(perPage)).String(), config.HTML_DATA_SELECTOR, districtWriter)
	}
	for i := 1; i <= int(v); i++ {
		scraper.ScrapeVillage(co.Clone(), utils.GetURL(types.FETCH_TYPE_VILLAGE, uint(i), uint(perPage)).String(), config.HTML_DATA_SELECTOR, villageWriter)
	}
}

package main

import (
	"log"
	"math"
	"net/http"
	"path"

	"github.com/gocolly/colly/v2"
	"github.com/ilyasa1211/indonesia-area/internal/common"
	"github.com/ilyasa1211/indonesia-area/internal/scraper"
	"github.com/ilyasa1211/indonesia-area/internal/writer"
)

func setupCallback(c *colly.Collector) {
	c.OnError(func(r *colly.Response, err error) {
		// for some unkown reason, it sometimes return error.
		// it could be from cloudflare's protection.
		// retry if error
		if r.StatusCode == http.StatusForbidden {
			log.Println("error with code forbidden, will retry")
			r.Request.Visit(r.Request.URL.String())

			return
		}

		log.Println("Something went wrong:", r.StatusCode, r, err)
	})

	c.OnRequest(func(r *colly.Request) {
		header := r.Headers

		header.Add("User-Agent", common.FAKE_USER_AGENT)
		header.Add("Accept", "text/html")
		header.Add("Accept-Language", "en-US,en;")
		header.Add("Upgrade-Insecure-Requests", "1")
		header.Add("Sec-Fetch-Dest", "document")
		header.Add("Sec-Fetch-Mode", "navigate")
		header.Add("Sec-Fetch-Site", "same-origin")
	})
}

func getAreaCount(co *colly.Collector, fetchType common.FetchType) int {
	return scraper.GetDataCount(
		co.Clone(),
		common.GetURL(fetchType, 1, 1).String(),
		common.HTML_DATA_COUNT_SELECTOR,
	)
}

func scrapeArea[T any, U any](
	co *colly.Collector,
	fetchType common.FetchType,
	perPage int,
	totalCount int,
	scrapeFunc func(*colly.Collector, string, string) []T,
	writers ...common.FileWriter[U],
) {
	pages := int(math.Ceil(float64(totalCount) / float64(perPage)))

	c := make(chan bool)

	go func() {
		res := make([]T, 0)
		for page := 1; page <= pages; page++ {
			url := common.GetURL(fetchType, page, perPage).String()
			c := co.Clone()
			setupCallback(c)
			data := scrapeFunc(c, url, common.HTML_DATA_SELECTOR)

			res = append(res, data...)
		}

		for _, w := range writers {
			file := common.OpenFile(path.Join(common.OUT_DIR, common.FetchTypeToName[fetchType]+w.GetExtension()))
			defer file.Close()

			err := w.Write(file, any(res).(U))

			if err != nil {
				log.Fatalln("Failed to write file:", err)
			}
		}

		c <- true
	}()

	<-c
}

func main() {
	co := colly.NewCollector(func(c *colly.Collector) {
	})

	setupCallback(co)

	provinceCount := getAreaCount(co, common.FETCH_TYPE_PROVINCE)
	cityCount := getAreaCount(co, common.FETCH_TYPE_CITY)
	districtCount := getAreaCount(co, common.FETCH_TYPE_DISTRICT)
	villageCount := getAreaCount(co, common.FETCH_TYPE_VILLAGE)

	perPage := common.MAX_DATA_PER_PAGE

	jsonWriter := writer.NewJSONWriter(common.JSON_INDENT_SIZE)
	csvWriter := writer.NewCSVWriter(common.CSV_COMMA, nil)

	csvWriter.SetHeaders([]string{"index", "name", "total_regency_and_city", "total_regency", "total_city", "total_district", "total_village", "total_island", "area_code"})
	scrapeArea(co, common.FETCH_TYPE_PROVINCE, perPage, provinceCount, scraper.ScrapeProvince, csvWriter, jsonWriter)

	csvWriter.SetHeaders([]string{"index", "name", "province", "total_district", "total_village", "area_code"})
	scrapeArea(co, common.FETCH_TYPE_CITY, perPage, cityCount, scraper.ScrapeCity, csvWriter, jsonWriter)

	csvWriter.SetHeaders([]string{"index", "name", "city", "province", "total_village", "area_code"})
	scrapeArea(co, common.FETCH_TYPE_DISTRICT, perPage, districtCount, scraper.ScrapeDistrict, csvWriter, jsonWriter)

	csvWriter.SetHeaders([]string{"index", "name", "district", "city", "province", "area_code"})
	scrapeArea(co, common.FETCH_TYPE_VILLAGE, perPage, villageCount, scraper.ScrapeVillage, csvWriter, jsonWriter)
}

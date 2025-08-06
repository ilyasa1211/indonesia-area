package main

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/ilyasa1211/indonesia-area/scraper"
	"github.com/ilyasa1211/indonesia-area/utils"
)

const (
	BASE_URL_INCOMPLETED     = "https://www.nomor.net/_kodepos.php?_i=&daerah=&jobs=&perhal=100&urut=&asc=001000&sby=010000"
	HTML_DATA_SELECTOR       = "body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(3) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(4) > tbody:nth-child(2)"
	HTML_DATA_COUNT_SELECTOR = "body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(3) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > center:nth-child(7) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > div:nth-child(1) > b:nth-child(2) > font:nth-child(1)"
	MAX_DATA_PER_PAGE        = 1000
)

type FetchType string

const (
	FETCH_TYPE_PROVINCE FetchType = "provinsi-kodepos"
	FETCH_TYPE_CITY     FetchType = "kota-kodepos"
	FETCH_TYPE_DISTRICT FetchType = "kecamatan-kodepos"
	FETCH_TYPE_VILLAGE  FetchType = "desa-kodepos"
)

func GetURL(fetchType FetchType, page uint, limit uint) *url.URL {
	if limit > MAX_DATA_PER_PAGE {
		log.Printf("Warning: limit couldn't be bigger than %d, will use %d instead.", MAX_DATA_PER_PAGE, MAX_DATA_PER_PAGE)
	}
	url, err := url.Parse(BASE_URL_INCOMPLETED)

	if err != nil {
		log.Fatalln("failed to parse url: ", err)
	}

	query := url.Query()

	query.Set("perhal", fmt.Sprintf("%d", limit))
	query.Set("_i", string(fetchType))

	if page == 1 {
		query.Set("no1", "2")

		url.RawQuery = query.Encode()
		return url
	}

	no1 := (page-2)*1000 + 1
	no2 := (page - 1) * 1000

	query.Set("no1", fmt.Sprintf("%d", no1))
	query.Set("no2", fmt.Sprintf("%d", no2))
	query.Set("kk", fmt.Sprintf("%d", page))

	url.RawQuery = query.Encode()
	return url
}

func GetDataCount(c *colly.Collector, url string, selector string) uint {
	var count uint
	done := make(chan struct{})

	fmt.Println(url)
	utils.SetProperHeader(c)

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 403 {
			log.Println("got forbidden, retrying")
			c.Visit(url) // Retry once on 403
		}
	})

	c.OnHTML(selector, func(h *colly.HTMLElement) {
		// Parse the text to extract a number
		text := strings.ReplaceAll(strings.TrimSpace(h.Text), ".", "")
		parsed, err := strconv.Atoi(text)
		if err == nil {
			count = uint(parsed)
		}
		close(done)
	})

	c.Visit(url)

	<-done // Wait for OnHTML to complete

	return count
}

func main() {
	co := colly.NewCollector(func(c *colly.Collector) {
	})

	// provinceFile := OpenFile("out/csv/provinces.csv")
	// cityFile := OpenFile("out/csv/cities.csv")
	// districtFile := OpenFile("out/csv/subdistricts.csv")
	// villageFile := OpenFile("out/csv/villages.csv")

	// defer provinceFile.Close()
	// defer cityFile.Close()
	// defer districtFile.Close()
	// defer villageFile.Close()

	// provinceWriter := csv.NewWriter(provinceFile)
	// cityWriter := csv.NewWriter(cityFile)
	// districtWriter := csv.NewWriter(districtFile)
	// villageWriter := csv.NewWriter(villageFile)

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

	// c.OnResponseHeaders(func(r *colly.Response) {
	// 	fmt.Println("Visited", r.Request.URL)
	// })

	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Println("Visited", r.Request.URL)
	// })

	provinceCount := GetDataCount(co.Clone(), GetURL(FETCH_TYPE_PROVINCE, 1, 1).String(), HTML_DATA_COUNT_SELECTOR)
	cityCount := GetDataCount(co.Clone(), GetURL(FETCH_TYPE_CITY, 1, 1).String(), HTML_DATA_COUNT_SELECTOR)
	districtCount := GetDataCount(co.Clone(), GetURL(FETCH_TYPE_DISTRICT, 1, 1).String(), HTML_DATA_COUNT_SELECTOR)
	villageCount := GetDataCount(co.Clone(), GetURL(FETCH_TYPE_VILLAGE, 1, 1).String(), HTML_DATA_COUNT_SELECTOR)

	fmt.Println([]int{
		int(provinceCount),
		int(cityCount),
		int(districtCount),
		int(villageCount),
	})

	perPage := MAX_DATA_PER_PAGE

	// 83400 / 1000 = 84 iteration

	p := math.Ceil(float64(provinceCount) / float64(perPage))
	c := math.Ceil(float64(cityCount) / float64(perPage))
	d := math.Ceil(float64(districtCount) / float64(perPage))
	v := math.Ceil(float64(villageCount) / float64(perPage))

	for i := 0; i < int(p); i++ {
		scraper.ScrapeProvince(co.Clone(), HTML_DATA_SELECTOR, i, uint(perPage))
	}
	for i := 0; i < int(c); i++ {
		scraper.ScrapeCity(co.Clone(), HTML_DATA_SELECTOR, i, uint(perPage))
	}
	for i := 0; i < int(d); i++ {
		scraper.ScrapeDistrict(co.Clone(), HTML_DATA_SELECTOR, i, uint(perPage))
	}
	for i := 0; i < int(v); i++ {
		scraper.ScrapeVillage(co.Clone(), HTML_DATA_SELECTOR, i, uint(perPage))
	}

	// c.OnHTML("", func(e *colly.HTMLElement) {
	// 	fetchType := e.Request.URL.Query().Get("_i")

	// 	switch fetchType {
	// 	case FETCH_TYPE_PROVINCE:
	// 		scraper.FetchProvince(e, provinceWriter)
	// 	case FETCH_TYPE_CITY:
	// 		FetchCity(e, cityWriter)
	// 	case FETCH_TYPE_DISTRICT:
	// 		FetchDistrict(e, districtWriter)
	// 	case FETCH_TYPE_VILLAGE:
	// 		FetchVillage(e, villageWriter)
	// 	}
	// })

}

package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gocolly/colly/v2"
	"github.com/ilyasa1211/indonesia-area/config"
	"github.com/ilyasa1211/indonesia-area/types"
)

func SetProperHeader(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		header := r.Headers

		header.Add("User-Agent", config.FAKE_USER_AGENT)
		header.Add("Accept", "text/html")
		header.Add("Accept-Language", "en-US,en;")
		header.Add("Upgrade-Insecure-Requests", "1")
		header.Add("Sec-Fetch-Dest", "document")
		header.Add("Sec-Fetch-Mode", "navigate")
		header.Add("Sec-Fetch-Site", "same-origin")
	})
}

func SetErrorHandling(c *colly.Collector) {
	c.OnError(func(r *colly.Response, err error) {
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
		log.Fatalln("Something went wrong:", r.StatusCode, r, err)
	})
}

func OpenFile(path string) *os.File {
	outDir := filepath.Dir(path)

	err := os.MkdirAll(outDir, 0755)

	if err != nil {
		log.Fatalln("failed to create directories: ", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalln("error while opening file", err)
	}

	return file
}

func GetURL(fetchType types.FetchType, page uint, limit uint) *url.URL {
	if page < 1 {
		log.Fatalln("page should not or below 0")
	}

	if limit > config.MAX_DATA_PER_PAGE {
		log.Printf("Warning: limit couldn't be bigger than %d, will use %d instead.", config.MAX_DATA_PER_PAGE, config.MAX_DATA_PER_PAGE)
	}
	url, err := url.Parse(config.BASE_URL_INCOMPLETED)

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

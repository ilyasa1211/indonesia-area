package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gocolly/colly/v2"
)

const FAKE_USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0"

func SetProperHeader(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		header := r.Headers

		header.Add("User-Agent", FAKE_USER_AGENT)
		header.Add("Accept", "text/html")
		header.Add("Accept-Language", "en-US,en;")
		header.Add("Upgrade-Insecure-Requests", "1")
		header.Add("Sec-Fetch-Dest", "document")
		header.Add("Sec-Fetch-Mode", "navigate")
		header.Add("Sec-Fetch-Site", "same-origin")
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

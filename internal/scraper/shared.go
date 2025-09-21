package scraper

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func GetDataCount(c *colly.Collector, url string, selector string) int {
	var count int
	done := make(chan struct{})

	fmt.Println(url)

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
			count = parsed
		}
		close(done)
	})

	c.Visit(url)

	<-done // Wait for OnHTML to complete

	return count
}

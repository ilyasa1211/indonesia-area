package common

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

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

func GetURL(fetchType FetchType, page int, limit int) *url.URL {
	if page < 1 {
		log.Fatalln("page should not or below 0")
	}

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

package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly/v2"
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

type CityType string

// const (
// 	CityTypeCity    CityType = "Kota"
// 	CityTypeRegency CityType = "Kab." // with dot suffix
// )

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

type DistrictData struct {
	Index        int
	Name         string
	PostalCode   []string // could be more than one
	TotalVillage int
	TotalIsland  *int
	AreaCode     string

	CityAreaCode     string
	ProvinceAreaCode string
}

type VillageData struct {
	Index      int // start from 1
	PostalCode string
	Name       string
	AreaCode   string

	DistrictAreaCode string
	CityAreaCode     string
	ProvinceAreaCode string
}

const (
	ROOT_URL        = "https://www.nomor.net/_kodepos.php?_i=provinsi-kodepos&daerah=&jobs=&perhal=60&urut=&asc=000011111&sby=000000"
	FAKE_USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0"
)

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
const (
	DISTRICT_INDEX = iota
	DISTRICT_NAME
	DISTRICT_POSTAL_CODE
	DISTRICT_TOTAL_VILLAGE
	DISTRICT_TOTAL_ISLAND
	DISTRICT_AREA_CODE
)
const (
	VILLAGE_INDEX = iota
	VILLAGE_POSTAL_CODE
	VILLAGE_NAME
	VILLAGE_AREA_CODE
)

const (
	FETCH_TYPE_PROVINCE = "provinsi-kodepos"
	FETCH_TYPE_CITY     = "kota-kodepos"
	FETCH_TYPE_DISTRICT = "kecamatan-kodepos"
	FETCH_TYPE_VILLAGE  = "desa-kodepos"
)

func FetchProvince(e *colly.HTMLElement, writer *csv.Writer) {
	provinceCount := e.DOM.Children().Length() - 1

	e.ForEach("tr", func(i int, h *colly.HTMLElement) {
		// Skip the last row and header
		if i >= provinceCount {
			return
		}

		child := h.DOM.Children()

		provinceName := child.Eq(PROVINCE_NAME).Text()
		cityUrl, exists := child.Eq(PROVINCE_NAME).Children().First().Attr("href")
		if !exists {
			log.Fatalln("City link doesn't exist for province:", provinceName)
		}

		index := child.Eq(PROVINCE_INDEX).Text()
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

		e.Request.Visit(cityUrl)
	})

	writer.Flush()
}

func FetchCity(e *colly.HTMLElement, writer *csv.Writer) {
	cityCount := e.DOM.Children().Length() - 1

	e.ForEach("tr", func(i int, h *colly.HTMLElement) {
		// skip last row
		if i >= cityCount {
			return
		}

		child := h.DOM.Children()

		name := child.Eq(CITY_NAME).Text()
		districtUrl, exists := child.Eq(CITY_NAME).Children().First().Attr("href")

		if !exists {
			log.Fatalln("link district doesnt exists", name)
		}

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

		e.Request.Visit(districtUrl)
	})

	writer.Flush()
}
func FetchDistrict(e *colly.HTMLElement, writer *csv.Writer) {
	districtCount := e.DOM.Children().Length() - 1

	e.ForEach("tr", func(i int, h *colly.HTMLElement) {
		// skip last row
		if i >= districtCount {
			return
		}

		child := h.DOM.Children()

		name := child.Eq(DISTRICT_NAME).Text()
		villageUrl, exists := child.Eq(DISTRICT_NAME).Children().First().Attr("href")

		if !exists {
			log.Fatalln("link district doesnt exists", name)
		}

		index := child.Eq(DISTRICT_INDEX).Text()
		totalVillage := strings.ReplaceAll(child.Eq(DISTRICT_TOTAL_VILLAGE).Text(), ".", "")
		totalIsland := strings.ReplaceAll(child.Eq(DISTRICT_TOTAL_ISLAND).Text(), ".", "")
		postalCode := strings.Split(child.Eq(DISTRICT_POSTAL_CODE).Text(), " - ")
		areaCode := child.Eq(DISTRICT_AREA_CODE).Text()

		err := writer.Write([]string{
			(index),
			name,
			(totalVillage),
			(totalIsland),
			areaCode,
			strings.Join(postalCode, "-"),
		})

		if err != nil {
			log.Fatalln("Failed to write CSV row:", err)
		}

		e.Request.Visit(villageUrl)
	})

	writer.Flush()
}
func FetchVillage(e *colly.HTMLElement, writer *csv.Writer) {
	villageCount := e.DOM.Children().Length() - 1

	e.ForEach("tr", func(i int, h *colly.HTMLElement) {
		// skip last row
		if i >= villageCount {
			return
		}
		child := h.DOM.Children()

		name := child.Eq(VILLAGE_NAME).Text()
		index := child.Eq(VILLAGE_INDEX).Text()
		postalCode := child.Eq(VILLAGE_POSTAL_CODE).Text()
		areaCode := child.Eq(VILLAGE_AREA_CODE).Text()

		err := writer.Write([]string{
			index,
			name,
			postalCode,
			areaCode,
		})

		if err != nil {
			log.Fatalln("Failed to write CSV row:", err)
		}
	})

	writer.Flush()
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

func main() {
	c := colly.NewCollector(func(c *colly.Collector) {
	})

	provinceFile := OpenFile("out/csv/provinces.csv")
	cityFile := OpenFile("out/csv/cities.csv")
	districtFile := OpenFile("out/csv/subdistricts.csv")
	villageFile := OpenFile("out/csv/villages.csv")

	defer provinceFile.Close()
	defer cityFile.Close()
	defer districtFile.Close()
	defer villageFile.Close()

	provinceWriter := csv.NewWriter(provinceFile)
	cityWriter := csv.NewWriter(cityFile)
	districtWriter := csv.NewWriter(districtFile)
	villageWriter := csv.NewWriter(villageFile)

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

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong:", r.StatusCode, r, err)
	})

	// c.OnResponseHeaders(func(r *colly.Response) {
	// 	fmt.Println("Visited", r.Request.URL)
	// })

	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Println("Visited", r.Request.URL)
	// })

	// c.OnResponse(func(r *colly.Response) {

	// })

	c.OnHTML("body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(3) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(4) > tbody:nth-child(2)", func(e *colly.HTMLElement) {
		fetchType := e.Request.URL.Query().Get("_i")
		switch fetchType {
		case FETCH_TYPE_PROVINCE:
			FetchProvince(e, provinceWriter)
		case FETCH_TYPE_CITY:
			FetchCity(e, cityWriter)
		case FETCH_TYPE_DISTRICT:
			FetchDistrict(e, districtWriter)
		case FETCH_TYPE_VILLAGE:
			FetchVillage(e, villageWriter)
		default:
			FetchProvince(e, provinceWriter)
		}
	})

	c.Visit(ROOT_URL)
}

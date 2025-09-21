package common

const (
	// This url is not completed, and couldn't be used without further modification.
	BASE_URL_INCOMPLETED = "https://www.nomor.net/_kodepos.php?_i=&daerah=&jobs=&perhal=100&urut=&asc=001000&sby=010000"

	// CSS selector to grab the table containing data
	HTML_DATA_SELECTOR = "body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(3) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(4) > tbody:nth-child(2)"

	// CSS selector to grab the total count of the data
	HTML_DATA_COUNT_SELECTOR = "body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(3) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > center:nth-child(7) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > div:nth-child(1) > b:nth-child(2) > font:nth-child(1)"

	// Used to tricks bot protection
	FAKE_USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0"

	MAX_DATA_PER_PAGE = 1000
	OUT_DIR           = "out/"
	JSON_INDENT_SIZE  = 2
	CSV_COMMA         = ','
)

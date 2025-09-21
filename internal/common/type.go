package common

type FetchType string

const (
	FETCH_TYPE_PROVINCE FetchType = "provinsi-kodepos"
	FETCH_TYPE_CITY     FetchType = "kota-kodepos"
	FETCH_TYPE_DISTRICT FetchType = "kecamatan-kodepos"
	FETCH_TYPE_VILLAGE  FetchType = "desa-kodepos"
)

var FetchTypeToName = map[FetchType]string{
	FETCH_TYPE_PROVINCE: "province",
	FETCH_TYPE_CITY:     "city",
	FETCH_TYPE_DISTRICT: "district",
	FETCH_TYPE_VILLAGE:  "village",
}

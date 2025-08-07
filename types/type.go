package types

type FetchType string

const (
	FETCH_TYPE_PROVINCE FetchType = "provinsi-kodepos"
	FETCH_TYPE_CITY     FetchType = "kota-kodepos"
	FETCH_TYPE_DISTRICT FetchType = "kecamatan-kodepos"
	FETCH_TYPE_VILLAGE  FetchType = "desa-kodepos"
)

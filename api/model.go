package api

type Thing struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Alias string `json:"alias"`
	Enabled bool   `json:"enabled"`
	LastSeen int32  `json:"last_seen"`
	StoreInfluxDb bool `json:"store_influxdb"`
	StoreMysqlDb bool `json:"store_mysqldb"`
}

type Org struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

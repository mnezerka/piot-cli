package api

type Thing struct {
	Id               string `json:"id" csv:"id"`
	Name             string `json:"name" csv:"name"`
	Type             string `json:"type" csv:"type"`
	Alias            string `json:"alias" csv:"alias"`
	Enabled          bool   `json:"enabled" csv:"enabled"`
	LastSeen         int32  `json:"last_seen" csv:"last_seen"`
	LastSeenInterval int32  `json:"last_seen_interval" csv:"last_seen_interval"`
	StoreInfluxDb    bool   `json:"store_influxdb" csv:"store_influxdb"`
	StoreMysqlDb     bool   `json:"store_mysqldb" csv:"store_mysqldb"`
}

type Org struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	InfluxDb string `json:"influxdb"`
}

type GqlLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type GqlError struct {
	Message   string        `json:"message"`
	Locations []GqlLocation `json:"locations"`
	Path      []interface{} `json:"path"`
}

type GqlResponse struct {
	Errors map[string]string `json:"errors"`
}

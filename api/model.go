package api

type Thing struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Alias            string `json:"alias"`
	Enabled          bool   `json:"enabled"`
	LastSeen         int32  `json:"last_seen"`
	LastSeenInterval int32  `json:"last_seen_interval"`
	StoreInfluxDb    bool   `json:"store_influxdb"`
	StoreMysqlDb     bool   `json:"store_mysqldb"`
}

type Org struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserProfile struct {
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	OrgId   string `json:"org_id"`
	Orgs    []Org  `json:"orgs"`
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

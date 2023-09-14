package model

import "time"

type Response struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		Geo GeoData `json:"geo"`
	} `json:"data"`
}

type GeoData struct {
	Host          string      `json:"host"`
	IP            string      `json:"ip"`
	RDNS          string      `json:"rdns"`
	ISP           string      `json:"isp"`
	CountryName   string      `json:"country_name"`
	CountryCode   string      `json:"country_code"`
	RegionName    string      `json:"region_name"`
	RegionCode    string      `json:"region_code"`
	City          string      `json:"city"`
	PostalCode    string      `json:"postal_code"`
	ContinentName string      `json:"continent_name"`
	ContinentCode string      `json:"continent_code"`
	Latitude      interface{} `json:"latitude"`
	Longitude     interface{} `json:"longitude"`
	MetroCode     string      `json:"metro_code"`
	Timezone      string      `json:"timezone"`
	Datetime      interface{} `json:"datetime"`
}

type LookupResponse struct {
	City        string `json:"city"`
	RegionName  string `json:"region"`
	CountryName string `json:"country"`
}

type LookupRequest struct {
	IpAddress    string    `json:"ip_address"`
	LookupTime   time.Time `json:"lookup_time"`
	LookupStatus bool      `json:"lookup_status"`
}

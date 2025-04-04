package dto

import "time"

type AuditLog struct {
	Models
	UserID    string    `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	Device    string    `json:"device"`
	Country   string    `json:"country"`
	Browser   string    `json:"browser"`
	Timestamp time.Time `json:"timestamp"`
}

type IpAPIResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

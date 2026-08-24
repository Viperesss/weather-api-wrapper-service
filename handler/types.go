package handler

type WeatherReport struct {
	Address string `json:"address"`
	Days    []Day  `json:"days"`
}

type Day struct {
	Datetime   string  `json:"datetime"`
	Temp       float64 `json:"temp"`
	Conditions string  `json:"conditions"`
}

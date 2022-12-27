package structs

type CommonConfig struct {
	AppVersion string   `json:"app_version"`
	UniqueIds  []string `json:"unique_ids"`
	ValuesMin  float64  `json:"values_min"`
	ValuesMax  float64  `json:"values_max"`
}

type Message struct {
	Id    string  `json:"id"`
	Value float64 `json:"value"`
}

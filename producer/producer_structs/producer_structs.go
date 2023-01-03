package producer_structs

type ProducerConfig struct {
	AppName         string   `json:"app_name"`
	MessageInterval int64    `json:"message_interval"`
	UniqueIds       []string `json:"unique_ids"`
	ValuesMin       float64  `json:"values_min"`
	ValuesMax       float64  `json:"values_max"`
}

type Message struct {
	Id    string  `json:"id"`
	Value float64 `json:"value"`
}

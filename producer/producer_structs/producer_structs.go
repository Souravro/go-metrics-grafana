package producer_structs

type ProducerConfig struct {
	AppName         string `json:"app_name"`
	MessageInterval int64  `json:"message_interval"`
}

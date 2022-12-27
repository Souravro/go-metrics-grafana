package consumer_structs

type ConsumerConfig struct {
	AppName       string `json:"app_name"`
	BadgerTempDir string `json:"badger_temp_dir"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

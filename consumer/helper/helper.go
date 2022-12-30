package helper

import (
	"consumer/consumer_structs"
	"encoding/json"
	"log"
	"os"
)

func LoadConsumerConfiguration(file string) (cConfig consumer_structs.ConsumerConfig) {
	configFile, err := os.Open(file)
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Printf("Error in closing consumer config file. Error: [%v]", err)
		}
	}(configFile)
	if err != nil {
		log.Printf("Error in opening consumer config file. Error: [%v]", err)
		return
	}
	jsonParser := json.NewDecoder(configFile)
	er := jsonParser.Decode(&cConfig)
	if er != nil {
		log.Printf("Error in json decoding consumer config. Error: [%v]", er)
		return
	}

	return cConfig
}

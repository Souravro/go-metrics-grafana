package helper

import (
	"encoding/json"
	"log"
	"os"
	"producer/producer_structs"
)

func LoadProducerConfiguration(file string) (pConfig producer_structs.ProducerConfig) {
	configFile, err := os.Open(file)
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Printf("Error in closing producer config file. Error: [%v]", err)
		}
	}(configFile)
	if err != nil {
		log.Printf("Error in opening producer config file. Error: [%v]", err)
		return
	}
	jsonParser := json.NewDecoder(configFile)
	er := jsonParser.Decode(&pConfig)
	if er != nil {
		log.Printf("Error in json decoding producer config. Error: [%v]", er)
		return
	}

	return pConfig
}

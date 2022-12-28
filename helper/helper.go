package helper

import (
	"encoding/json"
	"log"
	"os"
	"producer/consumer/consumer_structs"
	"producer/dashboard/dashboard_structs"
	"producer/producer/producer_structs"
	"producer/structs"
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

func LoadCommonConfiguration(file string) (common structs.CommonConfig) {
	configFile, err := os.Open(file)
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Printf("Error in closing common config file. Error: [%v]", err)
		}
	}(configFile)
	if err != nil {
		log.Printf("Error in opening common config file. Error: [%v]", err)
		return
	}
	jsonParser := json.NewDecoder(configFile)
	er := jsonParser.Decode(&common)
	if er != nil {
		log.Printf("Error in json decoding common config. Error: [%v]", er)
		return
	}

	return common
}

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

func LoadDashboardConfiguration(file string) (dConfig dashboard_structs.DashboardConfig) {
	configFile, err := os.Open(file)
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Printf("Error in closing dashboard config file. Error: [%v]", err)
		}
	}(configFile)
	if err != nil {
		log.Printf("Error in opening dashboard config file. Error: [%v]", err)
		return
	}
	jsonParser := json.NewDecoder(configFile)
	er := jsonParser.Decode(&dConfig)
	if er != nil {
		log.Printf("Error in json decoding dashboard config. Error: [%v]", er)
		return
	}

	return dConfig
}

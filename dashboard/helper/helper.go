package helper

import (
	"dashboard/dashboard_structs"
	"encoding/json"
	"log"
	"os"
)

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

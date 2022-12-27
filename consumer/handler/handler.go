package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"producer/consumer/consumer_structs"
	"producer/consumer/store"
	"producer/helper"
	"producer/structs"
)

var (
	commonConfig   structs.CommonConfig
	consumerConfig consumer_structs.ConsumerConfig
	storageSvc     store.StorageService
)

func GetValueForId(w http.ResponseWriter, r *http.Request) {
	log.Printf("In consumer.GetValueForId handler..")
	id := r.URL.Query().Get("id")
	log.Printf("ID in request: [%v]", id)

	// set common header
	w.Header().Set("content-type", "application/json")

	switch r.Method {
	case http.MethodGet:
		data, err := getValue(id)
		if err != nil {
			log.Printf("consumer.GetValueForId Error: [%v]", err)
			resp := consumer_structs.Response{
				Status:  "Failure",
				Message: err.Error(),
				Data:    nil,
			}
			respByte, _ := json.Marshal(resp)
			_, er := w.Write(respByte)
			if er != nil {
				return
			}

			return
		}

		log.Printf("consumer.GetValueForId Response: [%v]", data)
		resp := consumer_structs.Response{
			Status:  "Success",
			Message: "Data fetched successfully.",
			Data:    data,
		}
		respByte, _ := json.Marshal(resp)
		_, er := w.Write(respByte)
		if er != nil {
			return
		}
	default:
		resp := consumer_structs.Response{
			Status:  "Failure",
			Message: "Method not allowed",
			Data:    nil,
		}
		respByte, _ := json.Marshal(resp)
		_, er := w.Write(respByte)
		if er != nil {
			return
		}
	}
}

func getValue(id string) (structs.Message, error) {
	commonConfig = helper.LoadCommonConfiguration("config/common.json")
	consumerConfig = helper.LoadConsumerConfiguration("consumer/config/config.json")

	storageSvc = store.StorageService{
		ConsumerConfig: consumerConfig,
		CommonConfig:   commonConfig,
	}
	respMessage, err := storageSvc.GetValue(id)
	if err != nil {
		log.Printf("consumer.GetValueForId Error in getting value for id: [%v]. Error: [%v]", id, err)
		return structs.Message{}, err
	}

	return respMessage, nil
}

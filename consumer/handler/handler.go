package handler

import (
	"consumer/consumer_structs"
	"consumer/store"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"time"
)

var (
	storageSvc   *store.StorageService
	IdApiSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "consumer",
		Name:      "id_api_latency",
		Help:      "Latency for /getValueForId api, initiating from consumer_service",
	}, []string{"id"})
)

func GetValueForId(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	log.Printf("In consumer.GetValueForId handler..")
	id := r.URL.Query().Get("id")
	log.Printf("ID in request: [%v]", id)

	defer func() {
		elapsedTime := time.Since(startTime).Seconds()
		IdApiSummary.WithLabelValues(id).Observe(elapsedTime)
	}()

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

func getValue(id string) (consumer_structs.Message, error) {
	storageSvc = store.GetService()
	respMessage, err := storageSvc.GetValue(id)
	if err != nil {
		log.Printf("consumer.GetValueForId Error in getting value for id: [%v]. Error: [%v]", id, err)
		return consumer_structs.Message{}, err
	}

	return respMessage, nil
}

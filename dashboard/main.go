package main

import (
	"context"
	"dashboard/dashboard_structs"
	"dashboard/helper"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	GetRecordAPI = "getValueForId"
)

var (
	dashboardConfig dashboard_structs.DashboardConfig
	idApiSummary    = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "dashboard",
		Name:      "id_api_latency",
		Help:      "Latency for " + GetRecordAPI + " api, initiating from dashboard_service",
	}, []string{"id"})
)

func registerPrometheusMetrics() {
	prometheus.MustRegister(idApiSummary)
}

func main() {
	log.Println("Starting dashboard application...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dashboardConfig = helper.LoadDashboardConfiguration(os.Getenv("APP_HOME") + "/config/config.json")
	log.Printf("Dashboard config: [%v]", dashboardConfig)

	// Prometheus metric
	http.Handle("/metrics", promhttp.Handler())

	// Start http server
	fmt.Printf("Starting server at port 2121...\n")
	go func() {
		if err := http.ListenAndServe(":2121", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Register Prometheus custom metrics
	registerPrometheusMetrics()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := getRecord(); err != nil {
					log.Printf("Dashboard. Error while getting record. Error: [%v]", err)
					continue
				}
				time.Sleep(time.Duration(dashboardConfig.RequestInterval) * time.Millisecond)
			}
		}
	}()

	wg.Wait()
}

func getRecord() error {
	startTime := time.Now()

	idValue := dashboardConfig.UniqueIds[rand.Intn(len(dashboardConfig.UniqueIds))]
	queryParam := "id=" + idValue
	requestURL := dashboardConfig.DataHost + "/" + GetRecordAPI + "?" + queryParam
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Printf("Dashboard. Error in creating API request object hit to data host: [%v], Error: [%v]", dashboardConfig.DataHost, err)
		return err
	}

	res, er := http.DefaultClient.Do(req)
	if er != nil {
		log.Printf("Dashboard. Error while making API request to data host: [%v], Error: [%v]", dashboardConfig.DataHost, er)
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Dashboard: could not read response body. Error: [%v]", err)
		return err
	}

	log.Printf("Dashboard. Response: [%v]", string(resBody))

	elapsedTime := time.Since(startTime).Seconds()
	idApiSummary.WithLabelValues(idValue).Observe(elapsedTime)

	return nil
}

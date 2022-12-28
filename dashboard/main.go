package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"producer/dashboard/dashboard_structs"
	"producer/helper"
	"producer/structs"
	"sync"
	"time"
)

var (
	dashboardConfig dashboard_structs.DashboardConfig
	commonConfig    structs.CommonConfig
)

const (
	GetRecordAPI = "getValueForId"
)

func main() {
	log.Println("Starting dashboard application...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dashboardConfig = helper.LoadDashboardConfiguration("dashboard/config/config.json")
	log.Printf("Dashboard config: [%v]", dashboardConfig)

	commonConfig = helper.LoadCommonConfiguration("config/common.json")
	log.Printf("Common config: [%v]", commonConfig)

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
	queryParam := "id=" + commonConfig.UniqueIds[rand.Intn(len(commonConfig.UniqueIds))]
	requestURL := dashboardConfig.DataHost + "/" + GetRecordAPI + "?" + queryParam
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Printf("Dashboard. Error in creating API request object hit to data host: [%v], Error: [%v]", dashboardConfig.DataHost, err)
		return err
	}

	res, er := http.DefaultClient.Do(req)
	if er != nil {
		log.Printf("Dashboard. Error while making API request to data host: [%v], Error: [%v]", dashboardConfig.DataHost, err)
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Dashboard: could not read response body. Error: [%v]", err)
		return err
	}

	log.Printf("Dashboard. Response: [%v]", string(resBody))

	return nil
}

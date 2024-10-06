package main

import (
	"fmt"
	"log"
	"net/http"
	"web-monitor/database"
	webMonitorMetrics "web-monitor/web-monitor-metrics"
	webUrl "web-monitor/web-url"

	"github.com/spf13/viper"
)

func main() {
	database.InitDB()
	webUrl.CreateWebUrlTable()
	webMonitorMetrics.CreateWebMonitorMetricTable()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	monitoringProcessInterval := viper.GetString("web-monitor-scheduler.schedule")
	fmt.Println("Web monitor configuration")
	fmt.Println("Monitoring Process Schedule:", monitoringProcessInterval)
	go StartMonitoringScheduler(monitoringProcessInterval)
	mux := setupEndpoints()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("\nServer is running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func setupEndpoints() *http.ServeMux {
	mux := http.NewServeMux()
	controller := webUrl.WebController{}

	mux.HandleFunc("/web-url", controller.AddWebUrlForMonitoring)

	return mux
}

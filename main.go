package main

import (
	"net/http"
	"web-monitor/database"
	webMonitorMetrics "web-monitor/web-monitor-metrics"
	webUrl "web-monitor/web-url"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	database.InitDB(logger)
	webUrl.CreateWebUrlTable()
	webMonitorMetrics.CreateWebMonitorMetricTable()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("Fatal error config file", zap.Error(err))
	}
	monitoringProcessInterval := viper.GetString("web-monitor-scheduler.schedule")
	logger.Info("Monitoring process interval", zap.String("interval", monitoringProcessInterval))

	controller := webUrl.NewWebController(logger, database.GetDB())
	mux := setupEndpoints(logger, controller)
	go StartMonitoringScheduler(monitoringProcessInterval, logger, &controller.WebUrlRepository)
	serverPort := viper.GetString("server.port")
	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}
	logger.Info("Server started", zap.String("port", serverPort))
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

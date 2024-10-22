package main

import (
	"fmt"
	webUrl "web-monitor/web-url"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var scheduledIntervals []int
var logger *zap.Logger
var webUrlRepository *webUrl.WebUrlRepository

func StartMonitoringScheduler(monitoringProcessInterval string, _logger *zap.Logger, _webUrlRepository *webUrl.WebUrlRepository) {
	logger = _logger
	webUrlRepository = _webUrlRepository
	logger.Info("Starting monitoring scheduler")
	monitoringProcessCronJob := cron.New()
	monitoringProcessCronJob.AddFunc(monitoringProcessInterval, startMonitoring)
	monitoringProcessCronJob.Start()
}

func startMonitoring() {
	distinctIntervals, err := webUrlRepository.GetAllDistinctWebUrlIntervals()
	if err != nil {
		logger.Error("Error while fetching distinct intervals", zap.Error(err))
		return
	}
	unscheduledIntervals := getUnscheduledIntervals(distinctIntervals)
	logger.Info("Unscheduled intervals", zap.Any("unscheduledIntervals", unscheduledIntervals))
	scheduleMonitoringCronJobs(unscheduledIntervals)
}

func getUnscheduledIntervals(distinctIntervals []int) []int {
	scheduledIntervalsMap := make(map[int]bool)
	for _, value := range scheduledIntervals {
		scheduledIntervalsMap[value] = true
	}

	var result []int
	for _, value := range distinctIntervals {
		if !scheduledIntervalsMap[value] {
			result = append(result, value)
		}
	}

	return result
}

func scheduleMonitoringCronJobs(intervals []int) {
	for _, interval := range intervals {
		logger.Info("Scheduling monitoring cron job", zap.Int("interval", interval))
		webUrlsMonitoringCronJob := cron.New()
		webUrlsMonitoringCronJob.AddFunc(fmt.Sprintf("*/%d * * * *", interval), func() {
			logger.Info("Fetching web urls for interval", zap.Int("interval", interval))
			webUrls, err := fetchWebUrlsByInterval(interval)
			if err != nil {
				logger.Error("Error while fetching web urls, Retry in the next interval", zap.Error(err))
				return
			}
			CallWebUrls(webUrls)
		})
		webUrlsMonitoringCronJob.Start()
	}
	scheduledIntervals = append(scheduledIntervals, intervals...)
}

func fetchWebUrlsByInterval(interval int) (webUrls []webUrl.WebUrl, err error) {
	webUrls, err = webUrlRepository.GetWebUrlsByInterval(interval)
	if err != nil {
		logger.Error("Error while fetching web urls", zap.Error(err))
		return nil, err
	}
	return webUrls, nil
}

package main

import (
	"fmt"
	webUrl "web-monitor/web-url"

	"github.com/robfig/cron/v3"
)

var scheduledIntervals []int

func StartMonitoringScheduler(monitoringProcessInterval string) {
	fmt.Println("Starting web monitor scheduler")
	monitoringProcessCronJob := cron.New()
	monitoringProcessCronJob.AddFunc(monitoringProcessInterval, startMonitoring)
	monitoringProcessCronJob.Start()
}

func startMonitoring() {
	distinctIntervals, err := webUrl.GetAllDistinctWebUrlIntervals()
	if err != nil {
		fmt.Println("Error while fetching web urls", err)
		return
	}
	unscheduledIntervals := getUnscheduledIntervals(distinctIntervals)
	fmt.Println("Starting monitoring of unscheduled intervals", unscheduledIntervals)
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
		fmt.Println("Scheduling cron job for interval", interval)
		webUrlsMonitoringCronJob := cron.New()
		webUrlsMonitoringCronJob.AddFunc(fmt.Sprintf("*/%d * * * *", interval), func() {
			fmt.Println("Calling web urls")
			webUrls, err := fetchWebUrlsByInterval(interval)
			if err != nil {
				fmt.Println("Error while fetching web urls, Retry in the next interval", err)
				return
			}
			CallWebUrls(webUrls)
		})
		webUrlsMonitoringCronJob.Start()
	}
	scheduledIntervals = append(scheduledIntervals, intervals...)
}

func fetchWebUrlsByInterval(interval int) (webUrls []webUrl.WebUrl, err error) {
	webUrls, err = webUrl.GetWebUrlsByInterval(interval)
	if err != nil {
		fmt.Println("Error while fetching web urls", err)
		return nil, err
	}
	return webUrls, nil
}

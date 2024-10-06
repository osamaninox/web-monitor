package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"
	webMonitorMetrics "web-monitor/web-monitor-metrics"
	webUrl "web-monitor/web-url"
)

func CallWebUrls(webUrls []webUrl.WebUrl) {
	var wg sync.WaitGroup
	for _, webUrlToMonitor := range webUrls {
		wg.Add(1)
		go func(url webUrl.WebUrl) {
			defer wg.Done()
			monitorWebUrl(url)
		}(webUrlToMonitor)
	}
	wg.Wait()
}

func monitorWebUrl(webUrl webUrl.WebUrl) {
	fmt.Println("Monitoring web url", webUrl.Url)
	start := time.Now()
	httpResponse, err := http.Get(webUrl.Url)
	if err != nil {
		fmt.Println("Error while calling web url", err)
		recordMetrics(webUrl, 500, 0, false)
		return
	}
	responseTime := time.Since(start)
	fmt.Println("Time taken to call web url", webUrl, responseTime)
	fmt.Println("Response status code", httpResponse.StatusCode)
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		fmt.Println("Error while reading response body", err)
		recordMetrics(webUrl, httpResponse.StatusCode, responseTime, false)
		return
	}
	if webUrl.RegexPattern == nil {
		recordMetrics(webUrl, httpResponse.StatusCode, responseTime, false)
		return
	}
	isRegexPatternMatched := isRegexPatternMatched(*webUrl.RegexPattern, string(responseBody))
	fmt.Println("Is regex pattern matched", isRegexPatternMatched)
	recordMetrics(webUrl, httpResponse.StatusCode, responseTime, isRegexPatternMatched)
}

func isRegexPatternMatched(regexPattern string, body string) bool {
	matched, err := regexp.MatchString(regexPattern, body)
	if err != nil {
		fmt.Println("Invalid regex pattern", err)
		return false
	}
	return matched
}

func recordMetrics(webUrl webUrl.WebUrl, responseStatus int, responseTime time.Duration, isRegexPatternMatched bool) {
	fmt.Println("Recording metrics for web url", webUrl)
	result := webMonitorMetrics.CreateWebMonitorMetric(&webMonitorMetrics.WebMonitorMetric{
		Url:                   webUrl.Url,
		ResponseTime:          int(responseTime.Milliseconds()),
		ResponseStatus:        responseStatus,
		IsRegexPatternMatched: isRegexPatternMatched,
	})
	if result.Error != nil {
		fmt.Println("Error while recording metrics of web url", webUrl.Url, result.Error)
	}
}

package main

import (
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"
	webMonitorMetrics "web-monitor/web-monitor-metrics"
	webUrl "web-monitor/web-url"

	"go.uber.org/zap"
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
	logger.Info("Monitoring web url", zap.String("url", webUrl.Url))
	start := time.Now()
	httpResponse, err := http.Get(webUrl.Url)
	if err != nil {
		logger.Error("Error while calling web url", zap.Error(err))
		recordMetrics(webUrl, 500, 0, false)
		return
	}
	responseTime := time.Since(start)
	logger.Info("Time taken to call web url", zap.String("url", webUrl.Url), zap.Duration("responseTime", responseTime))
	logger.Info("Response status code", zap.Int("statusCode", httpResponse.StatusCode))
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		logger.Error("Error while reading response body", zap.Error(err))
		recordMetrics(webUrl, httpResponse.StatusCode, responseTime, false)
		return
	}
	if webUrl.RegexPattern == nil {
		recordMetrics(webUrl, httpResponse.StatusCode, responseTime, false)
		return
	}
	isRegexPatternMatched := isRegexPatternMatched(*webUrl.RegexPattern, string(responseBody))
	recordMetrics(webUrl, httpResponse.StatusCode, responseTime, isRegexPatternMatched)
}

func isRegexPatternMatched(regexPattern string, body string) bool {
	matched, err := regexp.MatchString(regexPattern, body)
	if err != nil {
		logger.Error("Error while matching regex pattern", zap.Error(err))
		return false
	}
	return matched
}

func recordMetrics(webUrl webUrl.WebUrl, responseStatus int, responseTime time.Duration, isRegexPatternMatched bool) {
	logger.Info("Recording metrics", zap.String("url", webUrl.Url), zap.Int("responseStatus", responseStatus), zap.Duration("responseTime", responseTime), zap.Bool("isRegexPatternMatched", isRegexPatternMatched))
	result := webMonitorMetrics.CreateWebMonitorMetric(&webMonitorMetrics.WebMonitorMetric{
		Url:                   webUrl.Url,
		ResponseTime:          int(responseTime.Milliseconds()),
		ResponseStatus:        responseStatus,
		IsRegexPatternMatched: isRegexPatternMatched,
	})
	if result.Error != nil {
		logger.Error("Error while recording metrics", zap.Error(result.Error))
	}
}

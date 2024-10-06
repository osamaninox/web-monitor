package main

import (
	"strings"
	"testing"
	"web-monitor/database"

	webMonitorMetrics "web-monitor/web-monitor-metrics"
	webUrl "web-monitor/web-url"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestWebUrlMonitorMetric(t *testing.T) {
	database.InitDB()

	var webUrlOne = "http://serverone.com"
	var webUrlTwo = "http://servertwo.com"

	defer gock.Off() // Flush pending mocks after test execution

	gock.New(webUrlOne).
		Get("").
		Reply(200).
		JSON("Hello, world!")

	gock.New(webUrlTwo).
		Get("").
		Reply(200).
		Body(strings.NewReader("Hello, world!"))

	t.Run("Should successfully record metrics", func(t *testing.T) {
		// Arrange
		var regexPattern = "Hello, world!"
		webUrls := []webUrl.WebUrl{
			{Url: webUrlOne, Interval: 5},
			{Url: webUrlTwo, Interval: 5, RegexPattern: &regexPattern},
		}
		// Act
		CallWebUrls(webUrls)
		// Assert
		webUrlOneMetric, err := webMonitorMetrics.GetWebMonitorMetrics(webUrlOne)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(webUrlOneMetric))
		assert.Equal(t, false, webUrlOneMetric[0].IsRegexPatternMatched)

		webUrlTwoMetric, err := webMonitorMetrics.GetWebMonitorMetrics(webUrlTwo)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(webUrlTwoMetric))
		assert.Equal(t, true, webUrlTwoMetric[0].IsRegexPatternMatched)

		// Clean up
		webMonitorMetrics.DeleteWebMonitorMetrics(webUrlOne)
		webMonitorMetrics.DeleteWebMonitorMetrics(webUrlTwo)
	})
}

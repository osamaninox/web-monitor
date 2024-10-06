package web_monitor_metrics

import (
	"web-monitor/database"

	"gorm.io/gorm"
)

func CreateWebMonitorMetric(webMonitorMetric *WebMonitorMetric) (result *gorm.DB) {
	return database.GetDB().Create(webMonitorMetric)
}

func GetWebMonitorMetrics(webUrl string) (webMonitorMetrics []WebMonitorMetric, err error) {
	result := database.GetDB().Where("url = ?", webUrl).Find(&webMonitorMetrics)
	if result.Error != nil {
		return nil, result.Error
	}
	return webMonitorMetrics, nil
}

func DeleteWebMonitorMetrics(webUrl string) (result *gorm.DB) {
	return database.GetDB().Where("url = ?", webUrl).Delete(&WebMonitorMetric{})
}

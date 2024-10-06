package web_monitor_metrics

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type WebMonitorMetric struct {
	gorm.Model
	Url                   string
	IsRegexPatternMatched bool
	ResponseTime          int
	ResponseStatus        int
}

func CreateWebMonitorMetricTable() {
	dsn := "host=postgresdb port=5432 user=postgres dbname=web-monitor-db sslmode=disable password=secretpassword"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database " + err.Error())
	}
	postgresDb, err := db.DB()
	if err != nil {
		panic("failed to connect database " + err.Error())
	}
	defer postgresDb.Close()

	// Migrate the schema
	db.AutoMigrate(&WebMonitorMetric{})
}

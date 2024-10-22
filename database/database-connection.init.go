package database

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(logger *zap.Logger) {
	dsn := "host=postgresdb port=5432 user=postgres dbname=web-monitor-db sslmode=disable password=secretpassword"
	for i := 0; i < 3; i++ {
		postgresDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			db = postgresDb
			break
		}
		logger.Error("Error while connecting to database", zap.Error(err))
		time.Sleep(5 * time.Second)
	}
}

func GetDB() *gorm.DB {
	return db
}

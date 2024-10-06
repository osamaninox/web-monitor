package web_url

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type WebUrl struct {
	gorm.Model
	Url          string  `gorm:"not null" json:"url"`
	RegexPattern *string `gorm:"default:NULL" json:"regexPattern"`
	Interval     int     `gorm:"not null" json:"interval"`
}

func CreateWebUrlTable() {
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
	db.AutoMigrate(&WebUrl{})
}

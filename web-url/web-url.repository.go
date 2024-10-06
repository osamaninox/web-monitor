package web_url

import (
	"fmt"
	"web-monitor/database"

	"gorm.io/gorm"
)

func CreateWebUrl(webUrl *WebUrl) (result *gorm.DB) {
	return database.GetDB().Create(webUrl)
}

func GetWebUrls() (groupedWebUrls map[int][]WebUrl, err error) {
	var webUrls []WebUrl
	result := database.GetDB().Find(&webUrls)
	if result.Error != nil {
		fmt.Println("Error while fetching web urls", result.Error)
		return nil, result.Error
	}
	groupedWebUrls = make(map[int][]WebUrl)
	for _, webUrl := range webUrls {
		groupedWebUrls[webUrl.Interval] = append(groupedWebUrls[webUrl.Interval], webUrl)
	}
	fmt.Println("Grouped web urls", groupedWebUrls)
	return groupedWebUrls, nil
}

func GetAllDistinctWebUrlIntervals() (intervals []int, err error) {
	result := database.GetDB().Model(&WebUrl{}).Select("DISTINCT interval").Find(&intervals)
	if result.Error != nil {
		fmt.Println("Error while fetching intervals", result.Error)
		return nil, result.Error
	}
	fmt.Println("Distinct intervals", intervals)
	return intervals, nil
}

func GetWebUrlsByInterval(interval int) (webUrls []WebUrl, err error) {
	result := database.GetDB().Where("interval = ?", interval).Find(&webUrls)
	if result.Error != nil {
		fmt.Println("Error while fetching web urls", result.Error)
		return nil, result.Error
	}
	fmt.Println("Web urls for interval", interval, webUrls)
	return webUrls, nil
}

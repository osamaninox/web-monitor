package web_url

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WebUrlRepository struct {
	database *gorm.DB
	logger   *zap.Logger
}

func (webUrlRepository *WebUrlRepository) CreateWebUrl(webUrl *WebUrl) (result *gorm.DB) {
	return webUrlRepository.database.Create(webUrl)
}

func (webUrlRepository *WebUrlRepository) GetWebUrls() (groupedWebUrls map[int][]WebUrl, err error) {
	var webUrls []WebUrl
	result := webUrlRepository.database.Find(&webUrls)
	if result.Error != nil {
		webUrlRepository.logger.Error("Error while fetching web urls", zap.Error(result.Error))
		return nil, result.Error
	}
	groupedWebUrls = make(map[int][]WebUrl)
	for _, webUrl := range webUrls {
		groupedWebUrls[webUrl.Interval] = append(groupedWebUrls[webUrl.Interval], webUrl)
	}
	webUrlRepository.logger.Info("Grouped web urls", zap.Any("groupedWebUrls", groupedWebUrls))
	return groupedWebUrls, nil
}

func (webUrlRepository *WebUrlRepository) GetAllDistinctWebUrlIntervals() (intervals []int, err error) {
	result := webUrlRepository.database.Model(&WebUrl{}).Select("DISTINCT interval").Find(&intervals)
	if result.Error != nil {
		webUrlRepository.logger.Error("Error while fetching intervals", zap.Error(result.Error))
		return nil, result.Error
	}
	webUrlRepository.logger.Info("Distinct intervals", zap.Any("intervals", intervals))
	return intervals, nil
}

func (webUrlRepository *WebUrlRepository) GetWebUrlsByInterval(interval int) (webUrls []WebUrl, err error) {
	result := webUrlRepository.database.Where("interval = ?", interval).Find(&webUrls)
	if result.Error != nil {
		webUrlRepository.logger.Error("Error while fetching web urls", zap.Error(result.Error))
		return nil, result.Error
	}
	webUrlRepository.logger.Info(fmt.Sprintf("Web urls for interval %d", interval), zap.Any("webUrls", webUrls))
	return webUrls, nil
}

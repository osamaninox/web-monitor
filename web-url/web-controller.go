package web_url

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WebController struct {
	logger           *zap.Logger
	WebUrlRepository WebUrlRepository
}

func NewWebController(logger *zap.Logger, database *gorm.DB) *WebController {
	return &WebController{logger: logger, WebUrlRepository: WebUrlRepository{database: database, logger: logger}}
}

func (controller *WebController) AddWebUrlForMonitoring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var webUrl WebUrl
	err = json.Unmarshal(bodyBytes, &webUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	controller.logger.Info("Web url to be added", zap.Any("webUrl", webUrl))
	if webUrl.Url == "" {
		controller.logger.Error("Url is required")
		http.Error(w, "Url is required", http.StatusBadRequest)
		return
	}
	if webUrl.Interval == 0 {
		controller.logger.Error("Interval is required")
		http.Error(w, "Interval is required", http.StatusBadRequest)
		return
	}
	result := controller.WebUrlRepository.CreateWebUrl(&webUrl)
	if result.Error != nil {
		controller.logger.Error("Error while adding web url", zap.Error(result.Error))
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	webUrlBytes, err := json.Marshal(webUrl)
	if err != nil {
		controller.logger.Error("Error while marshalling web url", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(webUrlBytes)
}

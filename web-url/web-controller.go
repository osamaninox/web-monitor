package web_url

import (
	"encoding/json"
	"io"
	"net/http"

	"errors"

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
	var logger = controller.logger
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
	err = validateWebUrl(webUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info("Web url to be added", zap.Any("webUrl", webUrl))

	result := controller.WebUrlRepository.CreateWebUrl(&webUrl)
	if result.Error != nil {
		logger.Error("Error while adding web url", zap.Error(result.Error))
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	webUrlBytes, err := json.Marshal(webUrl)
	if err != nil {
		logger.Error("Error while marshalling web url", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(webUrlBytes)
}

func validateWebUrl(webUrl WebUrl) error {
	if webUrl.Url == "" {
		return errors.New("url is required")
	}
	if webUrl.Interval == 0 {
		return errors.New("interval is required")
	}
	return nil
}

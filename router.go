package main

import (
	"net/http"
	webUrl "web-monitor/web-url"

	"go.uber.org/zap"
)

func setupEndpoints(logger *zap.Logger, controller *webUrl.WebController) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/web-url", controller.AddWebUrlForMonitoring)

	return mux
}

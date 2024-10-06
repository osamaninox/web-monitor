package web_url

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WebController struct {
}

func (controller WebController) AddWebUrlForMonitoring(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("Web url to be added", webUrl)
	if webUrl.Url == "" {
		http.Error(w, "Url is required", http.StatusBadRequest)
		return
	}
	if webUrl.Interval == 0 {
		http.Error(w, "Interval is required", http.StatusBadRequest)
		return
	}
	result := CreateWebUrl(&webUrl)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	webUrlBytes, err := json.Marshal(webUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(webUrlBytes)
}

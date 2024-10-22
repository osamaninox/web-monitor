package main

import (
	"bytes"
	"net/http/httptest"
	"web-monitor/database"

	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func runTestServer() *httptest.Server {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	database.InitDB(logger)
	server := httptest.NewServer(setupEndpoints(logger, database.GetDB()))
	return server
}

func TestWebUrl(t *testing.T) {
	server := runTestServer()
	defer server.Close()

	t.Run("Should successfully add web url", func(t *testing.T) {
		// Arrange
		url := server.URL + "/web-url"
		body := []byte(`{"url": "http://example.com", "interval": 5}`)

		// Act
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		res, err := (&http.Client{}).Do(req)
		if err != nil {
			t.Fatal(err)
		}
		// Assert
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("Should return bad request error when required field is missing from request body", func(t *testing.T) {
		// Arrange
		url := server.URL + "/web-url"
		body := []byte(`{"interval": 5}`)

		// Act
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		res, err := (&http.Client{}).Do(req)
		if err != nil {
			t.Fatal(err)
		}
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, 400, res.StatusCode)
	})
}

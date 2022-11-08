package sdk

import (
	"net/http"
	"time"
)

type ClientDatadog struct {
	client  *http.Client
	url     string
	api_key string
	app_key string
}

// https://app.datadoghq.eu
func NewClient(url string, api_key string, app_key string) *ClientDatadog {
	client := &http.Client{Timeout: 10 * time.Second}
	return &ClientDatadog{
		client:  client,
		url:     url,
		api_key: api_key,
		app_key: app_key,
	}
}

func (cl *ClientDatadog) auth(req *http.Request) {
	req.Header.Set("DD-API-KEY", cl.api_key)
	req.Header.Set("DD-APPLICATION-KEY", cl.app_key)
}

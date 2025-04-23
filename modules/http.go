package modules

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	Status int
	Body   map[string]interface{}
}

func HttpGet(url string) interface{} {
	resp, err := http.Get(url)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)

	return map[string]interface{}{
		"status": resp.StatusCode,
		"body":   body,
	}
}

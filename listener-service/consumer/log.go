package consumer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func logEvent(eventData string) error {
	request, err := http.NewRequest("POST", "http://loggerApp:5000/log", bytes.NewBuffer([]byte(eventData)))
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var req struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
			Data    string `json:"data,omitempty"`
		}
		if err := json.Unmarshal(responseData, &req); err != nil {
			return err
		}

		return errors.New(req.Data)
	}

	return nil
}

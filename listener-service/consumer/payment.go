package consumer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func paymentEvent(event Payload) error {
	jsonData, _ := json.Marshal(event)
	request, err := http.NewRequest("POST", "http://paymentApp:5000/create-payment-intent", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var jsonFromService struct {
		ClientSecret string `json:"client_secret"`
	}
	err = json.Unmarshal(responseBody, &jsonFromService)
	if err != nil {
		return err
	}

	return nil
}

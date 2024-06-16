package api

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	Low      = uint8(2)
	Medium   = uint8(5)
	Critical = uint8(9)
)

type rabbitPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (server *Server) Publish(routingKey string, data rabbitPayload, priority uint8) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return server.ch.Publish(
		"events",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Priority:    priority,
			Body:        jsonData,
		},
	)
}

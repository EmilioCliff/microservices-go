package consumer

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	var consumer Consumer
	consumer.conn = conn

	err := consumer.setConsumer()
	if err != nil {
		return &Consumer{}, err
	}

	return &consumer, nil
}

func (consumer *Consumer) setConsumer() error {
	// Create the exchange
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.ExchangeDeclare("events", "direct", true, false, false, false, nil)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		err = ch.QueueBind(q.Name, topic, "events", false, nil)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Printf("Failed to unmarshal payload: %v", err)
				d.Nack(false, true)
				return
			}

			err = handlePayload(payload)
			if err != nil {
				log.Printf("Failed to handle payload: %v", err)
				d.Nack(false, true)
			} else {
				log.Printf("Message acknowledged: %v", d.DeliveryTag)
				d.Ack(false)
			}
		}
	}()

	log.Println("Waiting for messages...")

	<-forever

	return nil
}

func handlePayload(payload Payload) error {
	switch payload.Name {
	case "log":
		err := logEvent(payload.Data)
		if err != nil {
			return err
		}

	case "payment":
		err := paymentEvent(payload.Data)
		if err != nil {
			return err
		}

	default:
		err := logEvent(payload.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

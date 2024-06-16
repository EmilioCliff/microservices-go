package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/EmilioCliff/listener-service/consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := connectToRabit(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Println("failed to connect to rabbitmq")
		return
	}
	defer conn.Close()

	consumer, err := consumer.NewConsumer(conn)
	if err != nil {
		log.Println("failed to create new consumer")
		return
	}

	err = consumer.Listen([]string{"log", "payment"})
	if err != nil {
		log.Println("failed to listen to channel")
		return
	}
}

func connectToRabit(uri string) (*amqp.Connection, error) {
	count := 0
	rollOff := 1 * time.Second
	var err error
	var connection *amqp.Connection
	for {
		connection, err = amqp.Dial(uri)
		if err != nil {
			log.Println("failed to connect to rabbitmq", err)
			if count > 12 {
				return nil, err
			}
			count++
			rollOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
			time.Sleep(rollOff)
			continue
		}
		fmt.Println("Connected to rabbitmq")
		break
	}

	return connection, nil
}

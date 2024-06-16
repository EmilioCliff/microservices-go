package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/EmilioCliff/broker-service/api"
	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "5000"

func main() {
	connection, err := connectToRabit(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Println("failed to connect to rabbitmq")
		return
	}
	defer connection.Close()

	ch, err := connection.Channel()
	if err != nil {
		log.Println("failed to open channel")
	}
	defer ch.Close()

	server := api.NewServer(ch)

	fmt.Println("Starting broker on port: ", webPort)
	server.Start(fmt.Sprintf("0.0.0.0:%s", webPort))
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

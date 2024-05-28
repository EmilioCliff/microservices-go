package main

import (
	"fmt"

	"github.com/EmilioCliff/broker-service/api"
)

const webPort = "5000"

func main() {
	server := api.NewServer()

	fmt.Println("Starting broker on port: ", webPort)
	server.Start(fmt.Sprintf("0.0.0.0:%s", webPort))
}

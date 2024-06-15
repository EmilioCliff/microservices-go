package main

import (
	"fmt"
	"os"

	"github.com/EmilioCliff/payment-service/api"
	"github.com/stripe/stripe-go/v78"
)

// const webPort = "5000"

func main() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	server := api.NewServer()
	fmt.Println("Starting server on port: ", os.Getenv("PORT"))

	server.Start(fmt.Sprintf("0.0.0.0:%v", os.Getenv("PORT")))
}

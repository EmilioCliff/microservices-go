package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/EmilioCliff/logger-service/api"
	db "github.com/EmilioCliff/logger-service/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

const webPort = "5000"

func main() {
	coonPool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Panic("failed to connect to db")
	}

	if err = coonPool.Ping(context.Background()); err != nil {
		log.Panic("failed to ping db")
	}

	store := db.New(coonPool)
	server := api.NewServer(store)

	log.Println("Starting logger at port: ", webPort)
	if err = server.Start(fmt.Sprintf("0.0.0.0:%s", webPort)); err != nil {
		log.Panic("failed to start server")
	}
}

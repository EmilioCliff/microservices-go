package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/EmilioCliff/auth-service/api"
	db "github.com/EmilioCliff/auth-service/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

const webPort = "5000"

var count = 0

func main() {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Println("could not connect to db", err)
	}

	if err = conn.Ping(context.Background()); err != nil {
		log.Println("Did not ping to db")
	} else {
		log.Println("pong")
	}

	store := db.New(conn)
	server := api.NewServer(store)

	log.Println("Starting auth server at port: ", fmt.Sprintf("0.0.0.0:%s", webPort))
	server.Start(fmt.Sprintf("0.0.0.0:%s", webPort))
}

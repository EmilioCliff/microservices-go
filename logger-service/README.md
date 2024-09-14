# Logger Service 🚀

Welcome to the **Logger Service**! This microservice is part of the larger Go Microservices project. It is designed to log specific request values/parameters to a postgres database.

## Endpoints ✨

- Feature 1
  `GET    /ping` used to test if the server is up and healthy
  `POST     /log` used to log the request into the database

## Technologies Used 🛠️

- **gin-gonic**: A server to listen and produce http request
- **sqlc**: Used to generate type safe code from sql queries
- **golang-migrate**: Used to apply updates to our database schema while helping in version control of our db schema

## Configuration ⚙️

- Config setting 1: `DB_URL` Set in the compose file

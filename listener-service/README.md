# Listener Service 🚀

Welcome to the **Listener Service**! This microservice is part of the larger Go Microservices project. It is designed as a dedicated service that only deals with the set up of the rabbitmq exchange, queues and binding routing-keys. The queues are forwarded to respective services in here.

## Technologies Used 🛠️

- **amqp091-go**: Used to for interaction with the AMQP that being RabbitMQ where request are queued.

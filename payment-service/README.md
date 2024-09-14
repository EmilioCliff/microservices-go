# Payment Service 🚀

Welcome to the **Payment Service**! This microservice is part of the larger Go Microservices project. It is designed to have a secure card payment method using the Stripe.

## Endpoints ✨

- Feature 1
  `GET    /config` route used to get the `STRIPE_PUBLISHABLE_KEY`
  `POST     /create-payment-intent` used to initiate the payment intent
  `POST     /webhook` used as a webhook to the stripe callback

## Technologies Used 🛠️

- **gin-gonic**: A server to listen and produce http request
- **stripe-go**: A library that helps in the stripe implementation in golang

## Configuration ⚙️

- Config setting 1: `DB_URL` Set in the compose file
- Confiq setting 2: `STRIPE_SECRET_KEY=` Provided by Stripe
- Confiq setting 3: `STRIPE_PUBLISHABLE_KEY=` Provided by Stripe
- Confiq setting 4: `WEBHOOK=` Used a live callback url. I used [ngrok](.https://ngrok.com/)
- Confiq setting 5: `PORT=5000`

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/v78/webhook"
)

type PublishableKeyResponse struct {
	PublishableKey string `json:"publishable_key"`
}

func (server *Server) getPublishableKey(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, PublishableKeyResponse{
		PublishableKey: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
	})
}

type CreatePaymentIntentResponse struct {
	ClientSecret string `json:"client_secret"`
}

type CreatePaymentIntentRequest struct {
	Amount int64 `json:"amount"`
}

func (server *Server) createPaymentIntent(ctx *gin.Context) {
	var req CreatePaymentIntentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err))
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, CreatePaymentIntentResponse{
		ClientSecret: pi.ClientSecret,
	})
}

func (server *Server) handleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	endpointSecret := "whsec_15c4e50f8c6d4ad1a714cca17bea3a476b8b9ad0384945c81391d13b56ea7830"

	event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.Status(http.StatusBadRequest)
			return
		}
		fmt.Printf("PaymentIntent was successful: %v\n", paymentIntent.ID)
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	c.Status(http.StatusOK)
}

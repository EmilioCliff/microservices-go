package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

type JSONRequst struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

type RequestPayload struct {
	Action  string         `json:"action"`
	Auth    Authpayload    `json:"auth,omitempty"`
	Logger  LoggerPayload  `json:"logger,omitempty"`
	Payment PaymentPayload `json:"payment,omitempty"`
}

type Authpayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoggerPayload struct {
	Email     string `json:"email"`
	Data      string `json:"data"`
	UserIP    string `json:"user_ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

type PaymentPayload struct {
	Amount  int64  `json:"amount,omitempty"`
	Publish string `json:"publish,omitempty"`
}

func (server *Server) testBroker(ctx *gin.Context) {
	res := JSONRequst{
		Error:   false,
		Message: "hit the broker",
	}

	ctx.JSON(http.StatusAccepted, res)
}

func (server *Server) handler(ctx *gin.Context) {
	var req RequestPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "could create request"))
		return
	}

	switch req.Action {
	case "auth":
		server.authenticate(ctx, req.Auth)

	case "logger":
		// server.log(ctx, req.Logger)
		server.logToRabit(ctx, req.Logger)

	case "payment":
		if req.Payment.Publish == "request publishable-key" {
			server.getPaymentKey(ctx)
		} else {
			server.processPayment(ctx, req.Payment)
		}

	default:
		ctx.JSON(http.StatusBadRequest, JSONRequst{
			Error:   true,
			Message: "unknown action",
		})
	}
}

func (server *Server) authenticate(ctx *gin.Context, payload Authpayload) {
	jsonData, _ := json.Marshal(payload)
	request, err := http.NewRequest("POST", "http://authApp:5000/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "could create request"))
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "couldn't send request"))
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "could read response"))
		return
	}

	var jsonFromService JSONRequst

	err = json.Unmarshal(responseBody, &jsonFromService)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "couldnt unmarshal response"))
		return
	}

	if response.StatusCode == http.StatusUnauthorized {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New(jsonFromService.Message), "password dont match"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(jsonFromService.Message), "request failed"))
		return
	}

	if jsonFromService.Error {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(jsonFromService.Message), "request failed"))
		return
	} else {
		ctx.JSON(http.StatusOK, JSONRequst{
			Error:   false,
			Message: "Authenticated",
			Data:    jsonFromService.Data,
		})
	}
}

func (server *Server) log(ctx *gin.Context, payload LoggerPayload) {
	payload.UserAgent = ctx.Request.UserAgent()
	payload.UserIP = ctx.ClientIP()

	jsonData, _ := json.Marshal(payload)
	request, err := http.NewRequest("POST", "http://loggerApp:5000/log", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "failed to create response"))
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "failed to send request"))
		return
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "failed to read response"))
		return
	}

	var req JSONRequst
	if err := json.Unmarshal(responseData, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "failed to unmarshal response"))
		return
	}

	ctx.JSON(http.StatusOK, JSONResponse{
		Error:   false,
		Message: req.Message,
		Data:    req.Data,
	})
}

func (server *Server) getPaymentKey(ctx *gin.Context) {
	resp, err := http.Get("http://paymentApp:5000/config")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "errrorrring 1"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "errrorrring 2"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "errrorrring 3"})
		return
	}

	var response struct {
		PublishableKey string `json:"publishable_key"`
	}

	if err = json.Unmarshal(body, &response); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "errrorrring 4"})
		return
	}

	ctx.JSON(http.StatusOK, response)
	return
}

func (server *Server) processPayment(ctx *gin.Context, payload PaymentPayload) {
	jsonData, _ := json.Marshal(payload)
	request, err := http.NewRequest("POST", "http://paymentApp:5000/create-payment-intent", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "errrorrring 1"})
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "errrorrring 2"})
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "errrorrring 3"})
		return
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "errrorrring 4"})
		return
	}

	var jsonFromService struct {
		ClientSecret string `json:"client_secret"`
	}
	err = json.Unmarshal(responseBody, &jsonFromService)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "errrorrring 5"})
		return
	}

	ctx.JSON(http.StatusOK, jsonFromService)
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

	endpointSecret := os.Getenv("WEBHOOK_SECRET")

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

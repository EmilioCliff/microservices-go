package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

type Server struct {
	router *gin.Engine
	ch     *amqp.Channel
}

func NewServer(ch *amqp.Channel) *Server {
	var server Server
	server.ch = ch
	server.setRoutes()

	return &server
}

func (server *Server) setRoutes() {
	router := gin.Default()

	router.Use(CORSMiddleware())
	// logger := router.Group("/").Use(LogRequest())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/", server.testBroker)
	router.POST("/handler", server.handler)
	router.POST("/webhook", server.handleWebhook)

	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errorResponse(err error, message string) JSONResponse {
	return JSONResponse{
		Error:   true,
		Message: message,
		Data:    err.Error(),
	}
}

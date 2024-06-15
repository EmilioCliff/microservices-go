package api

import "github.com/gin-gonic/gin"

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	var server Server
	server.setRoutes()
	return &server
}

func (server *Server) setRoutes() {
	router := gin.Default()

	router.GET("/config", server.getPublishableKey)
	router.POST("/create-payment-intent", server.createPaymentIntent)
	router.POST("/webhook", server.handleWebhook)

	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run()
}

func (server *Server) errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

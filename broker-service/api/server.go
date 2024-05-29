package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

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

	router.Use(CORSMiddleware())
	logger := router.Group("/").Use(LogRequest())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/", server.testBroker)
	logger.POST("/handler", server.handler)

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

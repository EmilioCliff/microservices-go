package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/handle", server.handleBroker)

	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

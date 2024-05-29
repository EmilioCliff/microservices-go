package api

import (
	"net/http"

	db "github.com/EmilioCliff/logger-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Queries
	router *gin.Engine
}

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func NewServer(store *db.Queries) *Server {
	var server Server

	server.setRoutes()

	server.store = store
	return &server
}

func (server *Server) setRoutes() {
	router := gin.Default()

	router.Use(server.CORSMiddleware())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"ping": "pong"})
	})

	router.POST("/log", server.logEntry)

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

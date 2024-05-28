package api

import (
	"net/http"

	db "github.com/EmilioCliff/auth-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

type Server struct {
	store  *db.Queries
	router *gin.Engine
}

func NewServer(store *db.Queries) *Server {
	var server Server

	server.setRoutes()

	server.store = store

	return &server
}

func (server *Server) setRoutes() {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	router.POST("/authenticate", server.authHandler)

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

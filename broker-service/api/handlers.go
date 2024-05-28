package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONRequst struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (server *Server) handleBroker(ctx *gin.Context) {
	res := JSONRequst{
		Error:   false,
		Message: "hit the broker",
	}

	ctx.JSON(http.StatusAccepted, res)
}

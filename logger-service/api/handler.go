package api

import (
	"encoding/json"
	"net/http"

	db "github.com/EmilioCliff/logger-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type LogEntryRequest struct {
	Email     string `json:"email" binding:"required"`
	Data      string `json:"data" binding:"required"`
	UserIP    string `json:"user_ip" binding:"required"`
	UserAgent string `json:"user_agent" binding:"required"`
}

func (server *Server) logEntry(ctx *gin.Context) {
	var req LogEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, "bad request"))
		return
	}

	log, err := server.store.CreateLog(ctx, db.CreateLogParams{
		Email:     req.Email,
		Data:      req.Data,
		UserAgent: req.UserAgent,
		UserIp:    req.UserIP,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "failed to create log"))
		return
	}

	jsonData, _ := json.Marshal(log)

	ctx.JSON(http.StatusOK, JSONResponse{
		Error:   false,
		Message: "log entered",
		Data:    string(jsonData),
	})

}

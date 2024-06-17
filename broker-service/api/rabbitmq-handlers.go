package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) logToRabit(ctx *gin.Context, payload LoggerPayload) {
	payload.UserAgent = ctx.Request.UserAgent()
	payload.UserIP = ctx.ClientIP()

	j, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "could push to exchanger"))
		return
	}

	p := rabbitPayload{
		Name: "log",
		Data: string(j),
	}
	err = server.Publish("log", p, Medium)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "could push to exchanger"))
		return
	}

	ctx.JSON(http.StatusOK, JSONResponse{
		Error:   false,
		Message: "logged via rabbitMQ",
		Data:    "rabbitmq successful",
	})
}

func (server *Server) paymentToRabit() {

}

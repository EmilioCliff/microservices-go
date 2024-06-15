package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization, Accept")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		ctx.Writer.Header().Set("Access-Control-Exposed-Headers", "Link")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

func LogRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RequestPayload
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}

		payload := LoggerPayload{
			UserIP:    ctx.ClientIP(),
			UserAgent: ctx.Request.UserAgent(),
		}
		if req.Action == "auth" {
			payload.Email = req.Action
			payload.Data = fmt.Sprintf("Authenticating: %s", req.Auth.Email)
		} else if req.Action == "logger" {
			payload = req.Logger
			payload.UserAgent = ctx.Request.UserAgent()
			payload.UserIP = ctx.ClientIP()
		} else if req.Action == "payment" {
		}

		jsonData, _ := json.Marshal(payload)
		request, err := http.NewRequest("POST", "http://loggerApp:5000/log", bytes.NewBuffer(jsonData))
		if err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}
		defer response.Body.Close()

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}

		if response.StatusCode != 200 {
			ctx.AbortWithError(http.StatusBadGateway, errors.New(fmt.Sprintf("failed with status code: %v and data %v my data = %v", response.StatusCode, string(responseData), req)))
			return
		}

		ctx.Set("payload", req)

		ctx.Next()
	}
}

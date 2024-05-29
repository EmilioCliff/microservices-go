package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONRequst struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

type RequestPayload struct {
	Action string        `json:"action"`
	Auth   Authpayload   `json:"auth,omitempty"`
	Logger LoggerPayload `json:"logger,omitempty"`
}

type Authpayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoggerPayload struct {
	Email     string `json:"email"`
	Data      string `json:"data"`
	UserIP    string `json:"user_ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

func (server *Server) testBroker(ctx *gin.Context) {
	res := JSONRequst{
		Error:   false,
		Message: "hit the broker",
	}

	ctx.JSON(http.StatusAccepted, res)
}

func (server *Server) handler(ctx *gin.Context) {
	// var req RequestPayload
	// if err := ctx.ShouldBindJSON(&req); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err, "bad request"))
	// 	return
	// }
	set, exists := ctx.Get("payload")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get request payload"})
		return
	}

	// Type assert to the correct type
	req, ok := set.(RequestPayload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request payload type"})
		return
	}

	switch req.Action {
	case "auth":
		server.authenticate(ctx, req.Auth)

	case "logger":
		server.log(ctx, req.Logger)

	default:
		ctx.JSON(http.StatusBadRequest, JSONRequst{
			Error:   true,
			Message: "unknown action",
		})
	}
}

func (server *Server) authenticate(ctx *gin.Context, payload Authpayload) {
	jsonData, _ := json.Marshal(payload)
	request, err := http.NewRequest("POST", "http://authApp:5000/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "could create request"))
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "couldn't send request"))
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "could read response"))
		return
	}

	var jsonFromService JSONRequst

	err = json.Unmarshal(responseBody, &jsonFromService)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "couldnt unmarshal response"))
		return
	}

	if response.StatusCode == http.StatusUnauthorized {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New(jsonFromService.Message), "password dont match"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(jsonFromService.Message), "request failed"))
		return
	}

	if jsonFromService.Error {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(jsonFromService.Message), "request failed"))
		return
	} else {
		ctx.JSON(http.StatusOK, JSONRequst{
			Error:   false,
			Message: "Authenticated",
			Data:    jsonFromService.Data,
		})
	}
}

func (server *Server) log(ctx *gin.Context, payload LoggerPayload) {
	payload.UserAgent = ctx.Request.UserAgent()
	payload.UserIP = ctx.ClientIP()
	jsonData, _ := json.Marshal(payload)
	request, err := http.NewRequest("POST", "http://loggerApp:5000/log", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "failed to create response"))
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "failed to send request"))
		return
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "failed to read response"))
		return
	}

	var req JSONRequst
	if err := json.Unmarshal(responseData, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "failed to unmarshal response"))
		return
	}

	ctx.JSON(http.StatusOK, JSONResponse{
		Error:   false,
		Message: req.Message,
		Data:    req.Data,
	})
}

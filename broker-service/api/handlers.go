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
	Action string      `json:"action"`
	Auth   Authpayload `json:"auth,omitempty"`
}

type Authpayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (server *Server) testBroker(ctx *gin.Context) {
	res := JSONRequst{
		Error:   false,
		Message: "hit the broker",
	}

	ctx.JSON(http.StatusAccepted, res)
}

func (server *Server) handler(ctx *gin.Context) {
	var req RequestPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, "bad request"))
		return
	}

	switch req.Action {
	case "auth":
		server.authenticate(ctx, req.Auth)
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

	// if err = json.NewDecoder(request.Body).Decode(&jsonFromService); err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err, "could decode response"))
	// 	return
	// }

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

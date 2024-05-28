package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EmilioCliff/auth-service/db"
	"github.com/gin-gonic/gin"
)

type reqPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (server *Server) authHandler(ctx *gin.Context) {
	var req reqPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, "bad request"))
		return
	}

	user, err := server.store.GetUser(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err, "email not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "internal server error"))
		return
	}

	if err = db.ComparePassword(user.Password, req.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err, "password do not match"))
		return
	}

	rspByte, _ := json.Marshal(user)

	ctx.JSON(http.StatusAccepted, JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    string(rspByte),
	})
}

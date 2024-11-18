package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/user"
)

func ResetPOST(ctx *gin.Context) {
	type ResetPasswordRequest struct {
		Email string `json:"email"`
	}

	var request ResetPasswordRequest

	if ctx.Bind(&request) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	if request.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email is blank. All fields must be filled"})
		return
	}

	_, err := user.GetUser(request.Email)
	if err == errors.ErrInvalidEmail {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address used in query"})
		return
	} else if err == errors.ErrNoUser {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find user with the specified email address"})
		return
	}

	err = user.ResetUserPassword(request.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully sent reset password email to user"})
}

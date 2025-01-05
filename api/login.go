package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/user"
)

/*
LoginPOST Gin handler for the POST request to the Login Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func LoginPOST(ctx *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var request LoginRequest

	err := ctx.Bind(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "err": sdkErrors.ErrInvalidObjectStructure.Error()})
		return
	}

	if request.Email == "" || request.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email and/or password is blank. Both fields must have content"})
		return
	}

	_, err = user.GetUser(request.Email)
	if errors.Is(err, sdkErrors.ErrNoUser) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find the user account with the requested email address", "err": err.Error()})
		return
	}

	accessToken, err := user.LoginUser(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token", "err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
}

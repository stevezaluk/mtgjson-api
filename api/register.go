package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/user"
)

/*
RegisterPOST Gin handler for POST request to the register endpoint. This should not be called directly and
should only be passed to the gin router
*/
func RegisterPOST(ctx *gin.Context) {
	type RegisterRequest struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var request RegisterRequest

	if ctx.Bind(&request) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	if request.Email == "" || request.Username == "" || request.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email, username, or password is blank. All fields must be filled"})
		return
	}

	_, err := user.RegisterUser(request.Username, request.Email, request.Password)
	if errors.Is(err, sdkErrors.ErrInvalidPasswordLength) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "User password is not long enough. Password must be at least 12 characters, 1 special character, and 1 number"})
		return
	} else if errors.Is(err, sdkErrors.ErrInvalidEmail) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "User email is not valid or is not an email address"})
		return
	} else if errors.Is(err, sdkErrors.ErrFailedToRegisterUser) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register the user with Auth0"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully registered"})
}

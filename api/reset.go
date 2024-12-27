package api

import (
	"errors"
	"mtgjson/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/user"
)

/*
ResetGET Gin handler for POST request to the reset endpoint. This should not be called directly and
should only be passed to the gin router
*/
func ResetGET(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	email := ctx.DefaultQuery("email", userEmail)

	if email != userEmail {
		if !auth.ValidateScope(ctx, "write:user") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to reset other users passwords", "requiredScope": "write:user"})
			return
		}
	}

	_, err := user.GetUser(email)
	if errors.Is(err, sdkErrors.ErrInvalidEmail) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address used in query"})
		return
	} else if errors.Is(err, sdkErrors.ErrNoUser) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find user with the specified email address"})
		return
	}

	err = user.ResetUserPassword(email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully sent reset password email to user"})
}

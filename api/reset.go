package api

import (
	"errors"
	"github.com/stevezaluk/mtgjson-sdk/server"
	"mtgjson/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/user"
)

/*
ResetGET Gin handler for the GET request to the Reset Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func ResetGET(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		email := ctx.DefaultQuery("email", userEmail)

		if email != userEmail {
			if !auth.ValidateScope(ctx, "write:user") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to reset other users passwords", "requiredScope": "write:user"})
				return
			}
		}

		_, err := user.GetUser(server.Database(), email)
		if errors.Is(err, sdkErrors.ErrInvalidEmail) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address used in query", "err": err.Error()})
			return
		} else if errors.Is(err, sdkErrors.ErrNoUser) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find user with the specified email address", "err": err.Error()})
			return
		}

		err = server.AuthenticationManager().ResetUserPassword(email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to reset user password", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully sent reset password email to user"})
	}
}

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
UserGET Gin handler for the GET request to the User Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func UserGET(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		email := ctx.DefaultQuery("email", userEmail)

		if email != userEmail { // externalize logic for fetching own profile to a separate endpoint
			if !auth.ValidateScope(ctx, "read:user") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to read other user's account data", "requiredScope": "read:user"})
				return
			}
		}

		if email == "" {
			limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
			result, err := user.IndexUsers(server.Database(), limit)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find users in database", "err": err.Error()})
				return
			}

			ctx.JSON(http.StatusOK, result)
			return
		}

		result, err := user.GetUser(server.Database(), email)
		if errors.Is(err, sdkErrors.ErrNoUser) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find user with the specified email address", "err": err.Error()})
			return
		} else if errors.Is(err, sdkErrors.ErrInvalidEmail) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address used in query", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}

/*
UserDELETE Gin handler for the DELETE request to the User Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func UserDELETE(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		email := ctx.DefaultQuery("email", userEmail)

		if email != userEmail {
			if !auth.ValidateScope(ctx, "write:user") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to delete other users", "requiredScope": "write:user"})
				return
			}
		}

		if email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "An email address must be used to delete an account", "err": sdkErrors.ErrUserMissingId.Error()})
			return
		}

		err := server.AuthenticationManager().DeactivateUser(email)
		if errors.Is(err, sdkErrors.ErrNoUser) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find user with the specified email address", "err": err.Error()})
			return
		} else if errors.Is(err, sdkErrors.ErrInvalidEmail) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address used in query", "err": err.Error()})
			return
		} else if errors.Is(err, sdkErrors.ErrUserDeleteFailed) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user from MongoDB. User account may still be active", "err": err.Error()})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to deactivate user", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deactivated user account"})
	}
}

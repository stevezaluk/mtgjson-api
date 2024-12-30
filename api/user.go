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
UserGET Gin handler for the GET request to the User Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func UserGET(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	email := ctx.DefaultQuery("email", userEmail)

	if email != userEmail {
		if !auth.ValidateScope(ctx, "read:user") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to read other user's account data", "requiredScope": "read:user"})
			return
		}
	}

	if email == "" {
		limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
		result, err := user.IndexUsers(limit)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find users in database"})
			return
		}

		ctx.JSON(http.StatusOK, result)
		return
	}

	result, err := user.GetUser(email)
	if errors.Is(err, sdkErrors.ErrNoUser) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find user with the specified email address"})
		return
	} else if errors.Is(err, sdkErrors.ErrInvalidEmail) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address used in query"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

/*
UserDELETE Gin handler for the DELETE request to the User Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func UserDELETE(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	email := ctx.DefaultQuery("email", userEmail)

	if email != userEmail {
		if !auth.ValidateScope(ctx, "write:user") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to delete other users", "requiredScope": "write:user"})
			return
		}
	}

	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "An email address must be used to delete an account"})
		return
	}

	err := user.DeactivateUser(email)
	if errors.Is(err, sdkErrors.ErrNoUser) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find user with the specified email address"})
		return
	} else if errors.Is(err, sdkErrors.ErrInvalidEmail) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address used in query"})
		return
	} else if errors.Is(err, sdkErrors.ErrUserDeleteFailed) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user from MongoDB. User account may still be active"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deactivated user account"})
}

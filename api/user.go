package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/user"
)

/*
Gin handler for GET request to the user endpoint. This should not be called directly and
should only be passed to the gin router
*/
func UserGET(ctx *gin.Context) {
	email := ctx.Query("email")

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
	if err == errors.ErrNoUser {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find user with the specified email address"})
		return
	} else if err == errors.ErrInvalidEmail {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address used in query"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stevezaluk/mtgjson-sdk/user"
	"net/http"
)

/*
StoreUserEmailHandler Gin handler for fetching and storing the users email address after there token has been validated. This function
is crucial in evaluating user ownership over objects
*/
func StoreUserEmailHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userEmail, err := user.GetEmailFromToken(ctx.GetString("tokenStr"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch email from access token. This is needed for establishing ownership in created objects", "err": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("userEmail", userEmail)
	}
}

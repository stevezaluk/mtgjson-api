package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"mtgjson/auth"
	"net/http"
	"strings"
)

/*
ValidateTokenHandler Gin handler for validating tokens received from your Auth0 tenant. An Authorization header is
required to be passed in the request for this to properly function. If the token is valid, then it
is stored in the gin context under 'token'. If the token is invalid, the request is aborted.
*/
func ValidateTokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing from request"}) // format this better
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		tokenValidator, err := auth.GetTokenValidator()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start token validator", "err": err.Error()}) // format this better
			ctx.Abort()
			return
		}

		token, err := tokenValidator.ValidateToken(
			context.Background(),
			tokenStr,
		)

		if token == nil || err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Token is not valid", "err": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("token", token)
		ctx.Set("tokenStr", tokenStr)
	}
}

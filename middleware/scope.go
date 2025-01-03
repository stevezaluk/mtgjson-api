package middleware

import (
	"github.com/gin-gonic/gin"
	"mtgjson/auth"
	"net/http"
)

/*
ValidateScopeHandler Gin handler for validating custom claims returned with the token. This is added as a handler in between the ValidateToken
handler and the core logic handler for the defined route.
*/
func ValidateScopeHandler(requiredScope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !auth.ValidateScope(ctx, requiredScope) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to access this resource", "requiredScope": requiredScope})
			ctx.Abort()
			return
		}
	}
}

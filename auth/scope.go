package auth

import (
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

/*
ValidateScope Fetch validated claims from the gin context and ensure that
the user has the desired scope
*/
func ValidateScope(ctx *gin.Context, requiredScope string) bool {
	token := ctx.Value("token").(*validator.ValidatedClaims)

	claims := token.CustomClaims.(*CustomClaims)
	if !claims.HasScope(requiredScope) {
		return false
	}

	return true
}

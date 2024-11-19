package auth

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetIssuerUrl() *url.URL {
	issuerUrl, err := url.Parse("https://" + viper.GetString("auth0.domain") + "/")
	if err != nil {
		panic(err) // fatal error
	}

	return issuerUrl
}

func GetValidator() (*validator.Validator, error) {
	issuer := GetIssuerUrl()
	provider := jwks.NewCachingProvider(issuer, 5*time.Minute)

	validator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuer.String(),
		[]string{viper.GetString("auth0.audience")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
	)

	if err != nil {
		return validator, err
	}

	return validator, nil
}

func ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if viper.GetBool("api.no_auth") { // if no auth is set, return to the next handler
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "Authorization header is missing from request"}) // format this better
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		validator, err := GetValidator()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to start token validator"}) // format this better
			ctx.Abort()
			return
		}

		token, err := validator.ValidateToken(
			context.Background(),
			tokenStr,
		)

		if token == nil || err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "Token is not valid", "err": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("token", token)
	}
}

func ValidateScope(requiredScope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Value("token").(*validator.ValidatedClaims)

		claims := token.CustomClaims.(*CustomClaims)
		if !claims.HasScope(requiredScope) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Missing required scoped permission to access this resource"})
			ctx.Abort()
			return
		}
	}
}

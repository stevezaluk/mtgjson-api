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

/*
Properly formats the domain of your Auth0 tenant to be used when creating a new JWT validator.
This value is pulled from viper under the property 'auth0.domain'
*/
func GetIssuerUrl() *url.URL {
	issuerUrl, err := url.Parse("https://" + viper.GetString("auth0.domain") + "/")
	if err != nil {
		panic(err) // fatal error
	}

	return issuerUrl
}

/*
Creates a new JWT token validator for use witin the ValidateToken middleware. The object that
this function returns provides logic for validating JWT tokens and unmarshaling custom claims
defined in your Auth0 tenant
*/
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

/*
Gin handler for validating tokens received from your Auth0 tenant. An Authorization header is
required to be passed in the request for this to properly function. If the token is valid, then it
is stored in the gin context under 'token'. If the token is invalid, the request is aborted.
Additionally, if the 'api.no_auth' flag is set, the validator returns to the next handler without any validation
*/
func ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if viper.GetBool("api.no_auth") { // if no auth is set, return to the next handler
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing from request"}) // format this better
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		validator, err := GetValidator()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start token validator"}) // format this better
			ctx.Abort()
			return
		}

		token, err := validator.ValidateToken(
			context.Background(),
			tokenStr,
		)

		if token == nil || err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Token is not valid", "err": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("token", token)
	}
}

/*
Gin handler for validating custom claims returned with the token. This is added as a handler in between the ValidateToken
handler and the core logic handler for the defined route.
*/
func ValidateScope(requiredScope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if viper.GetBool("api.no_scope") {
			return
		}

		token := ctx.Value("token").(*validator.ValidatedClaims)

		claims := token.CustomClaims.(*CustomClaims)
		if !claims.HasScope(requiredScope) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to access this resource", "missingScope": requiredScope})
			ctx.Abort()
			return
		}
	}
}

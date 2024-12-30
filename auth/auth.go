package auth

import (
	"context"
	"github.com/stevezaluk/mtgjson-sdk/user"
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
GetIssuerUrl Properly formats the domain of your Auth0 tenant to be used when creating a new JWT validator.
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
GetValidator Creates a new JWT token validator for use within the ValidateToken middleware. The object that
this function returns provides logic for validating JWT tokens and unmarshalling custom claims
defined in your Auth0 tenant
*/
func GetValidator() (*validator.Validator, error) {
	issuer := GetIssuerUrl()
	provider := jwks.NewCachingProvider(issuer, 5*time.Minute)

	tokenValidator, err := validator.New(
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
		return tokenValidator, err
	}

	return tokenValidator, nil
}

/*
ValidateTokenHandler Gin handler for validating tokens received from your Auth0 tenant. An Authorization header is
required to be passed in the request for this to properly function. If the token is valid, then it
is stored in the gin context under 'token'. If the token is invalid, the request is aborted.
Additionally, if the 'api.no_auth' flag is set, the validator returns to the next handler without any validation
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

		tokenValidator, err := GetValidator()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start token validator"}) // format this better
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

/*
ValidateScopeHandler Gin handler for validating custom claims returned with the token. This is added as a handler in between the ValidateToken
handler and the core logic handler for the defined route.
*/
func ValidateScopeHandler(requiredScope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !ValidateScope(ctx, requiredScope) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to access this resource", "missingScope": requiredScope})
			ctx.Abort()
			return
		}
	}
}

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

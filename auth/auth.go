package auth

import (
	"net/url"
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

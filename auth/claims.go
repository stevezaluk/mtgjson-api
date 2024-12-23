package auth

import (
	"context"
	"strings"
)

/*
CustomClaims Struct for unmarshalling Auth0 scopes during token validation
*/
type CustomClaims struct {
	Scope string `json:"scope"`
}

/*
Validate Currently this does nothing and is here to satisfy the CustomClaims interface
as defined in jwt-middleware
*/
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

/*
HasScope Validates that the expected scope exists within the CustomClaims struct
*/
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

package auth

import (
	"context"
	"strings"
)

type CustomClaims struct {
	Scope string `json:"scope"`
}

// need to satisfy interface requirements
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

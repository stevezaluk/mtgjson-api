package auth

import (
	"net/url"

	"github.com/spf13/viper"
)

func GetIssuerUrl() string {
	issuerUrl, err := url.Parse("https://" + viper.GetString("auth0.domain") + "/")
	if err != nil {
		panic(err) // fatal error
	}

	return issuerUrl.String()
}

package api

import (
	"errors"
	"github.com/stevezaluk/mtgjson-sdk/server"
	"net/http"

	"github.com/gin-gonic/gin"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/user"
)

/*
LoginPOST Gin handler for the POST request to the Login Endpoint. This function should not be called
directly and should only be passed to the gin router.
*/
func LoginPOST(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var request LoginRequest

		err := ctx.Bind(&request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "err": sdkErrors.ErrInvalidObjectStructure.Error()})
			return
		}

		if request.Email == "" || request.Password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email and/or password is blank. Both fields must have content"})
			return
		}

		_, err = user.GetUser(server.Database(), request.Email)
		if errors.Is(err, sdkErrors.ErrNoUser) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find the user account with the requested email address", "err": err.Error()})
			return
		}

		// this is returning 500 for all status codes. This has to get re-worked
		accessToken, err := server.AuthenticationManager().AuthenticateUser(request.Email, request.Password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to generate token", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, accessToken)
	}
}

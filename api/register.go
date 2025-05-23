package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	userModel "github.com/stevezaluk/mtgjson-models/user"
	"github.com/stevezaluk/mtgjson-sdk/server"
	"github.com/stevezaluk/mtgjson-sdk/user"
	"net/http"
)

/*
RegisterPOST Gin handler for the POST request to the Register Endpoint. This function should not be called
directly and should only be passed to the gin router. Revalidate this
*/
func RegisterPOST(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		type RegisterRequest struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var request RegisterRequest

		err := ctx.Bind(&request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "err": sdkErrors.ErrInvalidObjectStructure.Error()})
			return
		}

		if request.Email == "" || request.Username == "" || request.Password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email, username, or password is blank. All fields must be filled"})
			return
		}

		signUpResp, err := server.AuthenticationManager().RegisterUser(request.Username, request.Email, request.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user in Auth0", "err": err.Error()})
			return
		}

		err = user.NewUser(server.Database(), &userModel.User{
			Username: request.Username,
			Email:    request.Email,
			Auth0Id:  signUpResp.ID,
		})

		if errors.Is(err, sdkErrors.ErrInvalidPasswordLength) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "User password is not long enough. Password must be at least 12 characters, 1 special character, and 1 number", "err": err.Error()})
			return
		} else if errors.Is(err, sdkErrors.ErrInvalidEmail) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "User email is not valid or is not an email address", "err": err.Error()})
			return
		} else if errors.Is(err, sdkErrors.ErrFailedToRegisterUser) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register the user with Auth0", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User successfully registered"})
	}
}

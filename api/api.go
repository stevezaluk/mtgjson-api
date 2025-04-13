package api

import (
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/stevezaluk/mtgjson-sdk/server"
)

/*
API - An abstraction of the API as a whole
*/
type API struct {
	// server - The server structure that will be used for accessing the database and authentication api
	server *server.Server

	// router - The primary gin router used for routing endpoints on the API
	router *gin.Engine
}

/*
New - A constructor for the API structure
*/
func New(server *server.Server) *API {
	router := gin.New()
	router.Use(gin.Recovery(), sloggin.New(server.Log().Logger()))

	return &API{
		server: server,
	}
}

/*
FromConfig - Initialize the API structure using values from Viper
*/
func FromConfig() *API {
	return New(server.FromConfig())
}

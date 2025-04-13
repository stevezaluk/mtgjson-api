package api

import (
	"github.com/gin-gonic/gin"
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

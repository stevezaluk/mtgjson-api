package api

import (
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/spf13/viper"
	"github.com/stevezaluk/mtgjson-sdk/server"
	"log/slog"
	"mtgjson/middleware"
	"strconv"
)

// HandlerFunc - Wraps all handler functions to ensure that they can get passed a reference to the Server structure
type HandlerFunc func(server *server.Server) gin.HandlerFunc

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
		router: router,
	}
}

/*
FromConfig - Initialize the API structure using values from Viper
*/
func FromConfig() *API {
	return New(server.FromConfig())
}

/*
RegisterEndpoint - Registers an endpoint with the API. Method is the HTTP method that you want to
use on the path parameter, and the scope is the minimum required scope that will be required to
access the endpoint. If an empty string is provided to the scope, then one won't be required to
access it
*/
func (api *API) RegisterEndpoint(method string, path string, scope string, hasAuth bool, handler HandlerFunc) {
	handlers := []gin.HandlerFunc{
		handler(api.server),
	}

	if hasAuth {
		handlers = append(handlers, middleware.ValidateTokenHandler(api.server))
	}

	if scope != "" {
		handlers = append(handlers, middleware.ValidateScopeHandler(scope))
	}

	api.router.Handle(method, path, handlers...)
}

/*
Run - Connect to the MongoDB database and Start the API Server. The port parameter should describe
the port you want to expose the API on
*/
func (api *API) Run(port int) error {
	slog.Info("Initiating connection to MongoDB", "hostname", viper.GetString("mongo.hostname"))
	err := api.server.Database().Connect()
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "err", err)
		return err
	}

	slog.Info("Starting API Server", "port", port)
	err = api.router.Run(":" + strconv.Itoa(port))
	if err != nil {
		slog.Error("Failed to start API Server", "err", err)
		return err
	}

	return nil
}

/*
Shutdown - Gracefully stop the API Server and then disconnect from MongoDB
*/
func (api *API) Shutdown() error {
	slog.Info("Shutting down API")
	// stop receiving connections here
	// gin router doesn't provide this natively

	slog.Info("Disconnecting from MongoDB", "hostname", viper.GetString("mongo.hostname"))
	err := api.server.Database().Disconnect()
	if err != nil {
		slog.Error("Failed to disconnect from MongoDB", "err", err)
		return err
	}

	return nil
}

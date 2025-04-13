package api

import (
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/spf13/viper"
	"github.com/stevezaluk/mtgjson-sdk/server"
	"log/slog"
	"strconv"
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

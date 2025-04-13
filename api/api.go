package api

import (
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/spf13/viper"
	"github.com/stevezaluk/mtgjson-sdk/server"
	"log/slog"
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
Start - Connect to the MongoDB database and Start the API Server
*/
func (api *API) Start() error {
	slog.Info("Initiating connection to MongoDB", "hostname", viper.GetString("mongo.hostname"))
	err := api.server.Database().Connect()
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "err", err)
		return err
	}

	slog.Info("Starting API Server", "port", viper.GetInt("port"))
	err = api.router.Run(":" + viper.GetString("port"))
	if err != nil {
		slog.Error("Failed to start API Server", "err", err)
		return err
	}

	return nil
}

/*
Stop - Gracefully stop the API Server and then disconnect from MongoDB
*/
func (api *API) Stop() error {
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

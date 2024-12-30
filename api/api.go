package api

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/slog-gin"
	"github.com/stevezaluk/mtgjson-sdk/context"
	"log/slog"
	"mtgjson/auth"
	"strconv"
)

/*
API Abstraction of a Gin API. This stores the gin router and provides a scalable
way to add additional routes in the future. Call api.New() to create a new instance
of this object
*/
type API struct {
	Router *gin.Engine
}

/*
Init Initializes the database, logger, auth api, and management API and provides them to the gin router as middleware
*/
func (api API) Init() {
	context.InitLog()
	context.InitDatabase()
	context.InitAuthAPI()
	context.InitAuthManagementAPI()
}

/*
AddCardRoutes Add GET, POST, and DELETE routes to the API for the card namespace
*/
func (api API) addCardRoutes() {
	api.Router.GET("/api/v1/card", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("read:card"), CardGET)
	api.Router.POST("/api/v1/card", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:card"), CardPOST)
	api.Router.DELETE("/api/v1/card", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:card"), CardDELETE)
}

/*
AddDeckRoutes Add GET, POST, and DELETE routes to the API for the deck and deck content namespace
*/
func (api API) addDeckRoutes() {
	api.Router.GET("/api/v1/deck", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("read:deck"), DeckGET)
	api.Router.POST("/api/v1/deck", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:deck"), DeckPOST)
	api.Router.DELETE("/api/v1/deck", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:deck"), DeckDELETE)

	api.Router.GET("/api/v1/deck/content", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("read:deck"), DeckContentGET)
	api.Router.POST("/api/v1/deck/content", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:deck"), DeckContentPOST)
	api.Router.DELETE("/api/v1/deck/content", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:deck"), DeckContentDELETE)
}

func (api API) addSetRoutes() {
	api.Router.GET("/api/v1/set", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("read:set"), SetGET)
	api.Router.POST("/api/v1/set", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:set"), SetPOST)
	api.Router.DELETE("/api/v1/set", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:set"), SetDELETE)
	api.Router.GET("/api/v1/set/content", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("read:set"), SetContentGET)
	api.Router.POST("/api/v1/set/content", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:set"), SetContentPOST)
}

/*
AddUserRoutes Add GET and DELETE routes to the API for the user namespace
*/
func (api API) addUserRoutes() {
	api.Router.GET("/api/v1/user", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("read:profile"), UserGET)
	api.Router.DELETE("/api/v1/user", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("write:user"), UserDELETE)
}

/*
addManagementRoutes Add GET and POST routes to the API for the health and (eventually) the metrics endpoint
*/
func (api API) addManagementRoutes() {
	api.Router.GET("/api/v1/health", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("read:health"), HealthGET)
}

/*
AddAuthRoutes Add GET and POST routes to the API for the login, register, and reset password endpoints
*/
func (api API) addAuthRoutes() {
	api.Router.POST("/api/v1/login", LoginPOST)
	api.Router.POST("/api/v1/register", RegisterPOST)
	api.Router.GET("/api/v1/reset", auth.ValidateTokenHandler(), auth.StoreUserEmailHandler(), auth.ValidateScopeHandler("reset:password"), ResetGET)
}

func (api API) AddRoutes(routes []string) {
	api.addManagementRoutes()

	for _, route := range routes {
		if route == "card" {
			api.addCardRoutes()
		}

		if route == "deck" {
			api.addDeckRoutes()
		}

		if route == "set" {
			api.addSetRoutes()
		}

		if route == "user" {
			api.addUserRoutes()
		}

		if route == "auth" {
			api.addAuthRoutes()
		}
	}
}

/*
Start the API and add management routes to the router
*/
func (api API) Start(port int) {
	err := api.Router.Run(":" + strconv.Itoa(port))
	if err != nil {
		slog.Error("Failed to start api", "err", err.Error())
		return
	}
}

/*
Stop Destroy and release the database, and log file
*/
func (api API) Stop() {
	context.DestroyDatabase()
}

/*
New Creates a new instance of api.API and returns it
*/
func New() API {
	var router = gin.New()

	api := API{Router: router}
	api.Init()

	api.Router.Use(
		sloggin.New(context.GetLogger()),
		gin.Recovery(),
	)

	return api
}

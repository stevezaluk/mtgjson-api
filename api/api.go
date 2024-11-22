package api

import (
	"mtgjson/auth"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/slog-gin"
	"github.com/stevezaluk/mtgjson-sdk/context"
)

/*
Abstraction of a Gin API. This stores the gin router and provides a scalable
way to add additional routes in the future. Call api.New() to create a new instance
of this object
*/
type API struct {
	Router *gin.Engine
}

/*
Initializes the database, logger, auth api, and management API and provides them to the gin router as middleware
*/
func (api API) Init() {
	context.InitLog()
	context.InitDatabase()
	context.InitAuthAPI()
	context.InitAuthManagementAPI()
}

/*
Add GET, POST, and DELETE routes to the API for the card namespace
*/
func (api API) AddCardRoutes() {
	api.Router.GET("/api/v1/card", auth.ValidateToken(), auth.ValidateScope("read:card"), CardGET)
	api.Router.POST("/api/v1/card", auth.ValidateToken(), auth.ValidateScope("write:card"), CardPOST)
	api.Router.DELETE("/api/v1/card", auth.ValidateToken(), auth.ValidateScope("write:card"), CardDELETE)
}

/*
Add GET, POST, and DELETE routes to the API for the deck and deck content namespace
*/
func (api API) AddDeckRoutes() {
	api.Router.GET("/api/v1/deck", auth.ValidateToken(), auth.ValidateScope("read:deck"), DeckGET)
	api.Router.POST("/api/v1/deck", auth.ValidateToken(), auth.ValidateScope("write:deck"), DeckPOST)
	api.Router.DELETE("/api/v1/deck", auth.ValidateToken(), auth.ValidateScope("write:deck"), DeckDELETE)

	api.Router.GET("/api/v1/deck/content", auth.ValidateToken(), auth.ValidateScope("read:deck"), DeckContentGET)
	api.Router.POST("/api/v1/deck/content", auth.ValidateToken(), auth.ValidateScope("write:deck"), DeckContentPOST)
	api.Router.DELETE("/api/v1/deck/content", auth.ValidateToken(), auth.ValidateScope("write:deck"), DeckContentDELETE)
}

/*
Add GET and DELETE routes to the API for the user namespace
*/
func (api API) AddUserRoutes() {
	api.Router.GET("/api/v1/user", auth.ValidateToken(), auth.ValidateScope("read:user"), UserGET)
	api.Router.DELETE("/api/v1/user", auth.ValidateToken(), auth.ValidateScope("write:user"), UserDELETE)
}

/*
Add GET and POST routes to the API for the login, resgister, reset, and health endpoints
*/
func (api API) addManagementRoutes() {
	api.Router.GET("/api/v1/health", auth.ValidateToken(), auth.ValidateScope("read:health"), HealthGET)
	api.Router.POST("/api/v1/login", LoginPOST)
	api.Router.POST("/api/v1/register", RegisterPOST)
	api.Router.POST("/api/v1/reset", auth.ValidateToken(), ResetPOST)
}

/*
Start the API and add management routes to the router
*/
func (api API) Start(port int) {
	api.addManagementRoutes()
	api.Router.Run(":" + strconv.Itoa(port))
}

/*
Destrory and release the database, and log file
*/
func (api API) Stop() {
	context.DestroyDatabase()
}

/*
Creates a new instance of api.API and returns it
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

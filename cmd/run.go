package cmd

import (
	"mtgjson/api"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stevezaluk/mtgjson-sdk/context"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the API in the foreground and use STDOUT for logging",
	Long: `To start the API using the default config path:
$ mtgjson run

To start the API using a custom config file:
$ mtgjson run -c /path/to/config/.json

To start the API using environmental variables
$ mtgjson run --env`,
	PreRun: func(cmd *cobra.Command, args []string) {
		context.InitDatabase()

		debugMode := viper.GetBool("debug")
		if !debugMode {
			gin.SetMode(gin.ReleaseMode)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.New()

		logger := context.GetLogger()

		router.Use(
			sloggin.New(logger),
			gin.Recovery(),
		)

		router.GET("/api/v1/health", api.HealthGET)

		router.GET("/api/v1/card", api.CardGET)
		router.POST("/api/v1/card", api.CardPOST)
		router.DELETE("/api/v1/card", api.CardDELETE)

		router.GET("/api/v1/deck", api.DeckGET)
		router.POST("/api/v1/deck", api.DeckPOST)
		router.DELETE("/api/v1/deck", api.DeckDELETE)

		router.GET("/api/v1/deck/content", api.DeckContentGET)
		router.POST("/api/v1/deck/content", api.DeckContentPOST)
		router.DELETE("/api/v1/deck/content", api.DeckContentDELETE)

		router.Run()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		context.DestroyDatabase()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	context.InitConfig(cfgFile)
	context.InitLog()

	runCmd.Flags().BoolP("debug", "d", false, "Enable Gin debug mode. Release mode is set by default")
	viper.BindPFlag("debug", runCmd.Flags().Lookup("debug"))

	runCmd.Flags().StringP("log.path", "l", "/var/log/mtgjson-api", "Set the directory that the API should save logs too")
	viper.BindPFlag("log.path", runCmd.Flags().Lookup("log.path"))

	runCmd.Flags().IntP("api.port", "p", 8080, "Set the host port that the API should serve on")
	viper.BindPFlag("api.port", runCmd.Flags().Lookup("api.port"))

	runCmd.Flags().String("mongo.ip", "127.0.0.1", "Set the IP Address of your running MongoDB instance")
	viper.BindPFlag("mongo.ip", runCmd.Flags().Lookup("mongo.ip"))

	runCmd.Flags().String("mongo.port", "127.0.0.1", "Set the Port of your running MongoDB instance")
	viper.BindPFlag("mongo.port", runCmd.Flags().Lookup("mongo.port"))

	runCmd.Flags().String("mongo.user", "127.0.0.1", "Set the username to use for authentication with MongoDB")
	viper.BindPFlag("mongo.user", runCmd.Flags().Lookup("mongo.user"))

	runCmd.Flags().String("mongo.pass", "127.0.0.1", "Set the password to use for authentication with MongoDB")
	viper.BindPFlag("mongo.pass", runCmd.Flags().Lookup("mongo.pass"))

	runCmd.Flags().String("auth0.domain", "", "The domain of your Auth0 tenant")
	viper.BindPFlag("auth0.domain", runCmd.Flags().Lookup("auth0.domain"))

	runCmd.Flags().String("auth0.audience", "", "The identifier of your Auth0 API")
	viper.BindPFlag("auth0.audience", runCmd.Flags().Lookup("auth0.audience"))

	runCmd.Flags().String("auth0.client_id", "", "The Client ID for your Auth0 API")
	viper.BindPFlag("auth0.client_id", runCmd.Flags().Lookup("auth0.client_id"))

	runCmd.Flags().String("auth0.client_secret", "", "The Client Secret for your Auth0 APi")
	viper.BindPFlag("auth0.client_secret", runCmd.Flags().Lookup("auth0.client_secret"))
}

package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"mtgjson/api"
	"mtgjson/context"
	"mtgjson/server"
)

var defaultConfig string = "~/.config/mtgjson/config.json"
var ServerConfig server.Config

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
		context.InitConfig(cmd.PersistentFlags())
		context.InitDatabase()
	},
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()
		router.GET("/health", api.HealthGET)

		router.GET("/card", api.CardGET)
		router.POST("/card", api.CardPOST)

		router.GET("/deck", api.DeckGET)
		router.POST("/deck", api.DeckPOST)
		router.GET("/deck/content", api.DeckContentGET)
		router.PUT("/deck/content", api.DeckContentPUT)
		router.DELETE("/deck/content", api.DeckContentDELETE)

		router.Run()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		context.DestroyDatabase()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringP("config", "c", defaultConfig, "The path to your MTGJSON config file")
	runCmd.PersistentFlags().BoolP("env", "e", false, "Ignore the default config path and attempt to use Environmental Variables")
	runCmd.PersistentFlags().Int64P("port", "p", 2100, "Set the default port that the API listens on")
}

package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log/slog"
	"mtgjson/api"
	"os"

	"github.com/spf13/cobra"
)

// cfgFile - When -c or --config is called, the user supplied path is stored here
var cfgFile string

// rootCmd - The root command. Provides logic and help messages
var rootCmd = &cobra.Command{
	Use:   "mtgjson-api",
	Short: "",
	Long: `A command line tool for controlling the execution of the MTGJSON RESTful API.

MTGJSON API is unofficial Fan Content permitted under the Fan Content Policy. 
Not approved/endorsed by Wizards of the Coast. Portions of the materials used are property of Wizards of the Coast. 
© Wizards of the Coast LLC.

MTGJSON API not officially endorsed by MTGJSON`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("debug") {
			gin.SetMode(gin.DebugMode)
		}

		if viper.GetBool("verbose") {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		serv := api.FromConfig()

		// This is gross and needs to be reworked....
		serv.RegisterEndpoint("POST", "/api/v1/login", "", false, api.LoginPOST)
		serv.RegisterEndpoint("POST", "/api/v1/register", "", false, api.RegisterPOST)
		serv.RegisterEndpoint("GET", "/api/v1/reset", "read:profile", true, api.CardGET)

		serv.RegisterEndpoint("GET", "/api/v1/user", "read:user", true, api.UserGET)
		serv.RegisterEndpoint("DELETE", "/api/v1/user", "write:user", true, api.UserDELETE)

		serv.RegisterEndpoint("GET", "/api/v1/card", "read:card.wotc", true, api.CardGET)
		serv.RegisterEndpoint("POST", "/api/v1/card", "write:card.user", true, api.CardPOST)
		serv.RegisterEndpoint("DELETE", "/api/v1/card", "write:card.user", true, api.CardDELETE)

		serv.RegisterEndpoint("GET", "/api/v1/deck", "read:deck.wotc", true, api.DeckGET)
		serv.RegisterEndpoint("POST", "/api/v1/deck", "write:deck.user", true, api.DeckPOST)
		serv.RegisterEndpoint("DELETE", "/api/v1/deck", "write:deck.user", true, api.DeckDELETE)

		serv.RegisterEndpoint("GET", "/api/v1/deck/content", "read:deck.wotc", true, api.DeckContentGET)
		serv.RegisterEndpoint("POST", "/api/v1/deck/content", "write:deck.user", true, api.DeckContentPOST)
		serv.RegisterEndpoint("DELETE", "/api/v1/deck/content", "write:deck.user", true, api.DeckContentDELETE)

		serv.RegisterEndpoint("GET", "/api/v1/set", "read:set.wotc", true, api.SetGET)
		serv.RegisterEndpoint("POST", "/api/v1/set", "write:set.user", true, api.SetPOST)
		serv.RegisterEndpoint("DELETE", "/api/v1/set", "write:set.user", true, api.SetDELETE)

		serv.RegisterEndpoint("GET", "/api/v1/set/content", "read:set.wotc", true, api.SetContentGET)
		serv.RegisterEndpoint("POST", "/api/v1/set/content", "write:set.user", true, api.SetContentPOST)
		serv.RegisterEndpoint("DELETE", "/api/v1/set/content", "write:set.user", true, api.SetContentDELETE)

		err := serv.Run(viper.GetInt("port"))
		if err != nil {
			os.Exit(1) // not printing here as logic for logging is within the api.API struct
		}
	},
}

/*
Execute - Function automatically generated by cobra. Match CLI flags to there functionality. Should
not be called directly
*/
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
init - Function automatically created by cobra. Used to declare command line arguments, and functions that should be
executed when viper is initialized. Should not be called directly
*/
func init() {
	cobra.OnInitialize(initConfig)

	/*
		Universal CLI Flags - Any flags that can be used with any command
	*/
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mtgjson-api/config.json)")
	rootCmd.Flags().BoolP("debug", "d", false, "Put the gin engine in debug mode. (default is false [release mode])")
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbosity in logging (default is false)")
	rootCmd.Flags().IntP("port", "p", 8080, "The port the API should be exposed on (default is 8080)")

	/*
		MongoDB CLI Flags - Any flags used for identifying a MongoDB server
	*/
	rootCmd.Flags().String("mongo.hostname", "localhost", "The hostname of the MongoDB instance (default is localhost)")
	rootCmd.Flags().Int("mongo.port", 27017, "The port of the MongoDB instance (default is 27017)")
	rootCmd.Flags().String("mongo.username", "admin", "The username of the MongoDB user for authentication (default is admin)")
	rootCmd.Flags().String("mongo.password", "admin", "The hostname of the MongoDB instance (default is admin)")
	rootCmd.Flags().String("mongo.default_database", "mtgjson", "The MongoDB database to use by default (default is mtgjson)")

	/*
		Auth0 CLI Flags - Any flags used for identifying an Auth0 instance
	*/
	rootCmd.Flags().String("auth0.domain", "", "The domain of your Auth0 tenant")
	rootCmd.Flags().String("auth0.audience", "", "The audience of your Auth0 API")
	rootCmd.Flags().String("auth0.client_id", "", "The Client ID of your Auth0 application")
	rootCmd.Flags().String("auth0.client_secret", "", "The Client Secret of your Auth0 application")
	rootCmd.Flags().String("auth0.scope", "", "A space seperated string of Auth0 scopes that the API should recognize")

	/*
		Log CLI Flags - Any flags used for controlling slog logging features
	*/
	rootCmd.Flags().String("log.path", "/var/log/mtgjson-api", "The file path that log files should be saved to (default is /var/log/mtgjson-api)")

	/*
		Iterates through all the flags defined in rootCmd and binds them to viper values. The long name
		of the command is used by default
	*/
	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		fmt.Println("Error binding Cobra flags to viper: ", err.Error())
		fmt.Println("Viper config values may not work properly")
	}
}

/*
initConfig - Initialize viper with values from config files or environmental variables. Defaults
are not set here as CLI arguments are bound to viper config values. These provide defaults. Should
not be called directly, automatically called as a part of viper's initialization stack
*/
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetConfigType("json")
		viper.AddConfigPath(home + "/.config/mtgjson-api/")
		viper.SetConfigName("config.json")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}
}

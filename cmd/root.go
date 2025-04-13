package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "mtgjson-api",
	Short: "",
	Long: `A command line tool for controlling the execution of the MTGJSON RESTful API.

MTGJSON API is unofficial Fan Content permitted under the Fan Content Policy. 
Not approved/endorsed by Wizards of the Coast. Portions of the materials used are property of Wizards of the Coast. 
© Wizards of the Coast LLC.

MTGJSON API not officially endorsed by MTGJSON

Executing this binary with no CLI arguments will start the API with the default settings.
Any options can be configured with either command line arguments, a config file, or environment variables.

Developed and Tested on Debian-based Linux Distro's. Unconfirmed support on other operating systems`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mtgjson-api/config.json)")

	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbosity in logging")
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

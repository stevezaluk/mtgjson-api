package cmd

import (
	"fmt"
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mtgjson-api/config.json)")

	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbosity in logging")
}

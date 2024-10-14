package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var defaultConfig string = "~/.config/mtgjson/config.json"

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the API in the foreground and use STDOUT for logging",
	Long: `To start the API using the default config path:
$ mtgjson run

To start the API using a custom config file:
$ mtgjson run -c /path/to/config/.json

To start the API using environmental variables
$ mtgjson run --env`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[error] run not implemented yet")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("config", "c", defaultConfig, "The path to your MTGJSON config file")
	runCmd.Flags().BoolP("env", "e", false, "Ignore the default config path and attempt to use Environmental Variables")
	runCmd.Flags().Int64P("port", "p", 2100, "Set the default port that the API is listening on")
}

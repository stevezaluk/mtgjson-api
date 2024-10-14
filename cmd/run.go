package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the API in the foreground and use STDOUT for logging",
	Long: `To start the API using the default config path:
mtgjson run

To start the API using a custom config file:
mtgjson run -c /path/to/config/.json

To start the API using environmental variables
mtgjson run --env`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[error] run not implemented yet")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

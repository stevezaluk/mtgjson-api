package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "mtgjson",
	Short: "An RESTful API built on top of the MTGJSON data set",
	Long:  `This tool allows you to control the execution of the MTGJSON API`,
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

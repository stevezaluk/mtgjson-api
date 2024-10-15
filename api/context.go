package api

import (
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"mtgjson/server"
)

var ServerContext = context.Background()

func InitConfig(flags *pflag.FlagSet) {
	var config server.Config
	env, err := flags.GetBool("env")
	if err != nil {
		fmt.Println("[error] Error with env flag: ", err)
	}

	if env {
		config.ParseFromEnv()
	} else {
		configPath, err := flags.GetString("config")
		if err != nil {
			fmt.Println("[error] Error with config flag: ", err)
		}

		config.Parse(configPath)
	}

	ctx := context.WithValue(ServerContext, "config", config)
	ServerContext = ctx
}

func InitDatabase() {
	var database server.Database
	database.Config = ServerContext.Value("config").(server.Config)

	database.Connect() // externalize errors to here instead of printing

	ctx := context.WithValue(ServerContext, "database", database)
	ServerContext = ctx
}

func DestroyDatabase() {
	var database = ServerContext.Value("database").(server.Database)
	database.Disconnect()
}

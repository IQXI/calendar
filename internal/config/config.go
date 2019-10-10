package config

import (
	"github.com/spf13/viper"
	"log"
)

func GetConfig() map[string]interface{} {
	m := make(map[string]interface{})
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.AddConfigPath("./configs") // path to look for the config file in
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	m["outputs"] = viper.GetStringSlice("logger.outputs")
	m["logLevel"] = viper.GetString("logger.level")
	m["user"] = viper.GetString("db.user")
	m["password"] = viper.GetString("db.password")
	m["sslmode"] = viper.GetString("db.sslmode")
	m["host"] = viper.GetString("db.host")
	m["dbname"] = viper.GetString("db.dbname")

	return m
}

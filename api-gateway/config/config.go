package config

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path.Join(workDir, "config"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to load the config file: %s", err.Error())
	}
	switch viper.GetString("mode") {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	}
}

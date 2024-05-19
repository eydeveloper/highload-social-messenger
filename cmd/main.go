package main

import (
	"fmt"
	"github.com/eydeveloper/highload-social-messenger/internal/handler"
	"github.com/eydeveloper/highload-social-messenger/internal/service"

	messenger "github.com/eydeveloper/highload-social-messenger"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()
	
	services := service.NewService(redisClient)
	handlers := handler.NewHandler(services)

	err := new(messenger.Server).Run(viper.GetString("port"), handlers.InitRoutes())

	if err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

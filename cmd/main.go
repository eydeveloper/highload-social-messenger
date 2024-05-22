package main

import (
	"context"
	"github.com/eydeveloper/highload-social-messenger/internal/handler"
	"github.com/eydeveloper/highload-social-messenger/internal/service"

	messenger "github.com/eydeveloper/highload-social-messenger"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ctx = context.Background()

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:7001",
			"localhost:7002",
			"localhost:7003",
			"localhost:7004",
			"localhost:7005",
			"localhost:7006",
		},
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logrus.Fatalf("could not connected to redis cluster: %s", err.Error())
	}

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

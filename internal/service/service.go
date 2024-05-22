package service

import (
	"github.com/eydeveloper/highload-social-messenger/internal/entity"
	"github.com/redis/go-redis/v9"
)

type Messenger interface {
	SendMessage(message entity.Message) (int64, error)
	GetMessages(userId1, userId2 string) ([]string, error)
}

type Service struct {
	Messenger
}

func NewService(redisClient *redis.ClusterClient) *Service {
	return &Service{
		Messenger: NewMessengerService(redisClient),
	}
}

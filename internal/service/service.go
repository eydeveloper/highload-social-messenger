package service

import (
	"github.com/eydeveloper/highload-social-messenger/internal/entity"
	"github.com/redis/go-redis/v9"
)

type Messenger interface {
	SendMessage(message entity.Message) (int64, error)
	GetMessages(userID1, userID2 string) ([]string, error)
}

type Service struct {
	Messenger
}

func NewService(redisClient *redis.Client) *Service {
	return &Service{
		Messenger: NewMessengerService(redisClient),
	}
}

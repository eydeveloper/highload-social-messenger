package service

import (
	"context"
	"fmt"
	"github.com/eydeveloper/highload-social-messenger/internal/entity"
	"github.com/redis/go-redis/v9"
	"sort"
)

type MessengerService struct {
	redisClient *redis.Client
}

func NewMessengerService(redisClient *redis.Client) *MessengerService {
	return &MessengerService{
		redisClient: redisClient,
	}
}

func (s *MessengerService) SendMessage(message entity.Message) (int64, error) {
	ctx := context.Background()
	key := generateKey(message.SenderId, message.ReceiverId)
	id, err := s.redisClient.RPush(ctx, key, message.Text).Result()
	return id, err
}

func (s *MessengerService) GetMessages(userId1, userId2 string) ([]string, error) {
	ctx := context.Background()
	key := generateKey(userId1, userId2)
	return s.redisClient.LRange(ctx, key, 0, -1).Result()
}

func generateKey(userId1, userId2 string) string {
	ids := []string{userId1, userId2}
	sort.Strings(ids)
	return fmt.Sprintf("dialog:%s:%s", ids[0], ids[1])
}

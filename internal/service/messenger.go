package service

import (
	"context"
	"fmt"
	"github.com/eydeveloper/highload-social-messenger/internal/entity"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"sort"
)

type MessengerService struct {
	redisClient *redis.ClusterClient
}

func NewMessengerService(redisClient *redis.ClusterClient) *MessengerService {
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

func (s *MessengerService) RollbackSendMessage(id int64, receiverId string, senderId string) {
	ctx := context.Background()
	dialogId := generateKey(senderId, receiverId)

	placeholder := "__TO_DELETE__"
    err := s.redisClient.LSet(ctx, dialogId, -1, placeholder).Err()

    if err != nil {
		logrus.Error("error setting placeholder value: %w", err)
        return
    }

    err = s.redisClient.LRem(ctx, dialogId, 1, placeholder).Err()
    if err != nil {
		logrus.Error("error removing placeholder value: %w", err)
        return
    }
}

func generateKey(userId1, userId2 string) string {
	ids := []string{userId1, userId2}
	sort.Strings(ids)
	return fmt.Sprintf("%s:%s", ids[0], ids[1])
}

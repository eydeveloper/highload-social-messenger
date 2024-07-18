package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type CounterService struct {
	authToken string
}

func NewCounterService() *CounterService {
	return &CounterService{}
}

type counterResponse struct {
	Result bool `json:"result"`
}

func (s *CounterService) Increment(authToken string, requestId string, receiverId string, senderId string) error {
	url := "http://localhost:8002/api/counter/increment"
	requestBody := map[string]string{"receiver_id": receiverId, "dialog_id": generateKey(senderId, receiverId)}
	jsonRequestBody, _ := json.Marshal(requestBody)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequestBody))

	if err != nil {
		return err
	}

	request.Header.Set("X-Request-ID", requestId)
	request.Header.Set("Authorization", authToken)

	client := &http.Client{}
	response, err := client.Do(request)
	
	if err != nil {
		return err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)

	var counterResponse counterResponse

	if err = json.Unmarshal(responseBody, &counterResponse); err != nil {
		return err
	}

	return nil
}

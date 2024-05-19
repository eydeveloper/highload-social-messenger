package handler

import (
	"encoding/json"
	"github.com/eydeveloper/highload-social-messenger/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("api/messages")
	{
		api.POST("", h.authenticationMiddleware(), h.sendMessage)
		api.GET(":id", h.authenticationMiddleware(), h.getMessages)
	}

	return router
}

type AuthVerifyResponse struct {
	UserId string `json:"user_id"`
}

func (h *Handler) authenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		url := "http://localhost:8000/api/auth/verify"
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
			return
		}

		req.Header.Add("Authorization", authHeader)

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error sending request"})
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Received wrong status code"})
			return
		}

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error reading response body"})
			return
		}

		var response AuthVerifyResponse
		err = json.Unmarshal(body, &response)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing JSON response"})
			return
		}

		c.Set("userId", response.UserId)
		c.Next()
	}
}

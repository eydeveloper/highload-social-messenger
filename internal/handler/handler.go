package handler

import (
	"encoding/json"
	"github.com/eydeveloper/highload-social-messenger/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "highload_messenger_request_count",
			Help: "Total number of request to the messenger service",
		},
		[]string{"endpoint"},
	)
	requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "highload_messenger_request_duration_seconds",
            Help:    "Duration of messenger service requests in seconds",
        },
        []string{"endpoint"},
    )
    requestErrors = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "highload_messenger_request_errors",
            Help: "Total number of errors in messenger service requests",
        },
        []string{"endpoint"},
    )
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestErrors)
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := router.Group("api/messages")
	{
		api.POST("", h.authenticationMiddleware(), h.requestIdMiddleware(), trackMetrics(h.sendMessage))
		api.GET(":id", h.authenticationMiddleware(), h.requestIdMiddleware(), trackMetrics(h.getMessages))
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

func (h *Handler) requestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get("X-Request-ID")
		if requestId == "" {
			requestId = uuid.New().String()
		}

		c.Writer.Header().Set("X-Request-ID", requestId)
		c.Set("X-Request-ID", requestId)
		c.Next()
	}
}

func trackMetrics(next gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        endpoint := c.FullPath()

        // Continue to the next handler
        next(c)

        // After the handler is done
        duration := time.Since(start).Seconds()
        requestDuration.WithLabelValues(endpoint).Observe(duration)
        requestCount.WithLabelValues(endpoint).Inc()

        // Simulating error tracking for demonstration
        if rand.Float32() < 0.1 { // 10% chance of error
            requestErrors.WithLabelValues(endpoint).Inc()
        }

        // Ensure that the response is sent correctly
        if len(c.Errors) > 0 {
            c.JSON(-1, gin.H{"errors": c.Errors})
        }
    }
}

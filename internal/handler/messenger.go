package handler

import (
	"github.com/eydeveloper/highload-social-messenger/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) sendMessage(c *gin.Context) {
	senderId := c.MustGet("userId").(string)
	requestId := c.MustGet("X-Request-ID").(string)

	logrus.Info("Handling send message request with ID: " + requestId)

	var message entity.Message

	if err := c.BindJSON(&message); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	message.SenderId = senderId

	id, err := h.services.Messenger.SendMessage(message)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{"message_id": id})
}

func (h *Handler) getMessages(c *gin.Context) {
	senderId := c.MustGet("userId").(string)
	receiverId := c.Param("id")
	requestId := c.MustGet("X-Request-ID").(string)

	logrus.Info("Handling get messages request with ID: " + requestId)

	messages, err := h.services.GetMessages(senderId, receiverId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, messages)
}

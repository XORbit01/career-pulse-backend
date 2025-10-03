package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/XORbit01/jobseeker-backend/models"
	"github.com/XORbit01/jobseeker-backend/repos"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatRepo *repos.ChatRepository
}

func NewChatHandler(db *sql.DB) *ChatHandler {
	return &ChatHandler{
		chatRepo: repos.NewChatRepository(db),
	}
}

func RegisterChatRoutes(router *gin.RouterGroup, db *sql.DB) {
	h := NewChatHandler(db)
	chat := router.Group("/chats")
	{
		chat.GET("/", h.GetConversations)                     // List all conversations
		chat.POST("/:user_id/messages", h.SendMessage)        // Send message to user (creates conversation)
		chat.GET("/:conversation_id/messages", h.GetMessages) // Get messages in a conversation
		chat.PUT("/:conversation_id/read", h.MarkAsRead)      // Mark messages as read in conversation
	}
}

// GetConversations godoc
//
//	@Summary	Get all conversations for the current user
//	@Tags		Chat
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse{data=[]models.Conversation}
//	@Failure	401	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router		/chats/ [get]
//
// GetConversations godoc
//
//	@Summary	Get all conversations with message count
//	@Tags		Chat
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse{data=[]models.ConversationWithStats}
//	@Failure	401	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router		/chats/ [get]
func (h *ChatHandler) GetConversations(c *gin.Context) {
	userID := c.GetInt("userID")

	convs, count, err := h.chatRepo.GetConversationsWithStatsForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve conversations",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Conversations retrieved",
		Data: gin.H{
			"total":         count,
			"conversations": convs,
		},
	})
}

// SendMessage godoc
//
//	@Summary	Send a message to a user (creates/fetches conversation)
//	@Tags		Chat
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		user_id	path		int					true	"Receiver User ID"
//	@Param		input	body		models.MessageInput	true	"Message input"
//	@Success	200		{object}	models.SuccessResponse{data=models.Message}
//	@Failure	400		{object}	models.ErrorResponse
//	@Failure	401		{object}	models.ErrorResponse
//	@Failure	500		{object}	models.ErrorResponse
//	@Router		/chats/{user_id}/messages [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	senderID := c.GetInt("userID")
	receiverID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || receiverID <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   &models.ErrorInfo{Code: "INVALID_INPUT"},
		})
		return
	}

	var input models.MessageInput
	if err := c.ShouldBindJSON(&input); err != nil || input.Content == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid message input",
			Error:   &models.ErrorInfo{Code: "VALIDATION_ERROR", Details: err.Error()},
		})
		return
	}

	msg := models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    input.Content,
	}

	msgID, err := h.chatRepo.SaveMessage(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to send message",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	fullMsg, err := h.chatRepo.GetMessageByID(msgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch message",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Message sent successfully",
		Data:    fullMsg,
	})
}

// GetMessages godoc
//
//	@Summary	Get messages in a conversation
//	@Tags		Chat
//	@Security	BearerAuth
//	@Produce	json
//	@Param		conversation_id	path		int	true	"Conversation ID"
//	@Success	200				{object}	models.SuccessResponse{data=[]models.Message}
//	@Failure	400				{object}	models.ErrorResponse
//	@Failure	401				{object}	models.ErrorResponse
//	@Failure	500				{object}	models.ErrorResponse
//	@Router		/chats/{conversation_id}/messages [get]
func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID := c.GetInt("userID")
	convID, err := strconv.Atoi(c.Param("conversation_id"))
	if err != nil || convID <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid conversation ID",
			Error:   &models.ErrorInfo{Code: "INVALID_INPUT"},
		})
		return
	}

	messages, err := h.chatRepo.GetMessagesInConversation(convID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve messages",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Messages retrieved successfully",
		Data:    messages,
	})
}

// MarkAsRead godoc
//
//	@Summary	Mark all messages in a conversation as read
//	@Tags		Chat
//	@Security	BearerAuth
//	@Produce	json
//	@Param		conversation_id	path		int	true	"Conversation ID"
//	@Success	200				{object}	models.SuccessResponse
//	@Failure	400				{object}	models.ErrorResponse
//	@Failure	401				{object}	models.ErrorResponse
//	@Failure	500				{object}	models.ErrorResponse
//	@Router		/chats/{conversation_id}/read [put]
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetInt("userID")
	convID, err := strconv.Atoi(c.Param("conversation_id"))
	if err != nil || convID <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid conversation ID",
			Error:   &models.ErrorInfo{Code: "INVALID_INPUT"},
		})
		return
	}

	if err := h.chatRepo.MarkMessagesAsReadInConversation(convID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to mark messages as read",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Messages marked as read",
	})
}

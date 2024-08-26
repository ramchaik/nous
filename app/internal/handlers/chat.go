package handlers

import (
	"fmt"
	"net/http"

	"nous/internal/llmclient"
	"nous/internal/models"
	"nous/internal/store"

	"github.com/gin-gonic/gin"
)

type ChatHandler interface {
	Chat(*gin.Context)
	CreateChat(*gin.Context)
	GetChat(*gin.Context)
	// TODO: implement CRUD
	// UpdateChat(*gin.Context)
	// DeleteChat(*gin.Context)
}

type DefaultChatHandler struct {
	store     store.ChatStore
	llmClient llmclient.LLMClient
}

func NewChatHandler(store store.ChatStore, llmClient llmclient.LLMClient) ChatHandler {
	return &DefaultChatHandler{
		store:     store,
		llmClient: llmClient,
	}
}

func (h *DefaultChatHandler) Chat(c *gin.Context) {
	query := c.Query("query")

	predictResp, err := h.llmClient.Predict(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get prediction: " + err.Error()})
		return
	}

	fmt.Printf("Steps: %v\n", predictResp.Steps)

	c.HTML(http.StatusOK, "chat.html", gin.H{
		"query":    query,
		"response": predictResp.Response,
	})
}

func (h *DefaultChatHandler) CreateChat(c *gin.Context) {
	var chat models.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.Create(&chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

func (h *DefaultChatHandler) GetChat(c *gin.Context) {
	chatID := c.Param("id")

	chat, err := h.store.GetByID(chatID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// Implement UpdateChat and DeleteChat methods

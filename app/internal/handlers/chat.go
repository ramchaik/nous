package handlers

import (
	"net/http"

	"nous/internal/llmclient"
	"nous/internal/models"
	"nous/internal/store"

	"github.com/gin-gonic/gin"
)

type ChatHandler interface {
	CreateChat(*gin.Context)
	GetChat(*gin.Context)
	Predict(*gin.Context)
	// TODO: implement CRUD
	// UpdateChat(*gin.Context)
	// DeleteChat(*gin.Context)
}

type ChatAPIHandler struct {
	store     store.ChatStore
	llmClient llmclient.LLMClient
}

func NewChatAPIHandler(store store.ChatStore, llmClient llmclient.LLMClient) ChatHandler {
	return &ChatAPIHandler{store: store, llmClient: llmClient}
}

func (h *ChatAPIHandler) CreateChat(c *gin.Context) {
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

func (h *ChatAPIHandler) GetChat(c *gin.Context) {
	chatID := c.Param("id")

	chat, err := h.store.GetByID(chatID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *ChatAPIHandler) Predict(c *gin.Context) {
	var request struct {
		Query string `json:"query" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	predictResp, err := h.llmClient.Predict(c.Request.Context(), request.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get prediction: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, predictResp)
}

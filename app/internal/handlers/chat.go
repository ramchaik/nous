package handlers

import (
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

type ChatAPIHandler struct {
	store     store.ChatStore
	llmClient llmclient.LLMClient
}

type ChatUIHandler struct {
	llmClient llmclient.LLMClient
}

func NewChatUIHandler(llmClient llmclient.LLMClient) *ChatUIHandler {
	return &ChatUIHandler{llmClient: llmClient}
}

func NewChatAPIHandler(store store.ChatStore, llmClient llmclient.LLMClient) *ChatAPIHandler {
	return &ChatAPIHandler{store: store, llmClient: llmClient}
}

func (h *ChatUIHandler) RenderChatPage(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.HTML(http.StatusOK, "chat.html", gin.H{})
		return
	}

	predictResp, err := h.llmClient.Predict(query)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to get prediction: " + err.Error()})
		return
	}

	c.HTML(http.StatusOK, "chat.html", gin.H{
		"userMessage": query,
		"botResponse": predictResp.Response,
	})
}

func (h *ChatUIHandler) HandleChatMessage(c *gin.Context) {
	query := c.PostForm("query")
	if query == "" {
		c.String(http.StatusBadRequest, "Query is required")
		return
	}

	predictResp, err := h.llmClient.Predict(query)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get prediction: "+err.Error())
		return
	}

	c.HTML(http.StatusOK, "chat_messages.html", gin.H{
		"userMessage": query,
		"botResponse": predictResp.Response,
	})
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

func (h *ChatAPIHandler) PredictResponse(c *gin.Context) {
	var request struct {
		Query string `json:"query" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	predictResp, err := h.llmClient.Predict(request.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get prediction: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, predictResp)
}

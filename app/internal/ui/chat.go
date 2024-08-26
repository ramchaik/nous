package ui

import (
	"net/http"

	"nous/internal/llmclient"

	"github.com/gin-gonic/gin"
)

type ChatHandler interface {
	RenderChatPage(*gin.Context)
	HandleChatMessage(*gin.Context)
}

type ChatUIHandler struct {
	llmClient llmclient.LLMClient
}

func NewChatUIHandler(llmClient llmclient.LLMClient) ChatHandler {
	return &ChatUIHandler{llmClient: llmClient}
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

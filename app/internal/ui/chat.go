package ui

import (
	"log"
	"net/http"
	"net/url"

	"nous/internal/llmclient"
	"nous/internal/models"
	"nous/internal/store"
	"nous/internal/ui/components"

	"github.com/gin-gonic/gin"
)

type ChatHandler interface {
	InitiateChat(*gin.Context)
	RenderChatPage(*gin.Context)
	HandleChatMessage(*gin.Context)
}

type ChatUIHandler struct {
	llmClient llmclient.LLMClient
	chatStore store.ChatStore
}

func NewChatUIHandler(chatStore store.ChatStore, llmClient llmclient.LLMClient) ChatHandler {
	return &ChatUIHandler{
		chatStore: chatStore,
		llmClient: llmClient,
	}
}

func (h *ChatUIHandler) InitiateChat(c *gin.Context) {
	sessionID, err := h.getOrCreateSession(c)
	if err != nil {
		components.ErrorMessage("Failed to create session: "+err.Error()).Render(c, c.Writer)
		return
	}

	query := c.Query("query")
	sessionIDEncoded := url.QueryEscape(sessionID)
	urlEncoededQuery := url.QueryEscape(query)
	chatID := store.GenerateUUID()

	c.Redirect(http.StatusFound, "/chat/"+chatID+"?sid="+sessionIDEncoded+"&query="+urlEncoededQuery)
}

func (h *ChatUIHandler) RenderChatPage(c *gin.Context) {
	chatID := c.Param("chat_id")
	sessionID := c.Query("sid")
	query := c.Query("query")

	if sessionID == "" {
		components.ErrorMessage("Session ID is required").Render(c, c.Writer)
		return
	}

	chats, err := h.chatStore.GetChatsBySession(sessionID)
	if err != nil {
		components.ErrorMessage("Failed to get chat history: "+err.Error()).Render(c, c.Writer)
		return
	}

	components.ChatPage(sessionID, chatID, chats, query).Render(c, c.Writer)
}

func (h *ChatUIHandler) HandleChatMessage(c *gin.Context) {
	chatID := c.Param("chat_id")
	sessionID := c.PostForm("sid")
	query := c.PostForm("query")

	if query == "" || sessionID == "" {
		components.ErrorMessage("Query and sessionID are required").Render(c, c.Writer)
		return
	}

	// Save user message
	userChat := &models.Chat{
		ChatID:    chatID,
		SessionID: sessionID,
		Text:      query,
		Type:      "user",
	}
	err := h.chatStore.CreateChat(userChat)
	if err != nil {
		log.Println("Failed to save user message:", err)
		components.ErrorMessage("Failed to save user message:"+err.Error()).Render(c, c.Writer)
		return
	}

	// Get chat history
	chatHistory, err := h.chatStore.GetChatsBySession(sessionID)
	if err != nil {
		log.Println("Failed to get chat history:", err)
		components.ErrorMessage("Failed to get chat history:"+err.Error()).Render(c, c.Writer)
		return
	}

	// Prepare chat history for LLM
	var llmChatHistory []llmclient.ChatMessage
	for _, chat := range chatHistory {
		llmChatHistory = append(llmChatHistory, llmclient.ChatMessage{
			Role:    chat.Type,
			Content: chat.Text,
		})
	}

	// Get prediction from LLM
	predictResp, err := h.llmClient.Predict(c.Request.Context(), query, llmChatHistory)
	if err != nil {
		components.ErrorMessage("Failed to get prediction:"+err.Error()).Render(c, c.Writer)
		return
	}

	// Save bot response
	botChat := &models.Chat{
		ChatID:    chatID,
		SessionID: sessionID,
		Text:      predictResp.Response,
		Type:      "agent",
	}
	err = h.chatStore.CreateChat(botChat)
	if err != nil {
		components.ErrorMessage("Failed to save bot message:"+err.Error()).Render(c, c.Writer)
		return
	}

	components.ChatMessages(query, predictResp.Response).Render(c, c.Writer)
}

func (h *ChatUIHandler) getOrCreateSession(c *gin.Context) (string, error) {
	sessionID, err := c.Cookie("session_id")
	if err == nil && sessionID != "" {
		// Session exists, verify it
		_, err := h.chatStore.GetSession(sessionID)
		if err == nil {
			return sessionID, nil
		}
	}

	// Create new session
	sessionID, err = h.chatStore.CreateSession()
	if err != nil {
		return "", err
	}

	// Set session cookie
	c.SetCookie("session_id", sessionID, 3600*24*30, "/", "", false, true)
	return sessionID, nil
}

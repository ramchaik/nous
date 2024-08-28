package server

import (
	"log"
	"net/http"
	"nous/internal/config"
	"nous/internal/database"
	"nous/internal/handlers"
	"nous/internal/llmclient"
	"nous/internal/store"
	"nous/internal/ui"

	"github.com/gin-gonic/gin"
)

type Server interface {
	SetupRoutes()
	Run(string) error
}

type DefaultServer struct {
	router    *gin.Engine
	config    *config.Config
	db        database.Database
	llmClient llmclient.LLMClient
}

func New(cfg *config.Config, db database.Database, llmClient llmclient.LLMClient) Server {
	s := &DefaultServer{
		router:    gin.Default(),
		config:    cfg,
		db:        db,
		llmClient: llmClient,
	}

	s.SetupRoutes()

	return s
}

func (s *DefaultServer) SetupRoutes() {
	s.router.Use(s.globalErrorHandler())

	s.router.Static("/static", s.config.StaticPath)
	s.router.LoadHTMLGlob(s.config.TemplatesPath)

	chatStore := store.NewChatStore(s.db.GetDB())
	chatAPIHandler := handlers.NewChatAPIHandler(chatStore, s.llmClient)

	chatUIHandler := ui.NewChatUIHandler(chatStore, s.llmClient)
	homeUIHandler := ui.NewHomeUIHandler()

	// UI routes
	s.router.GET("/", homeUIHandler.RenderHomePage)
	s.router.GET("/chat", chatUIHandler.InitiateChat)
	s.router.GET("/chat/:chat_id", chatUIHandler.RenderChatPage)
	s.router.POST("/chat/:chat_id", chatUIHandler.HandleChatMessage)

	// API routes
	api := s.router.Group("/api")
	{
		api.POST("/chats", chatAPIHandler.CreateChat)
		api.GET("/chats/:id", chatAPIHandler.GetChat)
		api.POST("/predict", chatAPIHandler.Predict)
	}
}

func (s *DefaultServer) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *DefaultServer) globalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error (you can use a proper logging library here)
				log.Printf("Panic: %v", err)

				// Determine if the request expects HTML or JSON
				if c.Request.Header.Get("Accept") == "application/json" {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
				} else {
					c.HTML(http.StatusInternalServerError, "error.html", gin.H{
						"error": "An unexpected error occurred. Please try again later.",
					})
				}

				// Abort the request
				c.Abort()
			}
		}()

		// Process the request
		c.Next()
	}
}

package server

import (
	"nous/internal/config"
	"nous/internal/database"
	"nous/internal/handlers"
	"nous/internal/llmclient"
	"nous/internal/store"

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
	s.router.Static("/static", s.config.StaticPath)
	s.router.LoadHTMLGlob(s.config.TemplatesPath)

	chatStore := store.NewChatStore(s.db.GetDB())
	chatUIHandler := handlers.NewChatUIHandler(s.llmClient)
	chatAPIHandler := handlers.NewChatAPIHandler(chatStore, s.llmClient)

	// UI routes
	s.router.GET("/", handlers.Home)
	s.router.GET("/chat", chatUIHandler.RenderChatPage)
	s.router.POST("/chat", chatUIHandler.HandleChatMessage)

	// API routes
	api := s.router.Group("/api")
	{
		api.POST("/chats", chatAPIHandler.CreateChat)
		api.GET("/chats/:id", chatAPIHandler.GetChat)
		api.POST("/predict", chatAPIHandler.PredictResponse)
	}
}

func (s *DefaultServer) Run(addr string) error {
	return s.router.Run(addr)
}

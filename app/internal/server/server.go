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

	s.router.GET("/", handlers.Home)

	chatStore := store.NewChatStore(s.db.GetDB())
	chatHandler := handlers.NewChatHandler(chatStore, s.llmClient)

	s.router.GET("/chat", chatHandler.Chat)
	s.router.POST("/chat", chatHandler.CreateChat)
	s.router.GET("/chat/:id", chatHandler.GetChat)
}

func (s *DefaultServer) Run(addr string) error {
	return s.router.Run(addr)
}

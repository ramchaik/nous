package server

import (
	"nous/internal/config"
	"nous/internal/handlers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func New(cfg *config.Config) *Server {
	s := &Server{
		router: gin.Default(),
		config: cfg,
	}

	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	s.router.Static("/static", s.config.StaticPath)
	s.router.LoadHTMLGlob(s.config.TemplatesPath)

	s.router.GET("/", handlers.Home)
	s.router.GET("/chat", handlers.Chat)
}

func (s *Server) Run() error {
	return s.router.Run(s.config.Port)
}

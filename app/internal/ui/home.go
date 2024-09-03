package ui

import (
	"nous/internal/ui/components"

	"github.com/gin-gonic/gin"
)

type HomeHandler interface {
	RenderHomePage(*gin.Context)
}

type HomeUIHandler struct{}

func NewHomeUIHandler() HomeHandler {
	return &HomeUIHandler{}
}

func (h *HomeUIHandler) RenderHomePage(c *gin.Context) {
	component := components.HomePage()
	component.Render(c, c.Writer)
}

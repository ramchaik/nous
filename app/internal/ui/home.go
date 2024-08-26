package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeHandler interface {
	RenderHomePage(*gin.Context)
}

type HomeUIHandler struct {
}

func NewHomeUIHandler() HomeHandler {
	return &HomeUIHandler{}
}

func (h *HomeUIHandler) RenderHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Welcom to Nous!",
	})
}

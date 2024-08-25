package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Chat(c *gin.Context) {
	query := c.Query("query")
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"query": query,
	})
}

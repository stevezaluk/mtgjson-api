package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

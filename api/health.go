package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
Gin handler for GET request to the health endpoint. This should not be called directly and
should only be passed to the gin router. Currently this just returns healthy
*/
func HealthGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

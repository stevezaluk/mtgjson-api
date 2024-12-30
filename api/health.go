package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
HealthGET Gin handler for the GET request to the Health Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func HealthGET(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

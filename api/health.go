package api

import (
	"github.com/gin-gonic/gin"
)

func HealthGET(ctx *gin.Context) {
	ctx.HTML(200, "https", gin.H{"status": "Healthy"})
}

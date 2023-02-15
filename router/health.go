package router

import (
	"github.com/gin-gonic/gin"
	"github.com/romberli/go-template-web/api/v1/health"
)

// RegisterHealth is the sub-router for health
func RegisterHealth(group *gin.RouterGroup) {
	healthGroup := group.Group("/health")
	{
		healthGroup.GET("/status", health.Status)
		healthGroup.POST("/ping", health.Ping)
	}
}

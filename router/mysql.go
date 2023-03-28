package router

import (
	"github.com/gin-gonic/gin"
	"github.com/romberli/db-operator/api/v1/mysql"
)

// RegisterMySQL is the sub-router for mysql
func RegisterMySQL(group *gin.RouterGroup) {
	mysqlGroup := group.Group("/mysql")
	{
		mysqlGroup.POST("/install", mysql.Install)
	}
}

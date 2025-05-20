package http

import "github.com/gin-gonic/gin"

func Run() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	InitRouter(v1)
	router.Run("0.0.0.0:10011")
}

package main

import (
	"task-manager/config"
	"task-manager/internal/routers"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
	config.SyncDB()
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World, it's Task Management System",
		})
	})
	routers.TaskRouter(r)
	routers.UserRouter(r)
	r.Run()
}

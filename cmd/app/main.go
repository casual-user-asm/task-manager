package main

import (
	"task-manager/config"
	"task-manager/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Connect()

	r := gin.Default()

	r.POST("/create", handlers.CreateTask)
	r.GET("/tasks", handlers.GetTasks)

	r.Run(":8080")
}

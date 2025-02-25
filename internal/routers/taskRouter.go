package routers

import (
	"task-manager/internal/handlers"
	"task-manager/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func TaskRouter(c *gin.Engine) {
	task := c.Group("/task")
	{
		task.POST("/create", middlewares.AuthMiddleware, handlers.CreateTask)
		task.GET("/", middlewares.AuthMiddleware, handlers.GetTasks)
		task.DELETE("/delete", middlewares.AuthMiddleware, handlers.DeleteTask)
		task.PUT("/update", middlewares.AuthMiddleware, handlers.UpdateTasks)
	}
}

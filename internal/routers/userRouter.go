package routers

import (
	"task-manager/internal/handlers"
	"task-manager/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRouter(c *gin.Engine) {
	user := c.Group("/user")
	{
		user.POST("/register", handlers.UserRegistration)
		user.POST("/login", handlers.UserLogin)
		user.PUT("/logout", middlewares.AuthMiddleware, handlers.UserLogout)
		user.DELETE("/delete", middlewares.AuthMiddleware, handlers.UserDelete)
	}
}

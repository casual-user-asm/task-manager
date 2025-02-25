package handlers

import (
	"net/http"
	"task-manager/config"
	"task-manager/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Find user by ID
	var user models.User
	err := config.DB.First(&user, userID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Obtain data from request
	var body struct {
		Title       string `json:"title" binding:"required"`
		Description string `'json:"description" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fields are empty",
		})
		return
	}

	date := time.Now()

	newTask := models.Task{
		Title:       body.Title,
		Description: body.Description,
		CreatedBy:   userID,
		Date:        date,
	}

	result := config.DB.Create(&newTask)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task created successfully",
	})
}

func GetTasks(c *gin.Context) {
	task_id := c.Query("task_id")

	var tasks []models.Task
	err := config.DB.Find(&tasks, task_id).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Can't find tasks",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func UpdateTasks(c *gin.Context) {
	task_id := c.Query("task_id")

	var task models.Task
	err := config.DB.First(&task, task_id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fields are empty",
		})
		return
	}

	if body.Title != "" {
		task.Title = body.Title
	}

	if body.Description != "" {
		task.Title = body.Description
	}

	err = config.DB.Save(&task).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"task":    task,
	})
}

func DeleteTask(c *gin.Context) {
	task_id := c.Query("task_id")

	var task models.Task
	err := config.DB.First(&task, task_id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	err = config.DB.Delete(&models.Task{}, task_id).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error deleting task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

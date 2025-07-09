package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func MockCreateTask(c *gin.Context) {
	var body struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields are empty or invalid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task created successfully",
	})
}

func MockUpdateTask(c *gin.Context) {
	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Simulate non-existent task
	if taskID == "999" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields are empty"})
		return
	}

	// Check for empty fields
	if body.Title == "" && body.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields are empty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
	})
}

func TestCreateTask(t *testing.T) {
	router := setupTestRouter()
	router.POST("/create", MockCreateTask)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		shouldContain  string
	}{
		{
			name: "Valid Task Creation",
			requestBody: map[string]interface{}{
				"title":       "Test Task",
				"description": "This is a test task",
			},
			expectedStatus: http.StatusOK,
			shouldContain:  "Task created successfully",
		},
		{
			name: "Invalid Task Creation - Empty Fields",
			requestBody: map[string]interface{}{
				"title":       "",
				"description": "",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty or invalid",
		},
		{
			name: "Invalid Task Creation - Missing Fields",
			requestBody: map[string]interface{}{
				"title":       "Test Task",
				"description": "",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty or invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createJSONRequest("POST", "/create", tt.requestBody)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.shouldContain != "" {
				var response map[string]interface{}
				err = json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)

				if response["message"] != nil {
					assert.Contains(t, response["message"], tt.shouldContain)
				} else if response["error"] != nil {
					assert.Contains(t, response["error"], tt.shouldContain)
				}
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	router := setupTestRouter()
	router.PUT("/update", MockUpdateTask)

	tests := []struct {
		name           string
		taskID         string
		requestBody    map[string]interface{}
		expectedStatus int
		shouldContain  string
	}{
		{
			name:   "Valid Task Update",
			taskID: "1",
			requestBody: map[string]interface{}{
				"title":       "Updated Task",
				"description": "This is an updated task",
			},
			expectedStatus: http.StatusOK,
			shouldContain:  "Task updated successfully",
		},
		{
			name:   "Invalid Task Update - Empty Fields",
			taskID: "1",
			requestBody: map[string]interface{}{
				"title":       "",
				"description": "",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name:   "Invalid Task Update - Non-existent Task",
			taskID: "",
			requestBody: map[string]interface{}{
				"title":       "Non-existent Task",
				"description": "This task does not exist",
			},
			expectedStatus: http.StatusNotFound,
			shouldContain:  "Task not found",
		},
		{
			name:   "Invalid Task Update - Non-existent Task ID",
			taskID: "999",
			requestBody: map[string]interface{}{
				"title":       "Non-existent Task",
				"description": "This task does not exist",
			},
			expectedStatus: http.StatusNotFound,
			shouldContain:  "Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build URL with query parameter
			url := "/update"
			if tt.taskID != "" {
				url += "?task_id=" + tt.taskID
			}

			req, err := createJSONRequest("PUT", url, tt.requestBody)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.shouldContain != "" {
				var response map[string]interface{}
				err = json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)

				if response["message"] != nil {
					assert.Contains(t, response["message"], tt.shouldContain)
				} else if response["error"] != nil {
					assert.Contains(t, response["error"], tt.shouldContain)
				}
			}
		})
	}
}

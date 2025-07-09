package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func MockUserRegistration(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fields are empty",
		})
		return
	}

	// Simulate existing user check
	if body.Email == "existing@example.com" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User alreay exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func MockUserLogin(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fields are empty",
		})
		return
	}

	// Simulate user not found
	if body.Email == "nonexistent@example.com" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	// Simulate invalid password
	if body.Password == "wrongpassword" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfull",
		"token":   "mock-jwt-token",
	})
}

func MockUserLogout(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successfull",
	})
}

func MockUserDelete(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func TestUserRegistration(t *testing.T) {
	router := setupTestRouter()
	router.POST("/register", MockUserRegistration)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		shouldContain  string
	}{
		{
			name: "Valid Registration Data",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			shouldContain:  "",
		},
		{
			name: "Missing Username",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name: "Missing Email",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name: "Missing Password",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name: "Empty Username",
			requestBody: map[string]interface{}{
				"username": "",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name: "Empty Email",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name: "Empty Password",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name: "Existing User",
			requestBody: map[string]interface{}{
				"username": "existinguser",
				"email":    "existing@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "User alreay exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createJSONRequest("POST", "/register", tt.requestBody)
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

func TestUserLogin(t *testing.T) {
	router := setupTestRouter()
	router.POST("/login", MockUserLogin)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		shouldContain  string
	}{
		{
			name: "Valid Login Data",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			shouldContain:  "",
		},
		{
			name: "Missing Email",
			requestBody: map[string]interface{}{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name: "Missing Password",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name:           "Empty Body",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
		{
			name: "Empty Email",
			requestBody: map[string]interface{}{
				"email":    "",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Fields are empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createJSONRequest("POST", "/login", tt.requestBody)
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

func TestUserLogout(t *testing.T) {
	router := setupTestRouter()
	router.PUT("/logout", MockUserLogout)

	tests := []struct {
		name           string
		expectedStatus int
		shouldContain  string
	}{
		{
			name:           "Valid Logout",
			expectedStatus: http.StatusOK,
			shouldContain:  "Logout successfull",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createJSONRequest("PUT", "/logout", nil)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.shouldContain != "" {
				var response map[string]interface{}
				err = json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["message"], tt.shouldContain)
			}

			// Check that JWT cookie is cleared
			cookies := recorder.Result().Cookies()
			for _, cookie := range cookies {
				if cookie.Name == "jwt" {
					assert.Equal(t, "", cookie.Value)
					assert.Equal(t, -1, cookie.MaxAge)
				}
			}
		})
	}
}

func TestUserDelete(t *testing.T) {
	router := setupTestRouter()
	router.DELETE("/delete", MockUserDelete)

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "Delete User",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createJSONRequest("DELETE", "/delete", nil)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			var response map[string]interface{}
			err = json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")
		})
	}
}

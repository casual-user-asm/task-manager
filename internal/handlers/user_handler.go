package handlers

import (
	"net/http"
	"os"
	"task-manager/config"
	"task-manager/internal/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserRegistration(c *gin.Context) {
	// get fields from request
	var body struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fields are empty",
		})
		return
	}

	// Check if user already exist
	var user models.User
	config.DB.First(&user, "email = ?", body.Email)
	if user.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User alreay exist",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error hashing the password",
		})
		return
	}

	// Create user
	newUser := models.User{Username: body.Username, Email: body.Email, Password: string(hash)}
	result := config.DB.Create(&newUser)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating user",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func UserLogin(c *gin.Context) {
	// Get fields from request
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fields are empty",
		})
		return
	}

	// Find user by email
	var user models.User
	config.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erorr": "User not found",
		})
		return
	}

	// Compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and complete token using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error generating JWT",
		})
		return
	}

	// Respond
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfull",
		"token":   tokenString,
	})
}

//	url -X POST http://localhost:8080/user/login      -H "Content-Type: application/json"      -d '{
//		"email": "test@example.com",
//		"password": "securepassword"
//	  }'
func UserLogout(c *gin.Context) {
	// Clear JWT cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successfull",
	})
}

func UserDelete(c *gin.Context) {
	user_id := c.GetUint("user_id")

	// Clear JWT cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", "", -1, "", "", false, true)

	// Delete user
	err := config.DB.Delete(&models.User{}, user_id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

package user_r

import (
	"errors"
	"fmt"
	models "gocode/backend/backend/text-to-picture/models/init"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User struct to represent user data
type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_verified"`
}

// GetUserByEmail retrieves a user by name from the database
func GetUserByName(db *gorm.DB, name string) (*User, error) {
	var user User
	result := db.Where("name = ?", name).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", result.Error)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email from the database
func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", result.Error)
	}
	return &user, nil
}

func GetUserInfo(c *gin.Context) {
	name := c.Query("name")
	email := c.Query("email")
	var user *User
	var err error

	if name != "" {
		user, err = GetUserByName(models.DB, name)
	} else if email != "" {
		user, err = GetUserByEmail(models.DB, email)
	}

	if err != nil {
		if err.Error() == "User not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

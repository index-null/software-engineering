package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// User struct to represent user data
type User struct {
	ID         int
	Name       string
	Email      string
	Password   string
	IsVerified bool
}

var db *sql.DB

func initDB() {
	var err error
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Database connection successful")
}


func getUser(c *gin.Context) {
	username := c.Query("username")

	var user User
	err := db.QueryRow("SELECT id, username, email FROM users WHERE username = ?", username).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	r.GET("/user", getUser)

	r.Run(":8080")
}
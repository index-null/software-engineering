package auth

import (
	"net/http"

	db "text-to-picture/models/init"           // 给 init 包设置别名为 db
	"text-to-picture/models/repository/user_r" // 本地引用插入和查询函数
	"text-to-picture/models/user"              // 本地引用用户模型

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 注册处理函数
func Register(c *gin.Context) {
	var newUser user.Register // 使用 Register 结构体作为注册用户
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 调用插入函数
	err := user_r.InsertUserInformation(db.DB, &newUser) // 使用 db.DB 作为数据库连接
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// 登录处理函数
func Login(c *gin.Context) {
	var loginUser user.UserInformation // 使用 UserInformation 结构体作为登录用户
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 查询用户信息
	dbUser, err := user_r.GetUserByEmail(db.DB, loginUser.Email) // 使用 db.DB 作为数据库连接
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user"})
		}
		return
	}

	// 验证密码
	if dbUser.Password != loginUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	// 返回登录成功
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

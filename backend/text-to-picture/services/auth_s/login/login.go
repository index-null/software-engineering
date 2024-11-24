package auth_s

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	middlewire "text-to-picture/middlewire/jwt"
	models "text-to-picture/models/init"
	"text-to-picture/models/repository/user_r"
	userLogin "text-to-picture/models/user"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

// 注册
func Register(c *gin.Context) {
	// 定义用于接收 JSON 数据的结构体
	var input userLogin.UserInformation

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求数据格式错误"})
		return
	}

	//插入数据
	if err := user_r.InsertUserInformation(models.DB, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "用户创建失败",
			"error":   err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "注册成功",
	})
}

// 登录
func Login(c *gin.Context) {
	// 定义用于接收 JSON 数据的结构体
	var input struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求数据格式错误"})
		return
	}

	// 查找用户
	user, err := user_r.GetUserByName(models.DB, input.Name)
	if err != nil {
		log.Printf("name:%v,%v", input.Name, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "用户不存在"})
			return
		}
		//
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "数据库查询错误"})
		return
	}

	// 验证密码将用户的明文密码与数据库中的哈希密码对比
	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "密码错误"})
		return
	}

	//生成 JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middlewire.Claims{
		Username: input.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middlewire.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "生成 token 错误"})
		return
	}

	// 登录成功
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登录成功",
		"token":   tokenString,
	})
}

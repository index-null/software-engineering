package middlewire

import (
	//"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JwtKey 用于加密的密钥，换成你想要的秘钥
var JwtKey = []byte("wujinhao123")

// Claims 结构体定义了JWT中包含的用户信息
// Username 用户名
// StandardClaims 提供了JWT标准中定义的通用声明
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// JWTAuthMiddleware 是一个中间件，用于验证请求中的JWT令牌
// 该中间件检查传入请求的Authorization头，验证JWT令牌的有效性
// 如果令牌有效，将用户名设置到上下文中并允许请求继续
// 如果令牌无效或缺失，返回401状态码和错误信息
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的Authorization字段值
		tokenStr := c.GetHeader("Authorization")
		// 如果没有提供Token，则返回401状态码和错误信息
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请求头中缺少Token"})
			c.Abort()
			return
		}
		// 初始化Claims对象，用于解析Token中的用户信息
		claims := &Claims{}
		// 使用密钥解析Token，并验证其有效性
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil // jwtKey 是你的签名密钥
		})
		// 如果解析出错或Token无效，则返回401状态码和错误信息
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的Token"})
			c.Abort()
			return
		}
		// 如果Token有效，将用户名设置到上下文中，以便后续的处理函数使用
		c.Set("username", claims.Username)
		// 允许请求继续到下一个处理函数
		c.Next()
	}
}

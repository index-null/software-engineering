package main

import (
	"log"
	"text-to-picture/services/generate_s"

	"text-to-picture/api/auth"       // 导入注册和登录路由
	db "text-to-picture/models/init" // 为 init 包设置别名为 db

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	if err := db.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置路由
	r := gin.Default()

	// 配置 CORS 中间件
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // 允许的源，可以根据需要修改
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(config))

	r.POST("/register", auth.Register) // 注册路由
	r.POST("/login", auth.Login)       // 登录路由
	r.POST("/generate", generate_s.ReturnImage)
	r.Run("0.0.0.0:8080")
}

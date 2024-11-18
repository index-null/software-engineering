package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"text-to-picture/api/auth"
	"text-to-picture/api/generate"
	db "text-to-picture/models/init"
)

type DBConfig struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

func main() {

	//读取DBConfig.yaml文件
	yamlFile, err := os.ReadFile("config\\DBconfig\\DBconfig.yaml")
	if err != nil {
		fmt.Println("Error reading config.yaml file: %v", err)
	}

	//复制到config结构体
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		fmt.Println("Error parsing config.yaml file: %v", err)
	}

	//设置数据库连接的环境变量
	os.Setenv("DB_USER", dbconfig.DB.User)
	os.Setenv("DB_PASSWORD", dbconfig.DB.Password)
	os.Setenv("DB_HOST", dbconfig.DB.Host)
	os.Setenv("DB_PORT", dbconfig.DB.Port)
	os.Setenv("DB_NAME", dbconfig.DB.Name)

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

	// 创建 ImageGenerator 实例
	imgGen := generate.NewImageGenerator()

	// 注册路由
	r.POST("/register", auth.Register) // 注册路由
	r.POST("/login", auth.Login)       // 登录路由
	r.POST("/generate", func(c *gin.Context) {
		imgGen.ReturnImage(c)
	})

	// 添加静态文件服务，指向 docs 目录
	r.Static("/docs", "./docs")

	// 注册 Swagger 路由，并指定 doc.json 文件的路径
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	r.Run("0.0.0.0:8080")
}

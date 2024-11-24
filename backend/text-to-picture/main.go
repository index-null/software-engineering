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
	"text-to-picture/api/generate"
	middlewire "text-to-picture/middlewire/jwt"
	db "text-to-picture/models/init"
	image_r "text-to-picture/models/repository/image_r"
	user_r "text-to-picture/models/repository/user_r"
	"text-to-picture/services/auth_s/avator"
	auth_s "text-to-picture/services/auth_s/login"
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
	yamlFile, err := os.ReadFile("config/DBconfig/DBconfig.yaml")
	if err != nil {
		fmt.Printf("Error reading config.yaml file: %v\n", err)
	}

	//复制到config结构体
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		fmt.Printf("Error parsing config.yaml file: %v\n", err)
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
	config.AllowOrigins = []string{"http://localhost:8081"} // 允许的源，可以根据需要修改
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(config))

	// 创建 ImageGenerator 实例
	imgGen := generate.NewImageGenerator()

	// 注册路由
	r.POST("/register", auth_s.Register) // 注册路由
	r.POST("/login", auth_s.Login)       // 登录路由
	auth := r.Group("/auth", middlewire.JWTAuthMiddleware())
	{
		auth.POST("/generate", func(c *gin.Context) {
			imgGen.ReturnImage(c)
		})
		auth.POST("/setavator", avator.SetAvator) // 设置头像
		auth.GET("/getavator", avator.GetAvator)  // 获取头像
	}
	r.GET("/user/info", user_r.GetUserInfo)                        // 查询用户信息（根据id或username或email）
	r.GET("/user/images", image_r.GetUserImages)                   // 查询用户生成的所有图片（根据username或id）
	r.GET("/user/favoritedimages", image_r.GetUserFavoritedImages) // 查询用户收藏的图片(根据username或id)
	r.GET("/image", image_r.GetImages)                             // 查询指定的一张图片 (根据id 或图片的username属性的第一张图片)

	// 添加静态文件服务，指向 docs 目录
	r.Static("/docs", "./docs")

	// 注册 Swagger 路由，并指定 doc.json 文件的路径
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	r.Run("0.0.0.0:8080")
}

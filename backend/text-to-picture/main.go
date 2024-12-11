package main

import (
	"fmt"
	"log"
	"os"
	"text-to-picture/api/generate"
	middlewire "text-to-picture/middlewire/jwt"
	db "text-to-picture/models/init"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/yaml.v3"

	"text-to-picture/services/auth_s/avator"
	check_t "text-to-picture/services/auth_s/checkToken"
	user_d "text-to-picture/services/auth_s/delete"
	auth_s "text-to-picture/services/auth_s/login"
	user_q "text-to-picture/services/auth_s/query"
	user_up "text-to-picture/services/auth_s/update"
	favorited "text-to-picture/services/favorites_s"
	image_d "text-to-picture/services/image_s/delete"
	"text-to-picture/services/image_s/like"
	image_q "text-to-picture/services/image_s/query"
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
	if err := db.InitTestUser(); err != nil {
		log.Printf("Failed to initialize test user: %v", err)
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
	{ //Postman上测试时得在请求头上加上	Authorization：（登录时返回的Token）
		auth.POST("/generate", func(c *gin.Context) {
			imgGen.ReturnImage(c)
		})
		auth.POST("/like", like.LikeImage) // 参数?url=

		auth.POST("/setavator", avator.SetAvator) // 设置头像，参数json: url=
		auth.GET("/getavator", avator.GetAvator)  // 获取头像

		auth.POST("/addFavoritedImage", favorited.AddFavoritedImage)         // 收藏（参数：图像id或url）
		auth.DELETE("/deleteFavoritedImage", favorited.DeleteFavoritedImage) // 取消收藏（参数：?id或?url）

		// 下面3个Get请求无需参数
		auth.GET("/user/info", user_q.GetUserInfo)                        // 查询当前用户信息
		auth.GET("/user/images", image_q.GetUserImages)                   // 查询当前用户生成的所有图片
		auth.GET("/user/favoritedimages", image_q.GetUserFavoritedImages) // 查询当前用户收藏的图片

		auth.GET("/user/images/timeRange", image_q.GetImagesWithinTimeRange) // 获取当前用户指定时间段内的图像（start_time=YYYY-MM-DD&end_time=YYYY-MM-DD）
		// 或（任意一个都可）完整的时间戳格式：2006-01-02T15:04:05.000000Z

		auth.PUT("/user/update", user_up.UpdateUser)                         // 更新当前用户信息(拒绝改用户名)
		auth.DELETE("/user/deleteOne", user_d.DeleteUserByName) // 删除指定用户(?username=)

		auth.DELETE("/image/deleteOne", image_d.DeleteOneImage)     // 删除单个图像(?url=)
		auth.DELETE("/image/deleteByUser", image_d.DeleteUserImages) // 删除指定用户的所有图像(?username=)

		auth.GET("/score", user_up.AddScore) //签到增加积分接口
	}
	r.GET("/checkToken", check_t.CheckToken) //校验token是否有效

	// 以下接口暂时未需要
	r.GET("/user/all", user_q.GetAllUsersInfo)        // 获取所有用户信息

	r.GET("/image", image_q.GetImage)                        // 查询指定的一张图片 (根据id 或图片的username属性的第一张图片)
	r.GET("/image/all", image_q.GetAllImages)                // 获取所有图像信息
	r.DELETE("/image.deleteAll", image_d.DeleteAllImages)    // 删除所有图像(无参)

	// 添加静态文件服务，指向 docs 目录
	r.Static("/docs", "./docs")

	// 注册 Swagger 路由，并指定 doc.json 文件的路径
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	r.Run("0.0.0.0:8080")
}

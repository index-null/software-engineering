package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"text-to-picture/api/generate"
	middlewire "text-to-picture/middlewire/jwt"
	"text-to-picture/services/auth_s/avator"
	check_t "text-to-picture/services/auth_s/checkToken"
	user_d "text-to-picture/services/auth_s/delete"
	auth_s "text-to-picture/services/auth_s/login"
	user_q "text-to-picture/services/auth_s/query"
	user_up "text-to-picture/services/auth_s/update"
	favorited "text-to-picture/services/favorites_s"
	image_d "text-to-picture/services/image_s/delete"
	image_f "text-to-picture/services/image_s/findByFeature"
	"text-to-picture/services/image_s/like"
	image_q "text-to-picture/services/image_s/query"
)

type Service interface {
	Start()
}

type TextToPicture struct {
}

func (t *TextToPicture) Start() {
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

	// 注册登录界面
	r.POST("/register", auth_s.Register) // 注册路由
	r.POST("/login", auth_s.Login)       // 登录路由
	auth := r.Group("/auth", middlewire.JWTAuthMiddleware())
	{ //Postman上测试时得在请求头上加上	Authorization：（登录时返回的Token）
		//文生图界面
		auth.POST("/generate", func(c *gin.Context) {
			imgGen.ReturnImage(c)
		})

		//个人信息界面
		auth.POST("/setavator", avator.SetAvator)    // 设置头像，参数json: url=
		auth.GET("/getavator", avator.GetAvator)     // 获取头像
		auth.GET("/user/info", user_q.GetUserInfo)   // 查询当前用户信息
		auth.PUT("/user/update", user_up.UpdateUser) // 更新当前用户信息(拒绝改用户名)

		//历史记录界面
		auth.GET("/user/images", image_q.GetUserImages)                      // 查询当前用户生成的所有图片
		auth.GET("/user/images/timeRange", image_q.GetImagesWithinTimeRange) // 获取当前用户指定时间段内的图像（start_time=YYYY-MM-DD&end_time=YYYY-MM-DD）
		// 或（任意一个都可）完整的时间戳格式：2006-01-02T15:04:05.000000Z
		auth.POST("/like", like.LikeImage) // 参数?url=

		//收藏图像界面
		auth.GET("/user/favoritedimages", image_q.GetUserFavoritedImages)    // 查询当前用户收藏的图片
		auth.POST("/addFavoritedImage", favorited.AddFavoritedImage)         // 收藏（参数：图像id或url）
		auth.DELETE("/deleteFavoritedImage", favorited.DeleteFavoritedImage) // 取消收藏（参数：?id或?url）

		//积分功能
		auth.GET("/score", user_up.AddScore) //签到增加积分接口

		//图像界面

		auth.POST("/user/deleteImages", image_d.DeleteUserImagesBatch) // (批量)删除当前用户的图像

		//搜索界面
		auth.GET("/image/feature", image_f.FindByFeature) //查询所有图像中或当前用户的图像中 图像的Prompt中包含所输入的（可多个）关键字的所有图像

		// root操作
		auth.DELETE("/root/deleteOneUser", user_d.DeleteUserByName)       // 删除指定用户(?username=)
		auth.DELETE("/root/deleteOneImage", image_d.DeleteOneImage)       // 删除单个图像(?url=)
		auth.DELETE("/root/deleteImagesByUser", image_d.DeleteUserImages) // 删除指定用户的所有图像(?username=)
	}
	r.GET("/checkToken", check_t.CheckToken) //校验token是否有效

	// 以下接口暂时未需要
	r.GET("/user/all", user_q.GetAllUsersInfo) // 获取所有用户信息

	r.GET("/image", image_q.GetImage)                          // 查询指定的一张图片 (根据id 或图片的username属性的第一张图片)
	r.GET("/image/all", image_q.GetAllImages)                  // 获取所有图像信息
	r.DELETE("/root/deleteAllImages", image_d.DeleteAllImages) // 删除所有图像(无参)

	// 添加静态文件服务，指向 docs 目录
	r.Static("/docs", "./docs")

	// 注册 Swagger 路由，并指定 doc.json 文件的路径
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	r.Run("0.0.0.0:8080")
}

package findByFeature

import (
	"log"
	"net/http"

	d "text-to-picture/models/init"
	"text-to-picture/models/repository/image_r"

	//"text-to-picture/models/repository/user_r"

	//u "text-to-picture/models/user"

	"github.com/gin-gonic/gin"
)

func FindByFeature(c *gin.Context) {
	// 从查询参数中获取特征列表
	features := c.QueryArray("feature")
	isOwn := c.Query("isOwn")

	var username string
	username = ""
	if isOwn == "true" || isOwn == "True" || isOwn == "TRUE"{
		userName, exists := c.Get("username")
		if !exists {
			log.Printf("未找到用户名")
			c.JSON(401, gin.H{
				"success": false,
				"message": "未找到用户信息",
			})
			return
		}
		username = userName.(string)
	}

	// 调用业务逻辑层函数查找图片
	images, err := image_r.FindByFeature(d.DB, username, features)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "根据关键字查询图片失败",
			"error":err.Error(),
		})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"images":    images,
	})
}

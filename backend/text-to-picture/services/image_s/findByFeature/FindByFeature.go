package findbyfeature

import (
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

	// 调用业务逻辑层函数查找图片
	images, err := image_r.FindByFeature(d.DB, features)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"success": false,
			"message": "查询图片失败",
		})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"data":    images,
	})
}

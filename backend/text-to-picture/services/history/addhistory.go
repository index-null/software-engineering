package history

//接收前端传来的Prompt，width,height,seed,steps,picture(url)以及从上下文userName, exists := c.Get("username")获取的username，插入到图片信息表ImageInformation中
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	db "text-to-picture/models/init"
	"time"

	"text-to-picture/models/image"
)

type History struct {
	Prompt     string `json:"prompt"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Seed       int    `json:"seed"`
	Steps      int    `json:"steps"`
	PictureURL string `json:"pictureURL"`
}

// AddHistory 插入图片信息到数据库
func AddHistory(c *gin.Context) {
	var h History
	err := c.BindJSON(&h)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"success": false,
			"message": "参数解析失败",
		})
		return
	}
	// 从上下文中获取用户名
	userName, exists := c.Get("username")
	if !exists {
		log.Printf("未找到用户名")
		c.JSON(401, gin.H{
			"code":    401,
			"success": false,
			"message": "未找到用户信息",
		})
		return
	}

	// 获取数据库连接
	db := db.DB
	if db == nil {
		c.JSON(500, gin.H{
			"code":    500,
			"success": false,
			"message": "插入获取数据库连接失败",
		})
		return
	}

	// 创建 ImageInformation 结构体
	imageInfo := image.ImageInformation{
		UserName: userName.(string),
		Params: fmt.Sprintf("\"Prompt\": \"%s\", \"Width\": \"%d\", \"Height\": \"%d\", \"Steps\": \"%d\",\"Seed\": \"%d\"",
			h.Prompt, h.Width, h.Height, h.Steps, h.Seed),
		Picture:     h.PictureURL,
		Create_time: time.Now(),
	}

	// 插入数据到 ImageInformation 表
	result := db.Create(&imageInfo)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"success": false,
			"message": "插入数据到数据库失败",
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"success": true,
		"message": "插入数据到历史记录成功",
	})
}

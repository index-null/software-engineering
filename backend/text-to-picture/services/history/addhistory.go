package history

//接收前端传来的Prompt，width,height,seed,steps,picture(url)以及从上下文userName, exists := c.Get("username")获取的username，插入到图片信息表ImageInformation中
import (
	"fmt"
	"log"
	db "text-to-picture/models/init"
	"time"

	"github.com/gin-gonic/gin"

	"text-to-picture/models/image"
)

type History struct {
	Prompt     string `json:"prompt"`     // 前端传来的生成图片的文本提示
	Width      int    `json:"width"`      // 图片宽度
	Height     int    `json:"height"`     // 图片高度
	Seed       int    `json:"seed"`       // 随机种子
	Steps      int    `json:"steps"`      // 生成图片的步数
	PictureURL string `json:"pictureURL"` // 图片的URL
}

// AddHistory 插入图片信息到数据库
func AddHistory(c *gin.Context) {
	var h History
	// 解析前端传来的 JSON 数据到 History 结构体
	err := c.BindJSON(&h)
	if err != nil {
		// 如果解析失败，返回错误信息
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
		// 如果未找到用户名，返回未授权错误信息
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
		// 如果数据库连接失败，返回服务器错误信息
		c.JSON(500, gin.H{
			"code":    500,
			"success": false,
			"message": "插入获取数据库连接失败",
		})
		return
	}

	// 创建 ImageInformation 结构体用于保存图片信息
	imageInfo := image.ImageInformation{
		UserName: userName.(string), // 从上下文中提取的用户名
		Params: fmt.Sprintf("\"Prompt\": \"%s\", \"Width\": \"%d\", \"Height\": \"%d\", \"Steps\": \"%d\",\"Seed\": \"%d\"",
			h.Prompt, h.Width, h.Height, h.Steps, h.Seed), // 将图片生成参数序列化为字符串
		Picture:     h.PictureURL, // 图片 URL
		Create_time: time.Now(),   // 当前时间作为创建时间
	}

	// 插入数据到 ImageInformation 表
	result := db.Create(&imageInfo)
	if result.Error != nil {
		// 如果插入数据库失败，返回服务器错误信息
		c.JSON(500, gin.H{
			"code":    500,
			"success": false,
			"message": "插入数据到数据库失败",
		})
		return
	}

	// 如果操作成功，返回成功信息
	c.JSON(200, gin.H{
		"code":    200,
		"success": true,
		"message": "插入数据到历史记录成功",
	})
}

package generate_s

import (
	"fmt"
	"reflect"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 传入的图片参数
type ImageParaments struct {
	Prompt         string `json:"prompt" binding:"required" fault:"缺乏提示词"`
	Width          int    `json:"width" binding:"required,min=128,max=1024" fault:"宽度不在范围内"`
	Height         int    `json:"height" binding:"required,min=128,max=1024" fault:"高度不在范围内"`
	Steps          int    `json:"steps" binding:"required,min=1,max=100" fault:"步数不在范围内"`
	SamplingMethod string `json:"sampling_method" binding:"required,oneof=DDIM PLMS K-LMS" fault:"采样方法不在范围内"`
	Seed           string `json:"seed" binding:"required," fault:"缺乏种子"`
}

// 获取在Tag中的fault信息
func ParamentsError(err error, obj any) string {
	getObj := reflect.TypeOf(obj)
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, e := range err {
			if f, exists := getObj.Elem().FieldByName(e.Field()); exists {
				fault := f.Tag.Get("fault")
				return fault
			}
		}
	}
	return err.Error()
}

//判断OSS返回的url是否正确
func IsValidURL(inputURL string) bool {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}
// 接收传来的图片参数,并进行数据校验
func AcceptParaments(c *gin.Context) error {
	var imageParaments ImageParaments
	if err := c.ShouldBindJSON(&imageParaments); err != nil {
		fault := ParamentsError(err, &imageParaments)
		return fmt.Errorf(fault)
	}

	return nil
}

// 返回图片生成图片的url
func ReturnImage(c *gin.Context) {
	//校验参数
	if err := AcceptParaments(c); err != nil {
		c.JSON(400, gin.H{
			"success":false,
			"message": err,
		})
		return
	}

	imageUrl , err := GenerateImage()

	//校验生成图片
	 if err != nil {
		c.JSON(500, gin.H{
			"success":false,
			"message": "图片生成失败",
		})
		return
	}

	//校验图片url
	if !IsValidURL(imageUrl) {
		c.JSON(500, gin.H{
			"success":false,
			"message": "无效url",
		})
		return
	}

	c.JSON(200,gin.H{
		"success":true,
		"image_url":imageUrl,
	})
}

func GenerateImage() (string,error){
	//这里把图片上传到OSS,OSS会那里返回包含图片URL的json
	return "https://your-bucket-name.oss-region.aliyuncs.com/user/avatar.png",nil
}


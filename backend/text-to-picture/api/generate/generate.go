package generate

import (
	"github.com/gin-gonic/gin"
	"text-to-picture/services/generate_s"
)

type ImageGenerator interface {
	ReturnImage(c *gin.Context)
}

// NewImageGenerator 创建一个新的 ImageGenerator 实例
func NewImageGenerator() ImageGenerator {
	return &generate_s.ImageGeneratorImpl{}
}

package image_r

import (
	"fmt"
	"strings"

	i "text-to-picture/models/image"

	"gorm.io/gorm"
)

// IsImageFavoritedByUser checks if a user has favorited a specific image
func FindByFeature(db *gorm.DB, feature []string) ([]i.ImageInformation, error) {
	if len(feature) == 0 {
		return nil, nil // 如果没有提供特征，则返回空列表和无错误
	}

	// 构建SQL语句中的 OR 条件
	var conditions []string
	args := make([]interface{}, 0, len(feature))

	for _, f := range feature {
		// 对每个特征进行转义处理，防止SQL注入
		escapedFeature := strings.ReplaceAll(f, "'", "''")
		// 匹配规则：寻找从 "Prompt": " 开始到 ", "Width": 之间的内容，即Prompt的内容
		conditions = append(conditions, "params LIKE ?")
		args = append(args, fmt.Sprintf("%%\"Prompt\": \"%%%s%%\"%%", escapedFeature))
	}

	fmt.Println("查询条件：Conditions:", conditions, "\nArgs:", args)

	// 创建一个切片来存储结果
	var images []i.ImageInformation

	// 使用 GORM 构造查询并执行
	query := db.Table("imageinformation").Where(strings.Join(conditions, " OR "), args...)

	err := query.Find(&images).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// // 进一步过滤结果，确保只有 Prompt 符合条件的数据被保留
	// filteredImages := make([]i.ImageInformation, 0)
	// for _, img := range images {
	// 	for _, f := range feature {
	// 		if strings.Contains(img.Params, fmt.Sprintf("\"Prompt\": %s,", f)) {
	// 			filteredImages = append(filteredImages, img)
	// 			break // 只要找到一个符合条件的feature就跳出循环
	// 		}
	// 	}
	// }

	// return filteredImages, nil

	return images, nil
}

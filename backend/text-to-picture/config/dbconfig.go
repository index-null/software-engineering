// config/dbconfig.go
package config

import (
	"log"
	"path/filepath"
	"runtime"
)

func GetDBConfigPath() string {
	// 获取调用者的文件名（即 login_test.go 或 findByFeature.go）
	_, filename, _, ok := runtime.Caller(2) // 注意这里使用 Caller(2)
	if !ok {
		log.Fatal("无法获取运行时调用者信息")
	}

	// 获取当前文件所在的目录
	currentDir := filepath.Dir(filename)

	// 构建到项目根目录的相对路径
	dbConfigPath := filepath.Join(currentDir, "..", "..", "..", "config", "DBconfig", "DBconfig.yaml")

	// 将路径转换为绝对路径并简化路径（移除冗余的 '..'）
	absPath, err := filepath.Abs(dbConfigPath)
	if err != nil {
		log.Fatalf("无法获取绝对路径: %v", err)
	}

	simplifiedPath := filepath.Clean(absPath)

	return simplifiedPath
}

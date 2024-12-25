// config/dbconfig.go
package config

import (
	"log"
	"path/filepath"
	"runtime"
)

// getDBConfigPath 返回数据库配置文件的绝对路径。
// 该函数通过获取调用者文件的位置来确定当前项目的根目录，然后构建到数据库配置文件的路径。
// 这种方法允许在不同的测试或执行环境中灵活地定位配置文件。
func getDBConfigPath() string {
    // 获取调用者的文件名（即 login_test.go 或 findByFeature.go）
    // 注意这里使用 Caller(2)，因为直接调用者是 GetDBConfigPath 函数，我们需要的是 GetDBConfigPath 的调用者
    _, filename, _, ok := runtime.Caller(2)
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

// GetDBConfigPath 提供了一个公共接口来获取数据库配置文件的绝对路径。
// 该函数实际上调用了 getDBConfigPath 来执行实际的路径解析。
// 使用这个公共接口可以确保在不同的上下文中以一致的方式获取配置文件的路径。
func GetDBConfigPath() string {
    return getDBConfigPath()
}

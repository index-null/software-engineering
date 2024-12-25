// config/config.go
package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type OSSConfig struct {
	OSS_REGION            string `yaml:"OSS_REGION"`
	OSS_ACCESS_KEY_ID     string `yaml:"OSS_ACCESS_KEY_ID"`
	OSS_ACCESS_KEY_SECRET string `yaml:"OSS_ACCESS_KEY_SECRET"`
	OSS_BUCKET            string `yaml:"OSS_BUCKET"`
}

type ModelConfig struct {
	GEN_API_KEY string `yaml:"GEN_API_KEY"`
	Time        string `yaml:"timeout"`
}

type Config struct {
	DB    DBConfig    `yaml:"db"`
	OSS   OSSConfig   `yaml:"oss"`
	Model ModelConfig `yaml:"model"`
}

func getDBConfigPath() string {
	// 获取调用者的文件名（即 login_test.go 或 findByFeature.go）
	_, filename, _, ok := runtime.Caller(2) // 注意这里使用 Caller(2)
	if !ok {
		log.Fatal("无法获取运行时调用者信息")
	}

	// 获取当前文件所在的目录
	currentDir := filepath.Dir(filename)

	// 构建到项目根目录的相对路径
	dbConfigPath := filepath.Join(currentDir, "..", "config", "configs", "config.yaml")

	// 将路径转换为绝对路径并简化路径（移除冗余的 '..'）
	absPath, err := filepath.Abs(dbConfigPath)
	if err != nil {
		log.Fatalf("无法获取绝对路径: %v", err)
	}

	simplifiedPath := filepath.Clean(absPath)

	return simplifiedPath
}

func GetDBConfigPath() string {
	return getDBConfigPath()
}
func LoadConfig() (*Config, error) {
	configPath := GetDBConfigPath()
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("读取配置文件失败: %v", err)
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Printf("解析配置文件失败: %v", err)
		return nil, err
	}

	return &config, nil
}

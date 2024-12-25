package main

import (
	"log"
	"os"
	"text-to-picture/config"
	db "text-to-picture/models/init"
	"text-to-picture/services/router"
)

type DBConfig struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

func main() {

	//读取DBConfig.yaml文件
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	//设置数据库连接的环境变量
	os.Setenv("DB_USER", config.DB.User)
	os.Setenv("DB_PASSWORD", config.DB.Password)
	os.Setenv("DB_HOST", config.DB.Host)
	os.Setenv("DB_PORT", config.DB.Port)
	os.Setenv("DB_NAME", config.DB.Name)
	//设置模型连接的环境变量
	os.Setenv("GEN_API_KEY", config.Model.GEN_API_KEY)
	os.Setenv("TIMEOUT", config.Model.Time)
	//设置OSS连接的环境变量
	os.Setenv("OSS_REGION", config.OSS.OSS_REGION)
	os.Setenv("OSS_BUCKET", config.OSS.OSS_BUCKET)
	os.Setenv("OSS_ACCESS_KEY_ID", config.OSS.OSS_ACCESS_KEY_ID)
	os.Setenv("OSS_ACCESS_KEY_SECRET", config.OSS.OSS_ACCESS_KEY_SECRET)

	// 连接数据库
	if err := db.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	if err := db.InitTestUser(); err != nil {
		log.Printf("Failed to initialize test user: %v", err)
	}

	var serve router.TextToPicture
	serve.Start()
}

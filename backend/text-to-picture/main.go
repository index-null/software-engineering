package main

import (
	"log"
	"os"

	"text-to-picture/api/auth"       // 导入注册和登录路由
	db "text-to-picture/models/init" // 为 init 包设置别名为 db

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
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
	yamlFile, err := os.ReadFile("config\\DBconfig\\DBconfig.yaml")
	if err != nil {
		log.Fatalf("Error reading config.yaml file: %v", err)
	}

	//复制到config结构体
	var config DBConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error parsing config.yaml file: %v", err)
	}

	//设置数据库连接的环境变量
	os.Setenv("DB_USER", config.DB.User)
	os.Setenv("DB_PASSWORD", config.DB.Password)
	os.Setenv("DB_HOST", config.DB.Host)
	os.Setenv("DB_PORT", config.DB.Port)
	os.Setenv("DB_NAME", config.DB.Name)

	// 连接数据库
	if err := db.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置路由
	r := gin.Default()
	r.POST("/register", auth.Register) // 注册路由
	r.POST("/login", auth.Login)       // 登录路由

	r.Run() // 启动服务器
}

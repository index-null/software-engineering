package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	db "text-to-picture/models/init"
	"text-to-picture/services/router"

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
	yamlFile, err := os.ReadFile("config/DBconfig/DBconfig.yaml")
	if err != nil {
		fmt.Printf("Error reading config.yaml file: %v\n", err)
	}

	//复制到config结构体
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		fmt.Printf("Error parsing config.yaml file: %v\n", err)
	}

	envPath := filepath.Join("config", "oss", "oss.env")
	//envPath := "D:/软件工程项目/software-engineering/backend/text-to-picture/config/oss/oss.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}
	//设置数据库连接的环境变量
	os.Setenv("DB_USER", dbconfig.DB.User)
	os.Setenv("DB_PASSWORD", dbconfig.DB.Password)
	os.Setenv("DB_HOST", dbconfig.DB.Host)
	os.Setenv("DB_PORT", dbconfig.DB.Port)
	os.Setenv("DB_NAME", dbconfig.DB.Name)

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

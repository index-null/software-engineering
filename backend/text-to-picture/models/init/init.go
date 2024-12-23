package models

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	image2 "text-to-picture/models/image"
	user2 "text-to-picture/models/user"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS UserInformation (
    id SERIAL PRIMARY KEY,
	email VARCHAR(50) UNIQUE NOT NULL,
	username VARCHAR(30) UNIQUE NOT NULL,
    password VARCHAR(256) NOT NULL,
	avatar_url VARCHAR(255) NOT NULL,
    score INT DEFAULT 0,
	create_time TIMESTAMP DEFAULT NOW(),
    token VARCHAR(255)
);
CREATE TABLE IF NOT EXISTS UserScore (
    id SERIAL PRIMARY KEY,
	username VARCHAR(30) NOT NULL,
    record TEXT,
	create_time TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (username) REFERENCES UserInformation(username) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS ImageInformation (
    id SERIAL PRIMARY KEY,
	userName VARCHAR(30) NOT NULL,
    params TEXT,
    picture TEXT UNIQUE,
    likecount INT DEFAULT 0,
    create_time TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (userName) REFERENCES UserInformation(username) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS ImageLike (
    id SERIAL PRIMARY KEY,
    picture TEXT,
    username TEXT,
    num INT DEFAULT 0,
    create_time TIMESTAMP DEFAULT NOW(),
	FOREIGN KEY (username) REFERENCES UserInformation(username) ON DELETE CASCADE,
    FOREIGN KEY (picture) REFERENCES ImageInformation(picture) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS FavoritedImage (
	id SERIAL PRIMARY KEY,
	userName VARCHAR(30) NOT NULL,
	picture TEXT,
	create_time TIMESTAMP DEFAULT NOW(),
	FOREIGN KEY (userName) REFERENCES UserInformation(username) ON DELETE CASCADE,
    FOREIGN KEY (picture) REFERENCES ImageInformation(picture) ON DELETE CASCADE
);


`

// UserImformation中avatar_url为头像图片url

var DB *gorm.DB

func ConnectDatabase() error {
	var err error

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}
	return nil
}
func InitDB() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	tx := DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Exec(createTableSQL).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func InitTestUser() error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user user2.UserInformation
	currentTime := time.Now()
	result := tx.Where("username=?", "root").First(&user)

	if result.Error == nil {
		log.Printf("User already exists")

		// 更新用户分数
		user.Score = 10000
		if result := tx.Save(&user); result.Error != nil {
			log.Printf("Failed to update user score: %v", result.Error)
			tx.Rollback()
			return result.Error
		}

		tx.Commit()
		return nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Failed to find user: %v", result.Error)
		tx.Rollback()
		return result.Error
	}

	// 创建用户信息
	user = user2.UserInformation{
		Email:       "root@example.com",
		UserName:    "root",
		Password:    "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a", //111111
		Avatar_url:  "https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202407272335307.png",
		Score:       10000,
		Create_time: currentTime.AddDate(-1, 0, 0),
	}
	if result := tx.Create(&user); result.Error != nil {
		log.Printf("Failed to create user: %v", result.Error)
		tx.Rollback()
		return result.Error
	}

	// 创建用户积分记录
	userscore := user2.UserScore{
		Username:    "root",
		Record:      "积分+100",
		Create_time: currentTime.AddDate(0, -1, 0),
	}
	if result := tx.Create(&userscore); result.Error != nil {
		log.Printf("Failed to create record: %v", result.Error)
		tx.Rollback()
		return result.Error
	}

	// 处理图片文件路径
	ImgfilePath := filepath.Join("assets", "examples", "images", "image_urls.txt")
	file, err := os.Open(ImgfilePath)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		tx.Rollback()
		return err
	}
	defer file.Close()

	var imageUrls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		imageUrls = append(imageUrls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read file : %v", err)
		tx.Rollback()
		return fmt.Errorf("&v", err)
	}
	i := 0

	// 创建图像记录
	for _, url := range imageUrls {
		i++
		test := fmt.Sprintf("test%d", i)
		createTime := currentTime.AddDate(0, 0, -(11 + i))
		imageInfo := image2.ImageInformation{
			UserName:    "root", // 假设用户名为root
			Params:      "\"Prompt\": \"" + test + "\", \"Width\": \"512\", \"Height\": \"512\", \"Steps\": \"20\", \"SamplingMethod\": \"Euler a\"",
			Picture:     url,
			Create_time: createTime,
		}

		result := tx.Create(&imageInfo)
		if result.Error != nil {
			log.Printf("Faile to create image information: %v", result.Error)
			return result.Error
		}
	}

	// 如果超过10 张图像，随机选取10个用于点赞和收藏
	if len(imageUrls) > 10 {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(imageUrls), func(i, j int) { imageUrls[i], imageUrls[j] = imageUrls[j], imageUrls[i] })
		imageUrls = imageUrls[:10]
	}

	currentTime = time.Now()

	// 创建点赞和收藏记录
	for i, url := range imageUrls {
		createTime := currentTime.AddDate(0, 0, -(i + 1))
		// 创建图像收藏记录
		imagefavor := image2.FavoritedImages{
			UserName:    "root",
			Picture:     url,
			Create_time: createTime,
		}
		if result := tx.Create(&imagefavor); result.Error != nil {
			log.Printf("Failed to create image favor for URL %s: %v", url, result.Error)
			tx.Rollback()
			return result.Error
		}

	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

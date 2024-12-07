package models

import (
	"fmt"
	"os"

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
    FOREIGN KEY (id) REFERENCES UserInformation(id)
);
CREATE TABLE IF NOT EXISTS ImageInformation (
    id SERIAL PRIMARY KEY,
	userName VARCHAR(30) NOT NULL,
    params TEXT,
    picture TEXT UNIQUE,
    likecount INT DEFAULT 0,
    create_time TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (userName) REFERENCES UserInformation(username)
);
CREATE TABLE IF NOT EXISTS ImageLike (
    id SERIAL PRIMARY KEY,
    picture TEXT,
    username TEXT,
    num INT DEFAULT 0,
    create_time TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (picture) REFERENCES ImageInformation(picture)
);

CREATE TABLE IF NOT EXISTS FavoritedImage (
	id SERIAL PRIMARY KEY,
	userName VARCHAR(30) NOT NULL,
	picture TEXT,
	create_time TIMESTAMP DEFAULT NOW(),
	FOREIGN KEY (userName) REFERENCES UserInformation(username)
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

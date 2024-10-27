package models

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS user_login (
    id INT PRIMARY KEY,
	email VARCHAR(255) UNIQUE NOT NULL,
	username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS user_query (
    id INT PRIMARY KEY,
	user_id INT NOT NULL,
    params TEXT,
    result TEXT,
    time TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES user_login(id)
);

CREATE TABLE IF NOT EXISTS image (
	id INT PRIMARY KEY,
	user_id INT NOT NULL,
	result TEXT,
	FOREIGN KEY (user_id) REFERENCES user_login(id)
)
`

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

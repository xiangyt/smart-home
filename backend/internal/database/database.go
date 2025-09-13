package database

import (
	"fmt"
	"os"
	"smart-home-backend/internal/models"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// 根据数据库URL选择驱动
	if strings.Contains(databaseURL, "@tcp(") {
		// MySQL数据库 - 先创建数据库
		dbName := os.Getenv("DATABASE_NAME")
		if dbName == "" {
			dbName = "smart_home"
		}

		// 先连接到mysql系统数据库创建目标数据库
		systemDB, err := gorm.Open(mysql.Open(databaseURL), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
		}

		// 创建数据库
		createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
		if err := systemDB.Exec(createSQL).Error; err != nil {
			return nil, fmt.Errorf("failed to create database: %v", err)
		}

		// 关闭系统数据库连接
		sqlDB, _ := systemDB.DB()
		sqlDB.Close()

		// 重新构建连接字符串，连接到目标数据库
		targetURL := strings.Replace(databaseURL, "/mysql?", "/"+dbName+"?", 1)
		db, err = gorm.Open(mysql.Open(targetURL), &gorm.Config{})
	} else {
		// SQLite数据库
		db, err = gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&models.Device{},
		&models.DeviceHistory{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

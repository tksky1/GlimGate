package database

import (
	"fmt"
	"log"

	"github.com/tksky1/glimgate/internal/model"
	"github.com/tksky1/glimgate/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	dsn := config.AppConfig.Database.GetDSN()
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 自动迁移数据表
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("数据表迁移失败: %w", err)
	}

	log.Println("数据库连接成功")
	return nil
}

// autoMigrate 自动迁移数据表
func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Direction{},
		&model.Problem{},
		&model.SubmissionPoint{},
		&model.Submission{},
		&model.Score{},
	)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
package main

import (
	"log"

	"github.com/tksky1/glimgate/internal/model"
	"github.com/tksky1/glimgate/pkg/config"
	"github.com/tksky1/glimgate/pkg/database"
	"github.com/tksky1/glimgate/pkg/utils"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	db := database.GetDB()

	// 检查是否已存在管理员账户
	var adminCount int64
	if err := db.Model(&model.User{}).Where("is_admin = ?", true).Count(&adminCount).Error; err != nil {
		log.Fatalf("检查管理员账户失败: %v", err)
	}

	if adminCount > 0 {
		log.Println("管理员账户已存在，跳过初始化")
		return
	}

	// 创建默认管理员账户
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	admin := model.User{
		Username:  "admin",
		Password:  hashedPassword,
		Nickname:  "系统管理员",
		RealName:  "管理员",
		College:   "系统",
		StudentID: "ADMIN001",
		Email:     "admin@glimgate.com",
		IsAdmin:   true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("创建管理员账户失败: %v", err)
	}

	log.Println("默认管理员账户创建成功:")
	log.Println("用户名: admin")
	log.Println("密码: admin123")
	log.Println("请在生产环境中及时修改默认密码！")

	// 创建示例方向
	directions := []model.Direction{
		{
			Name:        "前端开发",
			Description: "负责前端页面开发和用户交互设计",
		},
		{
			Name:        "后端开发",
			Description: "负责服务器端逻辑和API开发",
		},
		{
			Name:        "移动开发",
			Description: "负责iOS和Android移动应用开发",
		},
		{
			Name:        "UI/UX设计",
			Description: "负责用户界面和用户体验设计",
		},
	}

	for _, direction := range directions {
		if err := db.Create(&direction).Error; err != nil {
			log.Printf("创建方向 %s 失败: %v", direction.Name, err)
		} else {
			log.Printf("创建方向 %s 成功", direction.Name)
		}
	}

	log.Println("数据库初始化完成！")
}
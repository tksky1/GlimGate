package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tksky1/glimgate/docs" // 添加这行
	"github.com/tksky1/glimgate/internal/middleware"
	"github.com/tksky1/glimgate/internal/router"
	"github.com/tksky1/glimgate/pkg/config"
	"github.com/tksky1/glimgate/pkg/database"
)

// @title GlimGate API
// @version 1.0
// @description GlimGate工作室招新交题和评分系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer token

func main() {
	// 加载配置
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建Gin引擎
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORSMiddleware())
	r.SetTrustedProxies(nil)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置路由
	router.SetupRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf("0.0.0.0:%d", config.AppConfig.Server.Port)
	log.Printf("服务器启动在端口 %d", config.AppConfig.Server.Port)
	log.Printf("API文档地址: http://localhost%s/swagger/index.html", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

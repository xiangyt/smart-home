package main

import (
	"log"
	"os"

	"smart-home-backend/internal/config"
	"smart-home-backend/internal/database"
	"smart-home-backend/internal/handlers"
	"smart-home-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化服务
	haService := services.NewHomeAssistantService(cfg.HomeAssistantURL, cfg.HomeAssistantToken)
	deviceService := services.NewDeviceService(db, haService)

	// 初始化处理器
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002", "http://127.0.0.1:3000", "http://127.0.0.1:3001", "http://127.0.0.1:3002"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API路由
	api := r.Group("/api")
	{
		// 设备相关路由
		devices := api.Group("/devices")
		{
			devices.GET("", deviceHandler.GetDevices)
			devices.GET("/:entity_id/state", deviceHandler.GetDeviceState)
			devices.POST("/:entity_id/control", deviceHandler.ControlDevice)
			devices.GET("/:entity_id/history", deviceHandler.GetDeviceHistory)
		}

		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

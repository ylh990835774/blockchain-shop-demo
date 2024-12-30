package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/configs"
	"github.com/ylh990835774/blockchain-shop-demo/internal/api"
	"github.com/ylh990835774/blockchain-shop-demo/internal/api/middleware"
	"github.com/ylh990835774/blockchain-shop-demo/internal/handlers"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"github.com/ylh990835774/blockchain-shop-demo/internal/service"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := configs.Load()
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 初始化日志
	logConfig := &logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
		Console:    cfg.Log.Console,
	}
	if err := logger.Setup(logConfig); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Sync()

	// 初始化数据库连接
	db, err := mysql.NewDB(&mysql.Config{
		Host:     cfg.MySQL.Host,
		Port:     cfg.MySQL.Port,
		Username: cfg.MySQL.Username,
		Password: cfg.MySQL.Password,
		Database: cfg.MySQL.Database,
	})
	if err != nil {
		logger.Fatal("初始化数据库失败", logger.Err(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("获取数据库连接失败", logger.Err(err))
	}
	defer sqlDB.Close()

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化 Gin 引擎
	router := gin.New()

	// 添加全局中间件
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())    // 错误处理中间件
	router.Use(middleware.ResponseHandler()) // 响应处理中间件

	// 初始化存储层
	userRepo := mysql.NewUserRepository(db)
	productRepo := mysql.NewProductRepository(db)
	orderRepo := mysql.NewOrderRepository(db)
	blockchainRepo := mysql.NewBlockchainRepository(db)

	// 初始化服务层
	jwtService := service.NewJWTService(cfg.JWT.SecretKey,
		cfg.JWT.Issuer, time.Hour*time.Duration(cfg.JWT.ExpireDurationHours))
	userService := service.NewUserService(userRepo, jwtService)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, blockchainRepo, db)

	// 初始化处理器
	h := handlers.NewHandlers(userService, jwtService, productService, orderService)

	// 设置路由
	api.SetupRouter(router, h, middleware.NewJWTMiddleware(jwtService))

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("启动服务器失败", logger.Err(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务器关闭失败", logger.Err(err))
	}

	logger.Info("服务器已关闭")
}

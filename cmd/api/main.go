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
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := configs.Load()
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 初始化日志
	logCfg := &logger.Config{
		Level:      cfg.Log.Level,
		Encoding:   cfg.Log.Encoding,
		OutputPath: cfg.Log.OutputPath,
	}
	log, err := logger.NewLogger(logCfg)
	if err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer log.Sync()

	// 初始化数据库连接
	db, err := mysql.NewDB(&mysql.Config{
		Host:     cfg.MySQL.Host,
		Port:     cfg.MySQL.Port,
		Username: cfg.MySQL.Username,
		Password: cfg.MySQL.Password,
		Database: cfg.MySQL.Database,
	})
	if err != nil {
		log.Fatal("初始化数据库失败", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取数据库连接失败", zap.Error(err))
	}
	defer sqlDB.Close()

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化 Gin 引擎
	router := gin.New()

	// 添加全局中间件
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(log))
	router.Use(middleware.ErrorHandler(log)) // 错误处理中间件
	router.Use(middleware.ResponseHandler()) // 响应处理中间件

	// 初始化存储层
	userRepo := mysql.NewUserRepository(db)
	productRepo := mysql.NewProductRepository(db)
	orderRepo := mysql.NewOrderRepository(db)
	blockchainRepo := mysql.NewBlockchainRepository(db)

	// 初始化服务层
	jwtService := service.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.Issuer, time.Hour*time.Duration(cfg.JWT.ExpireDurationHours))
	userService := service.NewUserService(userRepo, jwtService)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, blockchainRepo, db)

	// 初始化处理器
	h := handlers.NewHandlers(userService, jwtService, productService, orderService)

	// 设置路由 - 修复这里，传入正确的router实例
	api.SetupRouter(router, h, log, middleware.NewJWTMiddleware(jwtService))

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("启动服务器失败", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器关闭失败", zap.Error(err))
	}

	log.Info("服务器已关闭")
}

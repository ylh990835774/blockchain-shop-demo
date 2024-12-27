package main

import (
	"fmt"
	"log"

	"github.com/ylh990835774/blockchain-shop-demo/configs"
	"github.com/ylh990835774/blockchain-shop-demo/internal/api"
	"github.com/ylh990835774/blockchain-shop-demo/internal/api/middleware"
	"github.com/ylh990835774/blockchain-shop-demo/internal/blockchain"
	"github.com/ylh990835774/blockchain-shop-demo/internal/handlers"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"github.com/ylh990835774/blockchain-shop-demo/internal/service"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/logger"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("无法加载配置:", err)
	}

	// 初始化日志
	loggerCfg := &logger.Config{
		Level:      "info",
		Filename:   "./storage/logs/github.com/ylh990835774/blockchain-shop-demo.log",
		MaxSize:    500,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}
	logger, err := logger.NewLogger(loggerCfg)
	if err != nil {
		log.Fatal("初始化日志失败:", err)
	}
	defer logger.Sync()

	// 初始化 JWT
	utils.InitJWT(config.JWT.SecretKey)

	// 初始化数据库
	db, err := mysql.NewDB(&mysql.Config{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Username: config.Database.Username,
		Password: config.Database.Password,
		Database: config.Database.Database,
	})
	if err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}

	// 初始化仓储层
	userRepo := mysql.NewUserRepository(db)
	productRepo := mysql.NewProductRepository(db)
	orderRepo := mysql.NewOrderRepository(db)

	// 初始化区块链服务
	blockchainSvc, err := blockchain.NewBlockchainService()
	if err != nil {
		logger.Fatal("创建区块链服务失败", zap.Error(err))
	}

	// 初始化服务
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productService, blockchainSvc)
	jwtService := service.NewJWTService(config.JWT.SecretKey, config.JWT.Issuer)

	// 初始化 API 处理器
	apiHandlers := handlers.NewHandlers(userService, productService, orderService, blockchainSvc, jwtService)

	// 初始化 Gin 引擎
	r := gin.New()

	// 注册中间件
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.ErrorHandler(logger))

	// 注册 API 路由
	api.RegisterRoutes(r, apiHandlers)

	// 注册 JWT 中间件
	authMiddleware, err := middleware.NewJWTMiddleware(config.JWT.SecretKey, config.JWT.Issuer)
	if err != nil {
		logger.Fatal("初始化 JWT 中间件失败", zap.Error(err))
	}
	api.RegisterAuthRoutes(r, apiHandlers, authMiddleware)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", config.Server.Port)
	logger.Info("服务器启动", zap.String("address", serverAddr))

	if err := r.Run(serverAddr); err != nil {
		logger.Fatal("服务器启动失败", zap.Error(err))
	}
}

# Blockchain Shop

一个基于区块链技术的电商系统演示项目，实现了订单交易信息上链存证的功能。

## 功能特性

- 🔐 用户认证
  - 用户注册
  - 用户登录 (JWT 认证)
- 🛍️ 商品管理
  - 商品列表
  - 商品详情
- 📦 订单系统
  - 创建订单
  - 订单列表
  - 订单详情
- ⛓️ 区块链功能
  - 订单交易上链
  - 交易信息查询
  - 区块链存证验证

## 技术栈

- 后端框架：Gin
- 数据库：MySQL
- 区块链存储：LevelDB
- 认证：JWT
- API 文档：Swagger
- 日志：标准库 log

## 系统架构

```
├── cmd/                    # 应用程序入口
│   └── api/               # API服务入口
├── configs/               # 配置文件
├── internal/              # 内部包
│   ├── api/              # API相关
│   ├── blockchain/       # 区块链实现
│   ├── handlers/         # HTTP处理器
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   └── service/         # 业务逻辑层
├── pkg/                  # 公共包
│   ├── errors/          # 错误定义
│   └── logger/          # 日志工具
└── storage/             # 存储相关
    ├── db/              # 区块链数据
    └── logs/            # 应用日志
```

## 快速开始

### 前置要求

- Go 1.23 或更高版本
- MySQL 8.0 或更高版本
- Make (可选)
- goose (数据库迁移工具)

### 安装

1. 克隆项目

```bash
git clone https://github.com/yourusername/blockchain-shop.git
cd blockchain-shop
```

2. 安装依赖

```bash
go mod download
```

3. 安装数据库迁移工具

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

4. 配置数据库

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE IF NOT EXISTS blockchain_shop DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 修改配置文件
cp configs/config.example.yaml configs/config.yaml
# 编辑 config.yaml 设置数据库连接信息
```

5. 执行数据库迁移

```bash
# 使用本地 MySQL
goose -dir migrations mysql "root:123456@tcp(localhost:3306)/blockchain_shop?parseTime=true" up

# 如果使用 Docker 中的 MySQL，替换 localhost 为对应的主机地址
goose -dir migrations mysql "root:123456@tcp(host:3306)/blockchain_shop?parseTime=true" up
```

迁移文件说明：

- `000001_create_users_table.sql`: 创建用户表
- `000002_create_products_table.sql`: 创建商品表
- `000003_create_orders_table.sql`: 创建订单表
- `000004_create_transactions_table.sql`: 创建区块链交易表

6. 运行项目

```bash
make run
# 或者
go run cmd/api/main.go
```

## API 文档

### 用户相关

#### 注册用户

```http
POST /api/v1/register
Content-Type: application/json

{
    "username": "test",
    "password": "password123"
}
```

#### 用户登录

```http
POST /api/v1/login
Content-Type: application/json

{
    "username": "test",
    "password": "password123"
}
```

### 订单相关

#### 创建订单

```http
POST /api/v1/orders
Authorization: Bearer <token>
Content-Type: application/json

{
    "product_id": 1,
    "quantity": 2
}
```

#### 查询订单区块链交易

```http
GET /api/v1/orders/:id/transaction
Authorization: Bearer <token>
```

## 区块链实现

本项目使用简化的区块链实现，主要用于演示订单交易信息的存证功能：

- 使用 SHA256 进行哈希计算
- 使用 LevelDB 存储区块数据
- 实现了基本的区块验证
- 支持交易查询和验证

## 开发规范

- 使用 `go fmt` 格式化代码
- 使用 `go vet` 进行代码检查
- 遵循标准的 Go 项目布局
- 使用依赖注入管理服务依赖

## 测试

运行单元测试：

```bash
make test
# 或者
go test ./...
```

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交改动 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

项目维护者 - [@ylh990835774](https://github.com/ylh990835774)

项目链接: [https://github.com/ylh990835774/blockchain-shop-demo](https://github.com/ylh990835774/blockchain-shop-demo)

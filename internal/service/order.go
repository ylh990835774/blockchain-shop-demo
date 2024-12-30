package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"gorm.io/gorm"
)

type OrderService struct {
	repo           *mysql.OrderRepository
	productRepo    *mysql.ProductRepository
	blockchainRepo *mysql.BlockchainRepository
	db             *gorm.DB
}

func NewOrderService(repo *mysql.OrderRepository, productRepo *mysql.ProductRepository, blockchainRepo *mysql.BlockchainRepository, db *gorm.DB) IOrderService {
	return &OrderService{
		repo:           repo,
		productRepo:    productRepo,
		blockchainRepo: blockchainRepo,
		db:             db,
	}
}

func (s *OrderService) Create(userID, productID int64, quantity int) (*model.Order, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取商品信息（使用事务）
	product, err := s.productRepo.GetByIDWithTx(tx, productID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 检查库存
	if product.Stock < quantity {
		tx.Rollback()
		return nil, ErrInsufficientStock
	}

	// 创建订单
	order := &model.Order{
		UserID:     userID,
		ProductID:  productID,
		Quantity:   quantity,
		TotalPrice: product.Price * float64(quantity),
		Status:     model.OrderStatusPending,
	}

	if err := s.repo.CreateWithTx(tx, order); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 生成区块链交易
	transaction, err := s.blockchainRepo.CreateTransaction(order.ID, order.TotalPrice)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新订单的交易哈希
	order.TxHash = transaction.TxHash
	if err := s.repo.UpdateWithTx(tx, order); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新商品库存
	product.Stock -= quantity
	if err := s.productRepo.UpdateWithTx(tx, product); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return order, nil
}

func (s *OrderService) Get(id int64) (*model.Order, error) {
	return s.repo.GetByID(id)
}

func (s *OrderService) List(userID int64, page, pageSize int) ([]*model.Order, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(userID, offset, pageSize)
}

func (s *OrderService) GetTransaction(orderID int64) (*model.Transaction, error) {
	order, err := s.repo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.TxHash == "" {
		return nil, ErrTransactionNotFound
	}

	return s.blockchainRepo.GetTransaction(order.TxHash)
}

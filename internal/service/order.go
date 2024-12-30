package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
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

func (s *OrderService) Create(order *model.Order) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取商品信息（使用事务）
	product, err := s.productRepo.GetByIDWithTx(tx, order.ProductID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 检查库存
	if product.Stock < order.Quantity {
		tx.Rollback()
		return errors.ErrInsufficientStock
	}

	// 更新商品库存
	if err := s.productRepo.UpdateStockWithTx(tx, order.ProductID, -order.Quantity); err != nil {
		tx.Rollback()
		return err
	}

	// 创建订单
	if err := s.repo.CreateWithTx(tx, order); err != nil {
		tx.Rollback()
		return err
	}

	// 创建区块链交易
	transaction, err := s.blockchainRepo.CreateTransaction(order.ID, order.TotalPrice)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新订单的TxHash
	order.TxHash = transaction.TxHash
	if err := s.repo.UpdateWithTx(tx, order.ID, map[string]interface{}{
		"tx_hash": order.TxHash,
	}); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *OrderService) GetByID(id int64) (*model.Order, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	order, err := s.repo.GetByID(id)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return order, nil
}

func (s *OrderService) ListByUserID(userID int64, page, pageSize int) ([]*model.Order, int64, error) {
	if userID <= 0 {
		return nil, 0, errors.ErrInvalidInput
	}

	return s.repo.ListByUserID(userID, page, pageSize)
}

func (s *OrderService) GetTransaction(orderID int64) (*model.Transaction, error) {
	// 先获取订单信息
	order, err := s.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	// 如果订单没有交易哈希，返回错误
	if order.TxHash == "" {
		return nil, errors.ErrNotFound
	}

	// 获取交易信息
	transaction, err := s.blockchainRepo.GetTransaction(order.TxHash)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

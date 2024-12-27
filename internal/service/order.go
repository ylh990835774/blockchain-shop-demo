package service

import (
	"encoding/json"

	"github.com/ylh990835774/blockchain-shop-demo/internal/blockchain"
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

type OrderService struct {
	repo              *mysql.OrderRepository
	productService    *ProductService
	blockchainService blockchain.Service
}

func NewOrderService(repo *mysql.OrderRepository, productService *ProductService, blockchainService blockchain.Service) *OrderService {
	return &OrderService{
		repo:              repo,
		productService:    productService,
		blockchainService: blockchainService,
	}
}

func (s *OrderService) Create(order *model.Order) error {
	// 检查商品是否存在
	product, err := s.productService.GetByID(order.ProductID)
	if err != nil {
		return errors.ErrNotFound
	}

	// 检查库存
	if product.Stock < order.Quantity {
		return errors.ErrInsufficientStock
	}

	// 计算总价
	order.TotalPrice = float64(order.Quantity) * product.Price

	// 记录区块链交易
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}
	txHash, err := s.blockchainService.RecordTransaction(orderBytes)
	if err != nil {
		return err
	}
	order.TxHash = txHash

	// 开始事务
	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 创建订单
	if err := s.repo.CreateWithTx(tx, order); err != nil {
		tx.Rollback()
		return err
	}

	// 更新库存
	product.Stock -= order.Quantity
	if err := s.productService.UpdateWithTx(tx, product); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return err.Error
	}

	return nil
}

func (s *OrderService) GetByID(id int64) (*model.Order, error) {
	return s.repo.GetByID(id)
}

func (s *OrderService) ListByUserID(userID int64, page, pageSize int) ([]*model.Order, int64, error) {
	return s.repo.ListByUserID(userID, page, pageSize)
}

// GetOrderTransaction 获取订单的区块链交易信息
func (s *OrderService) GetOrderTransaction(orderID int64) ([]byte, error) {
	order, err := s.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.TxHash == "" {
		return nil, errors.ErrNotFound
	}

	return s.blockchainService.GetTransaction(order.TxHash)
}

package mysql

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"gorm.io/gorm"
)

type BlockchainRepository struct {
	db *gorm.DB
}

func NewBlockchainRepository(db *gorm.DB) *BlockchainRepository {
	return &BlockchainRepository{db: db}
}

func (r *BlockchainRepository) GetTransaction(txHash string) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.Where("tx_hash = ?", txHash).First(&tx).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &tx, nil
}

func (r *BlockchainRepository) SaveTransaction(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *BlockchainRepository) CreateTransaction(orderID int64, amount float64) (*model.Transaction, error) {
	// 这里应该调用实际的区块链接口生成交易
	// 为了演示，我们先生成一个模拟的交易
	tx := &model.Transaction{
		TxHash:    generateTxHash(), // 这里需要实现一个生成交易哈希的函数
		From:      "shop_address",   // 商店的区块链地址
		To:        "user_address",   // 用户的区块链地址
		Value:     fmt.Sprintf("%.2f", amount),
		Status:    true,
		Timestamp: time.Now(),
		OrderID:   orderID,
	}

	if err := r.SaveTransaction(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func generateTxHash() string {
	// 生成一个随机的交易哈希
	// 在实际应用中，这个哈希应该由区块链网络生成
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	return fmt.Sprintf("0x%x", randomBytes)
}

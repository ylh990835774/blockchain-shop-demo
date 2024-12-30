package mysql

import (
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

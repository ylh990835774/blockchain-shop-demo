package mysql

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) CreateWithTx(tx *gorm.DB, order *model.Order) error {
	return tx.Create(order).Error
}

func (r *OrderRepository) Update(id int64, updates map[string]interface{}) error {
	result := r.db.Model(&model.Order{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *OrderRepository) UpdateWithTx(tx *gorm.DB, id int64, updates map[string]interface{}) error {
	result := tx.Model(&model.Order{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *OrderRepository) GetByID(id int64) (*model.Order, error) {
	var order model.Order
	err := r.db.First(&order, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListByUserID(userID int64, page, pageSize int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}
